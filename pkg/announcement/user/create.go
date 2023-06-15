package user

import (
	"context"
	"fmt"

	npool "github.com/NpoolPlatform/message/npool/notif/mw/v1/announcement/user"
	cli "github.com/NpoolPlatform/notif-middleware/pkg/client/announcement/user"
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
	if h.AnnouncementID == nil {
		return fmt.Errorf("announcement id is empty")
	}
	return nil
}

func (h *Handler) CreateAnnouncementUser(ctx context.Context) (*npool.AnnouncementUser, error) {
	handler := &createHandler{
		Handler: h,
	}

	if err := handler.validate(); err != nil {
		return nil, err
	}

	info, err := cli.CreateAnnouncementUser(
		ctx,
		&npool.AnnouncementUserReq{
			AppID:          h.AppID,
			UserID:         h.UserID,
			AnnouncementID: h.AnnouncementID,
		},
	)
	if err != nil {
		return nil, err
	}

	h.ID = &info.ID
	return h.GetAnnouncementUser(ctx)
}
