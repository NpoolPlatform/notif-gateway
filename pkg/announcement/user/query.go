package user

import (
	"context"
	"fmt"

	"github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"
	mwpb "github.com/NpoolPlatform/message/npool/notif/mw/v1/announcement/user"
	npool "github.com/NpoolPlatform/message/npool/notif/mw/v1/announcement/user"
	mwcli "github.com/NpoolPlatform/notif-middleware/pkg/client/announcement/user"
)

func (h *Handler) GetAnnouncementUsers(ctx context.Context) ([]*npool.AnnouncementUser, uint32, error) {
	conds := &mwpb.Conds{
		AppID: &basetypes.StringVal{
			Op:    cruder.EQ,
			Value: *h.AppID,
		},
	}
	if h.AnnouncementID != nil {
		conds.AnnouncementID = &basetypes.StringVal{
			Op:    cruder.EQ,
			Value: *h.AnnouncementID,
		}
	}
	if h.UserID != nil {
		conds.UserID = &basetypes.StringVal{
			Op:    cruder.EQ,
			Value: *h.UserID,
		}
	}

	infos, total, err := mwcli.GetAnnouncementUsers(ctx, conds, h.Offset, h.Limit)
	if err != nil {
		return nil, 0, err
	}
	if len(infos) == 0 {
		return nil, total, nil
	}

	return infos, total, nil
}

func (h *Handler) GetAnnouncementUser(ctx context.Context) (*npool.AnnouncementUser, error) {
	info, err := mwcli.GetAnnouncementUser(ctx, *h.ID)
	if err != nil {
		return nil, err
	}
	if info == nil {
		return nil, fmt.Errorf("announcement not exist")
	}

	return info, nil
}
