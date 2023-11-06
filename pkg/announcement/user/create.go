package user

import (
	"context"
	"fmt"

	usermwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/user"
	"github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"
	npool "github.com/NpoolPlatform/message/npool/notif/gw/v1/announcement/user"
	mwpb "github.com/NpoolPlatform/message/npool/notif/mw/v1/announcement/user"
	cli "github.com/NpoolPlatform/notif-middleware/pkg/client/announcement/user"
)

func (h *Handler) CreateAnnouncementUser(ctx context.Context) (*npool.AnnouncementUser, error) {
	existUser, err := usermwcli.ExistUser(ctx, *h.AppID, *h.UserID)
	if err != nil {
		return nil, err
	}
	if !existUser {
		return nil, fmt.Errorf("invalid user")
	}

	exist, err := cli.ExistAnnouncementUserConds(
		ctx,
		&mwpb.Conds{
			AppID:          &basetypes.StringVal{Op: cruder.EQ, Value: *h.AppID},
			UserID:         &basetypes.StringVal{Op: cruder.EQ, Value: *h.UserID},
			AnnouncementID: &basetypes.StringVal{Op: cruder.EQ, Value: *h.AnnouncementID},
		},
	)
	if err != nil {
		return nil, err
	}
	if exist {
		return nil, fmt.Errorf("user exist")
	}

	info, err := cli.CreateAnnouncementUser(
		ctx,
		&mwpb.AnnouncementUserReq{
			AppID:          h.AppID,
			UserID:         h.UserID,
			AnnouncementID: h.AnnouncementID,
		},
	)
	if err != nil {
		return nil, err
	}

	h.ID = &info.ID
	h.EntID = &info.EntID
	return h.GetAnnouncementUser(ctx)
}
