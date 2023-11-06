package channel

import (
	"context"
	"fmt"

	appmwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/app"
	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"
	constant "github.com/NpoolPlatform/notif-middleware/pkg/const"
	"github.com/google/uuid"
)

type Handler struct {
	ID         *uint32
	EntID      *string
	AppID      *string
	Channel    *basetypes.NotifChannel
	EventType  *basetypes.UsedFor
	EventTypes []basetypes.UsedFor
	Offset     int32
	Limit      int32
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

func WithID(id *uint32, must bool) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if id == nil {
			if must {
				return fmt.Errorf("invalid id")
			}
			return nil
		}
		h.ID = id
		return nil
	}
}

func WithEntID(id *string, must bool) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if id == nil {
			if must {
				return fmt.Errorf("invalid entid")
			}
			return nil
		}
		_, err := uuid.Parse(*id)
		if err != nil {
			return err
		}
		h.EntID = id
		return nil
	}
}

func WithAppID(id *string, must bool) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if id == nil {
			if must {
				return fmt.Errorf("invalid appid")
			}
			return nil
		}
		_, err := uuid.Parse(*id)
		if err != nil {
			return err
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

func WithChannel(channel *basetypes.NotifChannel, must bool) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if channel == nil {
			if must {
				return fmt.Errorf("invalid channel")
			}
			return nil
		}
		switch *channel {
		case basetypes.NotifChannel_ChannelEmail:
		case basetypes.NotifChannel_ChannelSMS:
		case basetypes.NotifChannel_ChannelFrontend:
		default:
			return fmt.Errorf("channel %v invalid", *channel)
		}
		h.Channel = channel
		return nil
	}
}

//nolint:gocyclo
func WithEventType(_type *basetypes.UsedFor, must bool) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if _type == nil {
			if must {
				return fmt.Errorf("invalid eventtype")
			}
			return nil
		}
		switch *_type {
		case basetypes.UsedFor_WithdrawalRequest:
		case basetypes.UsedFor_WithdrawalCompleted:
		case basetypes.UsedFor_DepositReceived:
		case basetypes.UsedFor_KYCApproved:
		case basetypes.UsedFor_KYCRejected:
		case basetypes.UsedFor_Announcement:
		case basetypes.UsedFor_GoodBenefit1:
		case basetypes.UsedFor_UpdateEmail:
		case basetypes.UsedFor_UpdateMobile:
		case basetypes.UsedFor_UpdatePassword:
		case basetypes.UsedFor_UpdateGoogleAuth:
		case basetypes.UsedFor_NewLogin:
		case basetypes.UsedFor_OrderCompleted:
		default:
			return fmt.Errorf("EventType is invalid")
		}

		h.EventType = _type
		return nil
	}
}

//nolint:gocyclo
func WithEventTypes(_types []basetypes.UsedFor, must bool) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if len(_types) == 0 {
			if must {
				return fmt.Errorf("invalid eventtypes")
			}
			return nil
		}
		for _, _type := range _types {
			switch _type {
			case basetypes.UsedFor_WithdrawalRequest:
			case basetypes.UsedFor_WithdrawalCompleted:
			case basetypes.UsedFor_DepositReceived:
			case basetypes.UsedFor_KYCApproved:
			case basetypes.UsedFor_KYCRejected:
			case basetypes.UsedFor_Announcement:
			case basetypes.UsedFor_UpdateEmail:
			case basetypes.UsedFor_UpdateMobile:
			case basetypes.UsedFor_UpdatePassword:
			case basetypes.UsedFor_UpdateGoogleAuth:
			case basetypes.UsedFor_NewLogin:
			case basetypes.UsedFor_OrderCompleted:
			default:
				return fmt.Errorf("event type is invalid")
			}
		}

		h.EventTypes = _types
		return nil
	}
}

func WithOffset(offset int32) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		h.Offset = offset
		return nil
	}
}

func WithLimit(limit int32) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if limit == 0 {
			limit = constant.DefaultRowLimit
		}
		h.Limit = limit
		return nil
	}
}
