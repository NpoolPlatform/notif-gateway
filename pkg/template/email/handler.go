package email

import (
	"context"
	"fmt"

	appmwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/app"
	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"
	constant "github.com/NpoolPlatform/notif-gateway/pkg/const"
	"github.com/google/uuid"
)

type Handler struct {
	ID                *uint32
	EntID             *string
	AppID             *string
	LangID            *string
	UsedFor           *basetypes.UsedFor
	Sender            *string
	Subject           *string
	Body              *string
	DefaultToUsername *string
	ReplyTos          []string
	CCTos             []string
	Offset            int32
	Limit             int32
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

func WithUsedFor(usedFor *basetypes.UsedFor, must bool) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if usedFor == nil {
			if must {
				return fmt.Errorf("invalid usedfor")
			}
			return nil
		}
		_usedFor := false
		for key := range basetypes.UsedFor_value {
			if key == usedFor.String() && *usedFor != basetypes.UsedFor_DefaultUsedFor {
				_usedFor = true
			}
		}
		if !_usedFor {
			return fmt.Errorf("usedFor is invalid")
		}
		h.UsedFor = usedFor
		return nil
	}
}

func WithSender(sender *string, must bool) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if sender == nil {
			if must {
				return fmt.Errorf("invalid sender")
			}
			return nil
		}
		if *sender == "" {
			return fmt.Errorf("invalid sender")
		}
		h.Sender = sender
		return nil
	}
}

func WithSubject(subject *string, must bool) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if subject == nil {
			if must {
				return fmt.Errorf("invalid subject")
			}
			return nil
		}
		if *subject == "" {
			return fmt.Errorf("invalid subject")
		}
		h.Subject = subject
		return nil
	}
}

func WithDefaultToUsername(defaultToUsername *string, must bool) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if defaultToUsername == nil {
			if must {
				return fmt.Errorf("invalid defaulttousername")
			}
			return nil
		}
		if *defaultToUsername == "" {
			return fmt.Errorf("invalid defaultToUsername")
		}
		h.DefaultToUsername = defaultToUsername
		return nil
	}
}
func WithBody(body *string, must bool) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if body == nil {
			if must {
				return fmt.Errorf("invalid body")
			}
			return nil
		}
		h.Body = body
		return nil
	}
}

func WithReplyTos(replyTos []string, must bool) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		h.ReplyTos = replyTos
		return nil
	}
}

func WithCCTos(cCTos []string, must bool) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		h.CCTos = cCTos
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
