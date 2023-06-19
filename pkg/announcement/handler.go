package announcement

import (
	"context"
	"fmt"
	"time"

	appcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/app"
	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"
	npool "github.com/NpoolPlatform/message/npool/notif/mw/v1/announcement"
	constant "github.com/NpoolPlatform/notif-middleware/pkg/const"
	"github.com/google/uuid"
)

type Handler struct {
	ID      *string
	AppID   *string
	LangID  *string
	UserID  *string
	Title   *string
	Content *string
	Channel *basetypes.NotifChannel
	Type    *npool.AnnouncementType
	EndAt   *uint32
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
		_, err := uuid.Parse(*id)
		if err != nil {
			return err
		}
		h.ID = id
		return nil
	}
}
func WithUserID(userID *string) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if userID == nil {
			return fmt.Errorf("invalid user id")
		}
		_, err := uuid.Parse(*userID)
		if err != nil {
			return err
		}
		h.UserID = userID
		return nil
	}
}

func WithAppID(appID *string) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if appID == nil {
			return fmt.Errorf("invalid app id")
		}
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

func WithLangID(appID, langID *string) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if langID == nil {
			return fmt.Errorf("invalid lang id")
		}
		_, err := uuid.Parse(*langID)
		if err != nil {
			return err
		}
		// wait g11n-middleware merger to master
		// TODO: judge lang id
		// exist, err := appcli.ExistApp(ctx, *appID)
		// if err != nil {
		// 	return err
		// }
		// if !exist {
		// 	return fmt.Errorf("invalid app")
		// }

		h.LangID = langID
		return nil
	}
}

func WithTitle(title *string) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if title == nil {
			return nil
		}
		const leastTitleLen = 4
		if len(*title) < leastTitleLen {
			return fmt.Errorf("name %v too short", *title)
		}
		h.Title = title
		return nil
	}
}

func WithContent(content *string) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if content == nil {
			return nil
		}
		const leastContentLen = 4
		if len(*content) < leastContentLen {
			return fmt.Errorf("content %v too short", *content)
		}
		h.Content = content
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

func WithAnnouncementType(_type *npool.AnnouncementType) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if _type == nil {
			return nil
		}
		switch *_type {
		case npool.AnnouncementType_Broadcast:
		case npool.AnnouncementType_Multicast:
		default:
			return fmt.Errorf("type %v invalid", *_type)
		}
		h.Type = _type
		return nil
	}
}

func WithEndAt(endAt *uint32) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if endAt == nil {
			return nil
		}
		if *endAt < uint32(time.Now().Unix()) {
			return fmt.Errorf("invalid end at")
		}
		h.EndAt = endAt
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
