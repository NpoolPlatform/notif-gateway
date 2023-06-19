package contact

import (
	"context"
	"fmt"

	"github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"
	npool "github.com/NpoolPlatform/message/npool/notif/mw/v1/contact"
	cli "github.com/NpoolPlatform/notif-middleware/pkg/client/contact"
)

type updateHandler struct {
	*Handler
}

func (h *Handler) UpdateContact(ctx context.Context) (*npool.Contact, error) {
	exist, err := cli.ExistContactConds(
		ctx,
		&npool.ExistContactCondsRequest{
			Conds: &npool.Conds{
				ID: &basetypes.StringVal{
					Op:    cruder.EQ,
					Value: *h.ID,
				},
				AppID: &basetypes.StringVal{
					Op:    cruder.EQ,
					Value: *h.AppID,
				},
			},
		},
	)
	if err != nil {
		return nil, err
	}
	if !exist {
		return nil, fmt.Errorf("invalid id or app id")
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

	handler := &updateHandler{
		Handler: h,
	}

	info, err := handler.GetContact(ctx)
	if err != nil {
		return nil, err
	}

	return info, nil
}
