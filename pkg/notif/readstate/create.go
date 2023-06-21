package readstate

import (
	"context"
	"fmt"

	"github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"
	npool "github.com/NpoolPlatform/message/npool/notif/gw/v1/notif/readstate"
	mwpb "github.com/NpoolPlatform/message/npool/notif/mw/v1/notif/readstate"
	cli "github.com/NpoolPlatform/notif-middleware/pkg/client/notif/readstate"
)

func (h *Handler) CreateReadState(ctx context.Context) (*npool.ReadState, error) {
	infos, _, err := cli.GetReadStates(ctx, &mwpb.Conds{
		AppID: &basetypes.StringVal{
			Op:    cruder.EQ,
			Value: *h.AppID,
		},
		UserID: &basetypes.StringVal{
			Op:    cruder.EQ,
			Value: *h.UserID,
		},
		NotifID: &basetypes.StringVal{
			Op:    cruder.EQ,
			Value: *h.NotifID,
		},
	}, 0, 1)
	if err != nil {
		return nil, err
	}
	if len(infos) > 0 {
		return nil, fmt.Errorf("read state exist")
	}

	info, err := cli.CreateReadState(
		ctx,
		&mwpb.ReadStateReq{
			AppID:   h.AppID,
			UserID:  h.UserID,
			NotifID: h.NotifID,
		},
	)
	if err != nil {
		return nil, err
	}

	h.ID = &info.ID

	return h.GetReadState(ctx)
}
