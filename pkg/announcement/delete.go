package announcement

import (
	"context"
	"fmt"

	npool "github.com/NpoolPlatform/message/npool/notif/gw/v1/announcement"
	cli "github.com/NpoolPlatform/notif-middleware/pkg/client/announcement"
)

func (h *Handler) DeleteAnnouncement(ctx context.Context) (*npool.Announcement, error) {
	info, err := h.GetAnnouncement(ctx)
	if err != nil {
		return nil, err
	}
	if info == nil {
		return nil, fmt.Errorf("announcement not found")
	}

	if info.AppID != *h.AppID {
		return nil, fmt.Errorf("permission denied")
	}

	_, err = cli.DeleteAnnouncement(ctx, *h.ID)
	if err != nil {
		return nil, err
	}

	return info, nil
}
