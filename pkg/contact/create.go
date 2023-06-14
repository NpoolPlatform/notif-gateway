package contact

import (
	"context"
	"fmt"

	"github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"
	npool "github.com/NpoolPlatform/message/npool/notif/mw/v1/contact"
	cli "github.com/NpoolPlatform/notif-middleware/pkg/client/contact"
)

func (h *Handler) CreateContact(ctx context.Context) (*npool.Contact, error) {
	exist, err := cli.ExistContactConds(ctx, &npool.ExistContactCondsRequest{
		Conds: &npool.Conds{
			AppID: &basetypes.StringVal{
				Op:    cruder.EQ,
				Value: *h.AppID,
			},
			AccountType: &basetypes.Int32Val{
				Op:    cruder.EQ,
				Value: int32(*h.AccountType),
			},
			UsedFor: &basetypes.Int32Val{
				Op:    cruder.EQ,
				Value: int32(*h.UsedFor),
			},
		},
	})
	if err != nil {
		return nil, err
	}
	if exist {
		return nil, fmt.Errorf("contact exist")
	}

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
