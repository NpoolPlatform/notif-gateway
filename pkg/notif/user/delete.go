package user

import (
	"context"
	"fmt"

	"github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"
	npool "github.com/NpoolPlatform/message/npool/notif/gw/v1/notif/user"
	mwpb "github.com/NpoolPlatform/message/npool/notif/mw/v1/notif/user"
	mwcli "github.com/NpoolPlatform/notif-middleware/pkg/client/notif/user"
)

func (h *Handler) DeleteNotifUser(ctx context.Context) (*npool.NotifUser, error) {
	exist, err := mwcli.ExistNotifUserConds(ctx, &mwpb.Conds{
		ID:    &basetypes.Uint32Val{Op: cruder.EQ, Value: *h.ID},
		EntID: &basetypes.StringVal{Op: cruder.EQ, Value: *h.EntID},
		AppID: &basetypes.StringVal{Op: cruder.EQ, Value: *h.AppID},
	})
	if err != nil {
		return nil, err
	}
	if !exist {
		return nil, fmt.Errorf("notif user not exist")
	}

	info, err := mwcli.DeleteNotifUser(ctx, &mwpb.NotifUserReq{ID: h.ID})
	if err != nil {
		return nil, err
	}

	return h.GetNotifUserExt(ctx, info)
}
