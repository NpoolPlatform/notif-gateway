package notif

import (
	"context"
	"fmt"

	appmwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/app"
	appusermwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/user"
	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"
	templatemwpb "github.com/NpoolPlatform/message/npool/notif/mw/v1/template"
	constant "github.com/NpoolPlatform/notif-gateway/pkg/const"

	"github.com/google/uuid"
)

type Handler struct {
	ID          *string
	AppID       *string
	UserID      *string
	LangID      *string
	EventID     *string
	Notified    *bool
	EventType   *basetypes.UsedFor
	UseTemplate *bool
	Title       *string
	Content     *string
	Channel     *basetypes.NotifChannel
	Extra       *string
	NotifType   *basetypes.NotifType
	Vars        *templatemwpb.TemplateVars
	IDs         []string
	Offset      int32
	Limit       int32
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

func WithLangID(id *string) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if id == nil {
			return nil
		}
		if _, err := uuid.Parse(*id); err != nil {
			return err
		}
		h.LangID = id
		return nil
	}
}

func WithEventID(eventid *string) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if eventid == nil {
			return nil
		}
		_, err := uuid.Parse(*eventid)
		if err != nil {
			return err
		}
		h.EventID = eventid
		return nil
	}
}

func WithNotified(notified *bool) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		h.Notified = notified
		return nil
	}
}

func WithUseTemplate(usetemplate *bool) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		h.UseTemplate = usetemplate
		return nil
	}
}

func WithTitle(title *string) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if title == nil {
			return nil
		}
		if *title == "" {
			return fmt.Errorf("invalid title")
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
		if *content == "" {
			return fmt.Errorf("invalid content")
		}
		h.Content = content
		return nil
	}
}

func WithChannel(_channel *basetypes.NotifChannel) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if _channel == nil {
			return nil
		}
		switch *_channel {
		case basetypes.NotifChannel_ChannelFrontend:
		case basetypes.NotifChannel_ChannelEmail:
		case basetypes.NotifChannel_ChannelSMS:
		default:
			return fmt.Errorf("invalid channel")
		}
		h.Channel = _channel
		return nil
	}
}

func WithExtra(extra *string) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if extra == nil {
			return nil
		}
		if *extra == "" {
			return fmt.Errorf("invalid extra")
		}
		h.Extra = extra
		return nil
	}
}

func WithEventType(eventtype *basetypes.UsedFor) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if eventtype == nil {
			return nil
		}
		switch *eventtype {
		case basetypes.UsedFor_WithdrawalRequest:
		case basetypes.UsedFor_WithdrawalCompleted:
		case basetypes.UsedFor_DepositReceived:
		case basetypes.UsedFor_KYCApproved:
		case basetypes.UsedFor_KYCRejected:
		case basetypes.UsedFor_Announcement:
		default:
			return fmt.Errorf("invalid eventtype")
		}
		h.EventType = eventtype
		return nil
	}
}

func WithNotifType(_type *basetypes.NotifType) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if _type == nil {
			return nil
		}
		switch *_type {
		case basetypes.NotifType_NotifBroadcast:
		case basetypes.NotifType_NotifMulticast:
		case basetypes.NotifType_NotifUnicast:
			if h.UserID == nil {
				return fmt.Errorf("invalid userid")
			}
		default:
			return fmt.Errorf("invalid type")
		}
		h.NotifType = _type
		return nil
	}
}

func WithVars(vars *templatemwpb.TemplateVars) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		h.Vars = vars
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
