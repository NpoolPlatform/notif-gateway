package contact

import (
	"context"

	npool "github.com/NpoolPlatform/message/npool/notif/mw/v1/contact"
	cli "github.com/NpoolPlatform/notif-middleware/pkg/client/contact"
)

func (h *Handler) DeleteContact(ctx context.Context) (*npool.Contact, error) {
	info, err := h.GetContact(ctx)
	if err != nil {
		return nil, err
	}

	_, err = cli.DeleteContact(ctx, *h.ID)
	if err != nil {
		return nil, err
	}

	return info, nil
}
