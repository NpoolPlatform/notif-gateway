package usercode

import (
	"context"
	"fmt"

	usermwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/user"

	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"

	sendmwpb "github.com/NpoolPlatform/message/npool/third/mw/v1/send"
	sendmwcli "github.com/NpoolPlatform/third-middleware/pkg/client/send"

	tmplmwpb "github.com/NpoolPlatform/message/npool/notif/mw/v1/template"
	tmplmwcli "github.com/NpoolPlatform/notif-middleware/pkg/client/template"

	chanmgrpb "github.com/NpoolPlatform/message/npool/notif/mgr/v1/channel"

	usercodemwcli "github.com/NpoolPlatform/basal-middleware/pkg/client/usercode"
	usercodemwpb "github.com/NpoolPlatform/message/npool/basal/mw/v1/usercode"
)

func SendCode( //nolint
	ctx context.Context,
	appID, langID string,
	userID, account *string,
	accountType basetypes.SignMethod,
	usedFor basetypes.UsedFor,
	toUsername *string,
) error {
	switch usedFor {
	case basetypes.UsedFor_SetWithdrawAddress:
		fallthrough //nolint
	case basetypes.UsedFor_Withdraw:
		fallthrough //nolint
	case basetypes.UsedFor_CreateInvitationCode:
		fallthrough //nolint
	case basetypes.UsedFor_SetCommission:
		fallthrough //nolint
	case basetypes.UsedFor_SetTransferTargetUser:
		fallthrough //nolint
	case basetypes.UsedFor_Transfer:
		if userID != nil && *userID != "" {
			user, err := usermwcli.GetUser(ctx, appID, *userID)
			if err != nil {
				return err
			}
			if user == nil {
				return fmt.Errorf("invalid user")
			}
			switch accountType {
			case basetypes.SignMethod_Mobile:
				account = &user.PhoneNO
			case basetypes.SignMethod_Email:
				account = &user.EmailAddress
			}
			if toUsername == nil || *toUsername == "" {
				toUsername = &user.Username
			}
		}
	}

	if account == nil || *account == "" {
		return fmt.Errorf("invalid account")
	}

	channel := chanmgrpb.NotifChannel_ChannelEmail
	switch accountType {
	case basetypes.SignMethod_Email:
	case basetypes.SignMethod_Mobile:
		channel = chanmgrpb.NotifChannel_ChannelSMS
	default:
		return fmt.Errorf("invalid account type")
	}

	code, err := usercodemwcli.CreateUserCode(ctx, &usercodemwpb.CreateUserCodeRequest{
		Prefix:      basetypes.Prefix_PrefixUserCode.String(),
		AppID:       appID,
		Account:     *account,
		AccountType: accountType,
		UsedFor:     usedFor,
	})
	if err != nil {
		return err
	}

	info, err := tmplmwcli.GenerateText(ctx, &tmplmwpb.GenerateTextRequest{
		AppID:     appID,
		LangID:    langID,
		Channel:   channel,
		EventType: usedFor,
		Vars: &tmplmwpb.TemplateVars{
			Username: toUsername,
			Code:     &code.Code,
		},
	})
	if err != nil {
		return err
	}
	if info == nil {
		return fmt.Errorf("cannot generate text")
	}

	err = sendmwcli.SendMessage(ctx, &sendmwpb.SendMessageRequest{
		Subject:     info.Subject,
		Content:     info.Content,
		From:        info.From,
		To:          *account,
		ToCCs:       info.ToCCs,
		ReplyTos:    info.ReplyTos,
		AccountType: accountType,
	})
	if err != nil {
		return err
	}

	return nil
}
