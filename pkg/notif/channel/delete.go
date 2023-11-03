package channel

import (
	"context"

	"github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"
	mwpb "github.com/NpoolPlatform/message/npool/notif/mw/v1/notif/channel"
	mwcli "github.com/NpoolPlatform/notif-middleware/pkg/client/notif/channel"
)

func (h *Handler) DeleteChannel(ctx context.Context) (*mwpb.Channel, error) {
	exist, err := mwcli.ExistChannelConds(ctx, &mwpb.Conds{
		ID:    &basetypes.Uint32Val{Op: cruder.EQ, Value: *h.ID},
		EntID: &basetypes.StringVal{Op: cruder.EQ, Value: *h.EntID},
		AppID: &basetypes.StringVal{Op: cruder.EQ, Value: *h.AppID},
	})
	if err != nil {
		return nil, err
	}
	if !exist {
		return nil, nil
	}

	return mwcli.DeleteChannel(ctx, *h.ID)
}
