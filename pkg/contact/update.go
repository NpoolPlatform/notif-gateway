package contact

import (
	"context"

	npool "github.com/NpoolPlatform/message/npool/notif/mw/v1/contact"
	cli "github.com/NpoolPlatform/notif-middleware/pkg/client/contact"
)

type updateHandler struct {
	*Handler
}

func (h *Handler) UpdateContact(ctx context.Context) (*npool.Contact, error) {
	_, err := cli.UpdateContact(ctx, &npool.ContactReq{
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

	handler := &updateHandler{
		Handler: h,
	}

	info, err := handler.GetContact(ctx)
	if err != nil {
		return nil, err
	}

	return info, nil
}
