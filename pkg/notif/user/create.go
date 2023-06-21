package user

import (
	"context"
	"fmt"

	npool "github.com/NpoolPlatform/message/npool/notif/gw/v1/notif/user"
	mwpb "github.com/NpoolPlatform/message/npool/notif/mw/v1/notif/user"
	cli "github.com/NpoolPlatform/notif-middleware/pkg/client/notif/user"
)

type createHandler struct {
	*Handler
}

func (h *createHandler) validate() error {
	if h.AppID == nil {
		return fmt.Errorf("app id is empty")
	}
	if h.UserID == nil {
		return fmt.Errorf("user id is empty")
	}
	if h.EventType == nil {
		return fmt.Errorf("eventtype is empty")
	}
	return nil
}

func (h *Handler) CreateNotifUser(ctx context.Context) (*npool.NotifUser, error) {
	handler := &createHandler{
		Handler: h,
	}

	if err := handler.validate(); err != nil {
		return nil, err
	}

	info, err := cli.CreateNotifUser(
		ctx,
		&mwpb.NotifUserReq{
			AppID:     h.AppID,
			UserID:    h.UserID,
			EventType: h.EventType,
		},
	)
	if err != nil {
		return nil, err
	}

	h.ID = &info.ID
	return h.GetNotifUser(ctx)
}
