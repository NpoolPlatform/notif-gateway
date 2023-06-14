package announcement

import (
	"context"

	npool "github.com/NpoolPlatform/message/npool/notif/mw/v1/announcement"
	cli "github.com/NpoolPlatform/notif-middleware/pkg/client/announcement"
)

func (h *Handler) DeleteAnnouncement(ctx context.Context) (*npool.Announcement, error) {
	info, err := h.GetAnnouncement(ctx)
	if err != nil {
		return nil, err
	}

	_, err = cli.DeleteAnnouncement(ctx, *h.ID)
	if err != nil {
		return nil, err
	}

	return info, nil
}
