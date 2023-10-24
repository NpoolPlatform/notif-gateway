package usercode

import (
	"context"
	"fmt"

	appmwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/app"
	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"

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

func WithAppID(id *string, must bool) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if id == nil {
			if must {
				return fmt.Errorf("invalid appid")
			}
			return nil
		}
		exist, err := appmwcli.ExistApp(ctx, *id)
		if err != nil {
			return err
		}
		if !exist {
			return fmt.Errorf("invalid app")
		}
		h.AppID = id
		return nil
	}
}

func WithUserID(id *string, must bool) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if id == nil {
			if must {
				return fmt.Errorf("invalid userid")
			}
			return nil
		}
		if _, err := uuid.Parse(*id); err != nil {
			return err
		}
		h.UserID = id
		return nil
	}
}

func WithLangID(id *string, must bool) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if id == nil {
			if must {
				return fmt.Errorf("invalid langid")
			}
			return nil
		}
		_, err := uuid.Parse(*id)
		if err != nil {
			return err
		}

		h.LangID = id
		return nil
	}
}

func WithAccount(account *string, must bool) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if account == nil {
			if must {
				return fmt.Errorf("invalid account")
			}
			return nil
		}
		h.Account = account
		return nil
	}
}

func WithAccountType(accountType *basetypes.SignMethod, must bool) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if accountType == nil {
			if must {
				return fmt.Errorf("invalid accounttype")
			}
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

func WithUsedFor(usedFor *basetypes.UsedFor, must bool) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if usedFor == nil {
			if must {
				return fmt.Errorf("invalid usedfor")
			}
			return nil
		}
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

func WithToUsername(toUsername *string, must bool) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if toUsername == nil {
			if must {
				return fmt.Errorf("invalid tousername")
			}
			return nil
		}
		h.ToUsername = toUsername
		return nil
	}
}
