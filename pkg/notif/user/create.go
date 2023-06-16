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
	if h.NotifID == nil {
		return fmt.Errorf("notif id is empty")
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

	// TODO: judge whether exist
	info, err := cli.CreateUser(
		ctx,
		&mwpb.UserNotifReq{
			AppID:   h.AppID,
			UserID:  h.UserID,
			NotifID: h.NotifID,
		},
	)
	if err != nil {
		return nil, err
	}

	h.ID = &info.ID
	return h.GetNotifUser(ctx)
}
