package channel

import (
	"context"
	"fmt"

	appcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/app"
	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"
	constant "github.com/NpoolPlatform/notif-middleware/pkg/const"
	"github.com/google/uuid"
)

type Handler struct {
	ID         *string
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

func WithID(id *string) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		_, err := uuid.Parse(*id)
		if err != nil {
			return err
		}
		h.ID = id
		return nil
	}
}

func WithAppID(appID *string) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		_, err := uuid.Parse(*appID)
		if err != nil {
			return err
		}
		exist, err := appcli.ExistApp(ctx, *appID)
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

func WithChannel(channel *basetypes.NotifChannel) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
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

func WithEventType(_type *basetypes.UsedFor) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if _type == nil {
			return nil
		}
		switch *_type {
		case basetypes.UsedFor_WithdrawalRequest:
		case basetypes.UsedFor_WithdrawalCompleted:
		case basetypes.UsedFor_DepositReceived:
		case basetypes.UsedFor_KYCApproved:
		case basetypes.UsedFor_KYCRejected:
		case basetypes.UsedFor_Announcement:
		default:
			return fmt.Errorf("EventType is invalid")
		}

		h.EventType = _type
		return nil
	}
}

func WithEventTypes(_types []basetypes.UsedFor) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if _types == nil {
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
			default:
				return fmt.Errorf("EventType is invalid")
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
