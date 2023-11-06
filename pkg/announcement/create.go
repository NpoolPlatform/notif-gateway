package announcement

import (
	"context"
	"fmt"

	applangmwcli "github.com/NpoolPlatform/g11n-middleware/pkg/client/applang"
	"github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"
	"github.com/NpoolPlatform/message/npool/g11n/mw/v1/applang"
	npool "github.com/NpoolPlatform/message/npool/notif/gw/v1/announcement"
	mwpb "github.com/NpoolPlatform/message/npool/notif/mw/v1/announcement"
	cli "github.com/NpoolPlatform/notif-middleware/pkg/client/announcement"
)

func (h *Handler) CreateAnnouncement(ctx context.Context) (*npool.Announcement, error) {
	exist, err := applangmwcli.ExistAppLangConds(ctx, &applang.Conds{
		AppID:  &basetypes.StringVal{Op: cruder.EQ, Value: *h.AppID},
		LangID: &basetypes.StringVal{Op: cruder.EQ, Value: *h.LangID},
	})
	if err != nil {
		return nil, err
	}
	if !exist {
		return nil, fmt.Errorf("invalid applang")
	}

	info, err := cli.CreateAnnouncement(
		ctx,
		&mwpb.AnnouncementReq{
			AppID:            h.AppID,
			Title:            h.Title,
			Content:          h.Content,
			LangID:           h.LangID,
			Channel:          h.Channel,
			AnnouncementType: h.Type,
			StartAt:          h.StartAt,
			EndAt:            h.EndAt,
		},
	)
	if err != nil {
		return nil, err
	}

	h.ID = &info.ID
	h.EntID = &info.EntID
	return h.GetAnnouncement(ctx)
}
