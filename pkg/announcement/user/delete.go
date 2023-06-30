package user

import (
	"context"
	"fmt"

	npool "github.com/NpoolPlatform/message/npool/notif/gw/v1/announcement/user"
	cli "github.com/NpoolPlatform/notif-middleware/pkg/client/announcement/user"
)

func (h *Handler) DeleteAnnouncementUser(ctx context.Context) (*npool.AnnouncementUser, error) {
	info, err := h.GetAnnouncementUser(ctx)
	if err != nil {
		return nil, err
	}
	if info == nil {
		return nil, fmt.Errorf(" announcement user not found")
	}
	if info.AppID != *h.AppID {
		return nil, fmt.Errorf("permission denied")
	}

	_, err = cli.DeleteAnnouncementUser(ctx, *h.ID)
	if err != nil {
		return nil, err
	}

	return info, nil
}
