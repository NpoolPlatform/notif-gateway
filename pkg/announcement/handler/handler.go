package handler

import (
	"context"
	"fmt"

	appcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/app"
	appusercli "github.com/NpoolPlatform/appuser-middleware/pkg/client/user"
	constant "github.com/NpoolPlatform/notif-gateway/pkg/const"
	cli "github.com/NpoolPlatform/notif-middleware/pkg/client/announcement"
	"github.com/google/uuid"
)

type Handler struct {
	ID             *string
	AppID          *string
	UserID         *string
	AnnouncementID *string
	Offset         int32
	Limit          int32
}

func NewHandler(ctx context.Context, options ...interface{}) (*Handler, error) {
	handler := &Handler{}
	for _, opt := range options {
		_opt, ok := opt.(func(context.Context, *Handler) error)
		if !ok {
			continue
		}
		if err := _opt(ctx, handler); err != nil {
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

func WithUserID(appID, userID *string) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if userID == nil {
			return fmt.Errorf("invalid userid")
		}
		_, err := uuid.Parse(*userID)
		if err != nil {
			return err
		}

		exist, err := appusercli.ExistUser(ctx, *appID, *userID)
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

func WithAnnouncementID(appID, amtID *string) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if amtID == nil {
			return fmt.Errorf("invalid announcement id")
		}
		exist, err := cli.ExistAnnouncement(ctx, *amtID, *appID)
		if err != nil {
			return err
		}
		if !exist {
			return fmt.Errorf("announcement id not exist")
		}

		_, err = uuid.Parse(*amtID)
		if err != nil {
			return err
		}
		h.AnnouncementID = amtID
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
