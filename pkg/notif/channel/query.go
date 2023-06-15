package channel

import (
	"context"
	"fmt"

	"github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"
	npool "github.com/NpoolPlatform/message/npool/notif/mw/v1/notif/channel"
	cli "github.com/NpoolPlatform/notif-middleware/pkg/client/notif/channel"
)

func (h *Handler) GetChannels(ctx context.Context) ([]*npool.Channel, uint32, error) {
	infos, total, err := cli.GetChannels(
		ctx,
		&npool.Conds{
			AppID: &basetypes.StringVal{
				Op:    cruder.EQ,
				Value: *h.AppID,
			},
		},
		h.Offset,
		h.Limit,
	)
	if err != nil {
		return nil, 0, err
	}

	return infos, total, nil
}

func (h *Handler) GetChannel(ctx context.Context) (*npool.Channel, error) {
	if h.ID == nil {
		return nil, fmt.Errorf("invalid channel id")
	}

	info, err := cli.GetChannel(ctx, *h.ID)
	if err != nil {
		return nil, err
	}

	return info, nil
}
