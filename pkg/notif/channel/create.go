package channel

import (
	"context"
	"fmt"

	"github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"
	npool "github.com/NpoolPlatform/message/npool/notif/mw/v1/notif/channel"
	cli "github.com/NpoolPlatform/notif-middleware/pkg/client/notif/channel"
)

func (h *Handler) CreateChannel(ctx context.Context) (*npool.Channel, error) {
	exist, err := cli.ExistChannelConds(ctx, &npool.Conds{
		AppID:     &basetypes.StringVal{Op: cruder.EQ, Value: *h.AppID},
		EventType: &basetypes.Uint32Val{Op: cruder.EQ, Value: uint32(*h.EventType)},
		Channel:   &basetypes.Uint32Val{Op: cruder.EQ, Value: uint32(*h.Channel)},
	})
	if err != nil {
		return nil, err
	}
	if exist {
		return nil, fmt.Errorf("channel exist")
	}

	info, err := cli.CreateChannel(ctx, &npool.ChannelReq{
		AppID:     h.AppID,
		EventType: h.EventType,
		Channel:   h.Channel,
	},
	)
	if err != nil {
		return nil, err
	}

	h.ID = &info.ID
	h.EntID = &info.EntID
	return h.GetChannel(ctx)
}
