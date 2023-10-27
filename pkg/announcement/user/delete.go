package user

import (
	"context"

	"github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"
	npool "github.com/NpoolPlatform/message/npool/notif/gw/v1/announcement/user"
	mwpb "github.com/NpoolPlatform/message/npool/notif/mw/v1/announcement/user"
	mwcli "github.com/NpoolPlatform/notif-middleware/pkg/client/announcement/user"
)

func (h *Handler) DeleteAnnouncementUser(ctx context.Context) (*npool.AnnouncementUser, error) {
	exist, err := mwcli.ExistAnnouncementUserConds(ctx, &mwpb.Conds{
		ID:    &basetypes.Uint32Val{Op: cruder.EQ, Value: *h.ID},
		AppID: &basetypes.StringVal{Op: cruder.EQ, Value: *h.AppID},
	})
	if err != nil {
		return nil, err
	}
	if !exist {
		return nil, nil
	}

	info, err := mwcli.DeleteAnnouncementUser(ctx, *h.ID)
	if err != nil {
		return nil, err
	}

	return h.GetAnnouncementUserExt(ctx, info)
}
