package channel

import (
	"context"

	npool "github.com/NpoolPlatform/message/npool/notif/mw/v1/notif/channel"
	cli "github.com/NpoolPlatform/notif-middleware/pkg/client/notif/channel"
)

func (h *Handler) DeleteChannel(ctx context.Context) (*npool.Channel, error) {
	info, err := h.GetChannel(ctx)
	if err != nil {
		return nil, err
	}

	_, err = cli.DeleteChannel(ctx, *h.ID)
	if err != nil {
		return nil, err
	}

	return info, nil
}
