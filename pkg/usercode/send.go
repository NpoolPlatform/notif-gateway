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

	usercodemwcli "github.com/NpoolPlatform/basal-middleware/pkg/client/usercode"
	usercodemwpb "github.com/NpoolPlatform/message/npool/basal/mw/v1/usercode"
)

func (h *Handler) validateUser(ctx context.Context) error {
	user, err := usermwcli.GetUser(ctx, *h.AppID, *h.UserID)
	if err != nil {
		return err
	}
	if user == nil {
		return fmt.Errorf("invalid user")
	}
	switch *h.AccountType {
	case basetypes.SignMethod_Mobile:
		h.Account = &user.PhoneNO
	case basetypes.SignMethod_Email:
		h.Account = &user.EmailAddress
	}
	if h.ToUsername == nil || *h.ToUsername == "" {
		h.ToUsername = &user.Username
	}
	return nil
}

func (h *Handler) SendCode( //nolint
	ctx context.Context,
) error {
	switch *h.UsedFor {
	case basetypes.UsedFor_Signup:
	case basetypes.UsedFor_Update:
		if h.Account == nil || *h.Account == "" {
			if h.UserID == nil || *h.UserID == "" {
				return fmt.Errorf("invalid userid")
			}
			err := h.validateUser(ctx)
			if err != nil {
				return err
			}
		}
	case basetypes.UsedFor_Signin:
		fallthrough //nolint
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
		if h.UserID == nil || *h.UserID == "" {
			return fmt.Errorf("invalid userid")
		}
		err := h.validateUser(ctx)
		if err != nil {
			return err
		}
	}

	if h.Account == nil || *h.Account == "" {
		return fmt.Errorf("invalid account")
	}

	channel := basetypes.NotifChannel_ChannelEmail
	switch *h.AccountType {
	case basetypes.SignMethod_Email:
	case basetypes.SignMethod_Mobile:
		channel = basetypes.NotifChannel_ChannelSMS
	default:
		return fmt.Errorf("invalid account type")
	}

	code, err := usercodemwcli.CreateUserCode(ctx, &usercodemwpb.CreateUserCodeRequest{
		Prefix:      basetypes.Prefix_PrefixUserCode.String(),
		AppID:       *h.AppID,
		Account:     *h.Account,
		AccountType: *h.AccountType,
		UsedFor:     *h.UsedFor,
	})
	if err != nil {
		return err
	}

	info, err := tmplmwcli.GenerateText(ctx, &tmplmwpb.GenerateTextRequest{
		AppID:     *h.AppID,
		LangID:    *h.LangID,
		Channel:   channel,
		EventType: *h.UsedFor,
		Vars: &tmplmwpb.TemplateVars{
			Username: h.ToUsername,
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
		To:          *h.Account,
		ToCCs:       info.ToCCs,
		ReplyTos:    info.ReplyTos,
		AccountType: *h.AccountType,
	})
	if err != nil {
		return err
	}

	return nil
}
