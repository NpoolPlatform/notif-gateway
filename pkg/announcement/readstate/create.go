package readstate

import (
	"context"
	"fmt"

	usermwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/user"
	"github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"
	npool "github.com/NpoolPlatform/message/npool/notif/gw/v1/announcement/readstate"
	mwpb "github.com/NpoolPlatform/message/npool/notif/mw/v1/announcement/readstate"
	cli "github.com/NpoolPlatform/notif-middleware/pkg/client/announcement/readstate"
)

func (h *Handler) CreateReadState(ctx context.Context) (*npool.ReadState, error) {
	existUser, err := usermwcli.ExistUser(ctx, *h.AppID, *h.UserID)
	if err != nil {
		return nil, err
	}
	if !existUser {
		return nil, fmt.Errorf("invalid user")
	}

	exist, err := cli.ExistReadStateConds(ctx, &mwpb.Conds{
		AppID:          &basetypes.StringVal{Op: cruder.EQ, Value: *h.AppID},
		UserID:         &basetypes.StringVal{Op: cruder.EQ, Value: *h.UserID},
		AnnouncementID: &basetypes.StringVal{Op: cruder.EQ, Value: *h.AnnouncementID},
	})
	if err != nil {
		return nil, err
	}
	if exist {
		return nil, fmt.Errorf("read state exist")
	}

	info, err := cli.CreateReadState(
		ctx,
		&mwpb.ReadStateReq{
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

	return h.GetReadState(ctx)
}
