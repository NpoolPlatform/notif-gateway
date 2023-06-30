package email

import (
	"context"
	"fmt"

	appmwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/app"
	g11nmwcli "github.com/NpoolPlatform/g11n-middleware/pkg/client/applang"
	cruder "github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"
	"github.com/NpoolPlatform/message/npool/g11n/mw/v1/applang"
	constant "github.com/NpoolPlatform/notif-gateway/pkg/const"
	"github.com/google/uuid"
)

type Handler struct {
	ID                *string
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

func WithUsedFor(usedFor *basetypes.UsedFor) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
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

func WithSender(sender *string) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if sender == nil {
			return nil
		}
		if *sender == "" {
			return fmt.Errorf("invalid sender")
		}
		h.Sender = sender
		return nil
	}
}

func WithSubject(subject *string) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if subject == nil {
			return nil
		}
		if *subject == "" {
			return fmt.Errorf("invalid subject")
		}
		h.Subject = subject
		return nil
	}
}

func WithDefaultToUsername(defaultToUsername *string) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if defaultToUsername == nil {
			return nil
		}
		if *defaultToUsername == "" {
			return fmt.Errorf("invalid defaultToUsername")
		}
		h.DefaultToUsername = defaultToUsername
		return nil
	}
}
func WithBody(body *string) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if body == nil {
			return nil
		}
		h.Body = body
		return nil
	}
}

func WithReplyTos(replyTos []string) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		h.ReplyTos = replyTos
		return nil
	}
}

func WithCCTos(cCTos []string) func(context.Context, *Handler) error {
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
