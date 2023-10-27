package contact

import (
	"context"
	"fmt"

	"github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"
	npool "github.com/NpoolPlatform/message/npool/notif/mw/v1/contact"
	cli "github.com/NpoolPlatform/notif-middleware/pkg/client/contact"
)

func (h *Handler) UpdateContact(ctx context.Context) (*npool.Contact, error) {
	exist, err := cli.ExistContactConds(ctx, &npool.Conds{
		ID:    &basetypes.Uint32Val{Op: cruder.EQ, Value: *h.ID},
		AppID: &basetypes.StringVal{Op: cruder.EQ, Value: *h.AppID},
	})
	if err != nil {
		return nil, err
	}
	if !exist {
		return nil, fmt.Errorf("contact not found")
	}

	info, err := cli.UpdateContact(ctx, &npool.ContactReq{
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
