package contact

import (
	"context"
	"fmt"

	npool "github.com/NpoolPlatform/message/npool/notif/mw/v1/contact"
	cli "github.com/NpoolPlatform/notif-middleware/pkg/client/contact"
)

func (h *Handler) UpdateContact(ctx context.Context) (*npool.Contact, error) {
	info, err := h.GetContact(ctx)
	if err != nil {
		return nil, err
	}

	if info == nil {
		return nil, fmt.Errorf("contact not found")
	}
	if info.AppID != *h.AppID {
		return nil, fmt.Errorf("permission denied")
	}

	_, err = cli.UpdateContact(ctx, &npool.ContactReq{
		ID:          h.ID,
		AppID:       h.AppID,
		Account:     h.Account,
		AccountType: h.AccountType,
		Sender:      h.Sender,
	},
	)
	if err != nil {
		return nil, err
	}
	h.EntID = &info.EntID

	return h.GetContact(ctx)
}
