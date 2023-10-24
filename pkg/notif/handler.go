package notif

import (
	"context"
	"fmt"

	appmwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/app"
	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"
	notifmw "github.com/NpoolPlatform/message/npool/notif/mw/v1/notif"
	templatemwpb "github.com/NpoolPlatform/message/npool/notif/mw/v1/template"
	constant "github.com/NpoolPlatform/notif-gateway/pkg/const"

	"github.com/google/uuid"
)

type Handler struct {
	ID          *uint32
	EntID       *string
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
	Reqs        []*notifmw.NotifReq
	EntIDs      []string
	IDs         []uint32
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
			return fmt.Errorf("invalid lang id")
		}
		_, err := uuid.Parse(*id)
		if err != nil {
			return err
		}

		h.LangID = id
		return nil
	}
}

func WithEventID(eventid *string, must bool) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if eventid == nil {
			if must {
				return fmt.Errorf("invalid eventid")
			}
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

func WithNotified(notified *bool, must bool) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if notified == nil {
			if must {
				return fmt.Errorf("invalid notified")
			}
			return nil
		}
		h.Notified = notified
		return nil
	}
}

func WithUseTemplate(usetemplate *bool, must bool) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if usetemplate == nil {
			if must {
				return fmt.Errorf("invalid usetemplate")
			}
			return nil
		}
		h.UseTemplate = usetemplate
		return nil
	}
}

func WithTitle(title *string, must bool) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if title == nil {
			if must {
				return fmt.Errorf("invalid title")
			}
			return nil
		}
		if *title == "" {
			return fmt.Errorf("invalid title")
		}
		h.Title = title
		return nil
	}
}

func WithContent(content *string, must bool) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if content == nil {
			if must {
				return fmt.Errorf("invalid content")
			}
			return nil
		}
		if *content == "" {
			return fmt.Errorf("invalid content")
		}
		h.Content = content
		return nil
	}
}

func WithChannel(_channel *basetypes.NotifChannel, must bool) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if _channel == nil {
			if must {
				return fmt.Errorf("invalid channel")
			}
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

func WithExtra(extra *string, must bool) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if extra == nil {
			if must {
				return fmt.Errorf("invalid extra")
			}
			return nil
		}
		if *extra == "" {
			return fmt.Errorf("invalid extra")
		}
		h.Extra = extra
		return nil
	}
}

//nolint:gocyclo
func WithEventType(eventtype *basetypes.UsedFor, must bool) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if eventtype == nil {
			if must {
				return fmt.Errorf("invalid eventtype")
			}
			return nil
		}
		switch *eventtype {
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
			return fmt.Errorf("invalid eventtype")
		}
		h.EventType = eventtype
		return nil
	}
}

func WithNotifType(_type *basetypes.NotifType, must bool) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if _type == nil {
			if must {
				return fmt.Errorf("invalid notiftype")
			}
			return nil
		}
		switch *_type {
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

func WithVars(vars *templatemwpb.TemplateVars, must bool) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if vars == nil {
			if must {
				return fmt.Errorf("invalid vars")
			}
			return nil
		}
		h.Vars = vars
		return nil
	}
}

func WithEntIDs(ids []string, must bool) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		for _, id := range ids {
			if _, err := uuid.Parse(id); err != nil {
				return err
			}
		}
		h.EntIDs = ids
		return nil
	}
}

func WithIDs(ids []uint32, must bool) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		h.IDs = ids
		return nil
	}
}

func WithReqs(reqs []*notifmw.NotifReq) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if len(reqs) == 0 {
			return fmt.Errorf("infos is empty")
		}
		for _, req := range reqs {
			if req.ID == nil {
				return fmt.Errorf("invalid id")
			}
			if req.Notified == nil {
				return fmt.Errorf("invalid notified")
			}
			if !*req.Notified {
				return fmt.Errorf("invalid notified %v", *req.Notified)
			}
		}
		h.Reqs = reqs
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
