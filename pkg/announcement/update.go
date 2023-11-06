package announcement

import (
	"context"
	"fmt"

	"github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"
	npool "github.com/NpoolPlatform/message/npool/notif/gw/v1/announcement"
	mwpb "github.com/NpoolPlatform/message/npool/notif/mw/v1/announcement"
	cli "github.com/NpoolPlatform/notif-middleware/pkg/client/announcement"
)

func (h *Handler) UpdateAnnouncement(ctx context.Context) (*npool.Announcement, error) {
	info, err := cli.GetAnnouncementOnly(ctx, &mwpb.Conds{
		ID:    &basetypes.Uint32Val{Op: cruder.EQ, Value: *h.ID},
		EntID: &basetypes.StringVal{Op: cruder.EQ, Value: *h.EntID},
		AppID: &basetypes.StringVal{Op: cruder.EQ, Value: *h.AppID},
	})
	if err != nil {
		return nil, err
	}
	if info == nil {
		return nil, fmt.Errorf("announcement not found")
	}

	if h.StartAt != nil && h.EndAt != nil {
		if *h.StartAt >= *h.EndAt {
			return nil, fmt.Errorf("start at less than end at")
		}
	}
	if h.StartAt != nil && h.EndAt == nil {
		if *h.StartAt > info.EndAt {
			return nil, fmt.Errorf("start at less than end at")
		}
	}
	if h.EndAt != nil && h.StartAt == nil {
		if *h.EndAt < info.StartAt {
			return nil, fmt.Errorf("start at less than end at")
		}
	}

	_, err = cli.UpdateAnnouncement(ctx, &mwpb.AnnouncementReq{
		ID:               h.ID,
		Title:            h.Title,
		Content:          h.Content,
		EndAt:            h.EndAt,
		StartAt:          h.StartAt,
		AnnouncementType: h.Type,
	},
	)
	if err != nil {
		return nil, err
	}

	h.EntID = &info.EntID

	return h.GetAnnouncement(ctx)
}
