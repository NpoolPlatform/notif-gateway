package user

import (
	"context"

	npool "github.com/NpoolPlatform/message/npool/notif/mw/v1/announcement/user"
	cli "github.com/NpoolPlatform/notif-middleware/pkg/client/announcement/user"
)

func (h *Handler) DeleteAnnouncementUser(ctx context.Context) (*npool.AnnouncementUser, error) {
	info, err := h.GetAnnouncementUser(ctx)
	if err != nil {
		return nil, err
	}

	_, err = cli.DeleteAnnouncementUser(ctx, *h.ID)
	if err != nil {
		return nil, err
	}

	return info, nil
}
