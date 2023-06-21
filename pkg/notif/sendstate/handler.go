package sendstate

import (
	"context"
	"fmt"

	appmwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/app"
	appusermwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/user"
	"github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"
	notifmwpb "github.com/NpoolPlatform/message/npool/notif/mw/v1/notif"
	constant "github.com/NpoolPlatform/notif-gateway/pkg/const"
	notifmwcli "github.com/NpoolPlatform/notif-middleware/pkg/client/notif"

	"github.com/google/uuid"
)

type Handler struct {
	ID      *string
	AppID   *string
	UserID  *string
	EventID *string
	Channel *basetypes.NotifChannel
	IDs     []string
	Offset  int32
	Limit   int32
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
		if id == nil {
			return nil
		}
		if _, err := uuid.Parse(*id); err != nil {
			return err
		}
		h.ID = id
		return nil
	}
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

func WithEventID(appID, eventID *string) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if eventID == nil {
			return fmt.Errorf("invalid eventid")
		}
		if _, err := uuid.Parse(*eventID); err != nil {
			return err
		}
		exist, err := notifmwcli.ExistNotifConds(ctx, &notifmwpb.Conds{
			AppID: &basetypes.StringVal{
				Op:    cruder.EQ,
				Value: *appID,
			},
			EventID: &basetypes.StringVal{
				Op:    cruder.EQ,
				Value: *eventID,
			},
		})
		if err != nil {
			return err
		}
		if !exist {
			return fmt.Errorf("invalid notif")
		}
		h.EventID = eventID
		return nil
	}
}

func WithChannel(_channel *basetypes.NotifChannel) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if _channel == nil {
			return nil
		}
		switch *_channel {
		case basetypes.NotifChannel_ChannelEmail:
		case basetypes.NotifChannel_ChannelSMS:
		default:
			return fmt.Errorf("channel %v invalid", *_channel)
		}
		h.Channel = _channel
		return nil
	}
}

func WithIDs(ids []string) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		for _, id := range ids {
			if _, err := uuid.Parse(id); err != nil {
				return err
			}
		}
		h.IDs = ids
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
