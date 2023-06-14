package contact

import (
	"context"

	npool "github.com/NpoolPlatform/message/npool/notif/mw/v1/contact"
	cli "github.com/NpoolPlatform/notif-middleware/pkg/client/contact"
)

func (h *Handler) CreateContact(ctx context.Context) (*npool.Contact, error) {
	info, err := cli.CreateContact(
		ctx,
		&npool.ContactReq{
			AppID:       h.AppID,
			Account:     h.Account,
			AccountType: h.AccountType,
			UsedFor:     h.UsedFor,
			Sender:      h.Sender,
		},
	)
	if err != nil {
		return nil, err
	}

	h.ID = &info.ID
	return h.GetContact(ctx)
}
