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
	exist, err := cli.ExistChannelConds(ctx, &npool.ExistChannelCondsRequest{
		Conds: &npool.Conds{
			AppID: &basetypes.StringVal{
				Op:    cruder.EQ,
				Value: *h.AppID,
			},
			EventType: &basetypes.Uint32Val{
				Op:    cruder.EQ,
				Value: uint32(*h.EventType),
			},
			Channel: &basetypes.Uint32Val{
				Op:    cruder.EQ,
				Value: uint32(*h.Channel),
			},
		},
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
	return h.GetChannel(ctx)
}

func (h *Handler) CreateChannels(ctx context.Context) (infos []*npool.Channel, err error) {
	reqs := []*npool.ChannelReq{}
	for _, _type := range h.EventTypes {
		reqs = append(reqs, &npool.ChannelReq{
			AppID:     h.AppID,
			EventType: &_type,
			Channel:   h.Channel,
		})
	}
	infos, err = cli.CreateChannels(ctx, reqs)
	if err != nil {
		return nil, err
	}

	return infos, nil
}
