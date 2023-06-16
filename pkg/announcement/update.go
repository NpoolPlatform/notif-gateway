package announcement

import (
	"context"

	npool "github.com/NpoolPlatform/message/npool/notif/gw/v1/announcement"
	mwpb "github.com/NpoolPlatform/message/npool/notif/mw/v1/announcement"
	cli "github.com/NpoolPlatform/notif-middleware/pkg/client/announcement"
)

func (h *Handler) UpdateAnnouncement(ctx context.Context) (*npool.Announcement, error) {
	_, err := cli.UpdateAnnouncement(ctx, &mwpb.AnnouncementReq{
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
