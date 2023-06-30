package usercode

import (
	"context"
	"fmt"

	appmwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/app"
	appusermwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/user"
	g11nmwcli "github.com/NpoolPlatform/g11n-middleware/pkg/client/applang"
	cruder "github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"
	"github.com/NpoolPlatform/message/npool/g11n/mw/v1/applang"

	"github.com/google/uuid"
)

type Handler struct {
	AppID       *string
	LangID      *string
	UserID      *string
	Account     *string
	AccountType *basetypes.SignMethod
	UsedFor     *basetypes.UsedFor
	ToUsername  *string
}

func NewHandler(ctx context.Context, options ...func(context.Context, *Handler) error) (*Handler, error) {
	handler := &Handler{}
	for _, opt := range options {
		if err := opt(ctx, handler); err != nil {
			return nil, err
		}
	}
	return handler, nil
}

func WithAppID(appID *string) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if appID == nil {
			return nil
		}
		if _, err := uuid.Parse(*appID); err != nil {
			return err
		}
		exist, err := appmwcli.ExistApp(ctx, *appID)
		if err != nil {
			return err
		}
		if !exist {
			return fmt.Errorf("invalid app")
		}
		h.AppID = appID
		return nil
	}
}

func WithUserID(appID, userID *string) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if appID == nil || userID == nil {
			return nil
		}
		_, err := uuid.Parse(*userID)
		if err != nil {
			return err
		}
		exist, err := appusermwcli.ExistUser(ctx, *appID, *userID)
		if err != nil {
			return err
		}
		if !exist {
			return fmt.Errorf("invalid user")
		}

		h.UserID = userID
		return nil
	}
}

func WithLangID(appID, langID *string) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if langID == nil {
			return fmt.Errorf("invalid lang id")
		}
		_, err := uuid.Parse(*langID)
		if err != nil {
			return err
		}

		exist, err := g11nmwcli.ExistAppLangConds(ctx, &applang.Conds{
			AppID: &basetypes.StringVal{
				Op:    cruder.EQ,
				Value: *appID,
			},
			LangID: &basetypes.StringVal{
				Op:    cruder.EQ,
				Value: *langID,
			},
		})

		if err != nil {
			return err
		}
		if !exist {
			return fmt.Errorf("invalid lang id")
		}

		h.LangID = langID
		return nil
	}
}

func WithAccount(account *string) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		h.Account = account
		return nil
	}
}

func WithAccountType(accountType *basetypes.SignMethod) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if accountType == nil {
			return nil
		}
		switch *accountType {
		case basetypes.SignMethod_Email:
		case basetypes.SignMethod_Mobile:
		default:
			return fmt.Errorf("invalid accountType")
		}
		h.AccountType = accountType
		return nil
	}
}

func WithUsedFor(usedFor *basetypes.UsedFor) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		switch *usedFor {
		case basetypes.UsedFor_Signup:
		case basetypes.UsedFor_Signin:
		case basetypes.UsedFor_Update:
		case basetypes.UsedFor_SetWithdrawAddress:
		case basetypes.UsedFor_Withdraw:
		case basetypes.UsedFor_CreateInvitationCode:
		case basetypes.UsedFor_SetCommission:
		case basetypes.UsedFor_SetTransferTargetUser:
		case basetypes.UsedFor_Transfer:
		default:
			return fmt.Errorf("invalid usedFor")
		}
		h.UsedFor = usedFor
		return nil
	}
}

func WithToUsername(toUsername *string) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		h.ToUsername = toUsername
		return nil
	}
}
