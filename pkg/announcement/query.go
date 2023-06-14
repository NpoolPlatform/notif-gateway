package announcement

import (
	"context"
	"fmt"

	"github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"
	mwpb "github.com/NpoolPlatform/message/npool/notif/mw/v1/announcement"
	npool "github.com/NpoolPlatform/message/npool/notif/mw/v1/announcement"
	cli "github.com/NpoolPlatform/notif-middleware/pkg/client/announcement"
	mwcli "github.com/NpoolPlatform/notif-middleware/pkg/client/announcement"
)

func (h *Handler) GetAnnouncements(ctx context.Context) ([]*npool.Announcement, uint32, error) {
	conds := &mwpb.Conds{
		AppID: &basetypes.StringVal{
			Op:    cruder.EQ,
			Value: *h.AppID,
		},
	}
	if h.UserID != nil {
		conds.UserID = &basetypes.StringVal{
			Op:    cruder.EQ,
			Value: *h.UserID,
		}
	}
	if h.LangID != nil {
		conds.LangID = &basetypes.StringVal{
			Op:    cruder.EQ,
			Value: *h.LangID,
		}
	}

	infos, total, err := mwcli.GetAnnouncementStates(ctx, conds, h.Offset, h.Limit)
	if err != nil {
		return nil, 0, err
	}
	if len(infos) == 0 {
		return nil, total, nil
	}

	return infos, total, nil
}

func (h *Handler) GetAnnouncement(ctx context.Context) (*npool.Announcement, error) {
	info, err := cli.GetAnnouncement(ctx, *h.ID)
	if err != nil {
		return nil, err
	}
	if info == nil {
		return nil, fmt.Errorf("announcement not exist")
	}

	return info, nil
}
