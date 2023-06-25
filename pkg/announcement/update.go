package announcement

import (
	"context"
	"fmt"

	npool "github.com/NpoolPlatform/message/npool/notif/gw/v1/announcement"
	mwpb "github.com/NpoolPlatform/message/npool/notif/mw/v1/announcement"
	cli "github.com/NpoolPlatform/notif-middleware/pkg/client/announcement"
)

func (h *Handler) UpdateAnnouncement(ctx context.Context) (*npool.Announcement, error) {
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

	_, err = cli.UpdateAnnouncement(ctx, &mwpb.AnnouncementReq{
		ID:               h.ID,
		Title:            h.Title,
		Content:          h.Content,
		EndAt:            h.EndAt,
		AnnouncementType: h.Type,
	},
	)
	if err != nil {
		return nil, err
	}

	return h.GetAnnouncement(ctx)
}
