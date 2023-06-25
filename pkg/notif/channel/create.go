package channel

import (
	"context"

	"github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"
	npool "github.com/NpoolPlatform/message/npool/notif/mw/v1/notif/channel"
	cli "github.com/NpoolPlatform/notif-middleware/pkg/client/notif/channel"
)

func (h *Handler) CreateChannel(ctx context.Context) ([]*npool.Channel, error) {
	infos := []*npool.Channel{}
	for _, eventType := range h.EventTypes {
		exist, err := cli.ExistChannelConds(ctx, &npool.ExistChannelCondsRequest{
			Conds: &npool.Conds{
				AppID: &basetypes.StringVal{
					Op:    cruder.EQ,
					Value: *h.AppID,
				},
				EventType: &basetypes.Uint32Val{
					Op:    cruder.EQ,
					Value: uint32(eventType),
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
			continue
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

		infos = append(infos, info)
	}

	return infos, nil
}
