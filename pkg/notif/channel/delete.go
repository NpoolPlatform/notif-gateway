package channel

import (
	"context"
	"fmt"

	npool "github.com/NpoolPlatform/message/npool/notif/mw/v1/notif/channel"
	cli "github.com/NpoolPlatform/notif-middleware/pkg/client/notif/channel"
)

func (h *Handler) DeleteChannel(ctx context.Context) (*npool.Channel, error) {
	info, err := h.GetChannel(ctx)
	if err != nil {
		return nil, err
	}
	if info == nil {
		return nil, nil
	}
	if info.AppID != *h.AppID {
		return nil, fmt.Errorf("permission denied")
	}

	_, err = cli.DeleteChannel(ctx, *h.ID)
	if err != nil {
		return nil, err
	}

	return info, nil
}
