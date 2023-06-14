package contact

import (
	"context"
	"fmt"

	"github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"
	npool "github.com/NpoolPlatform/message/npool/notif/mw/v1/contact"
	cli "github.com/NpoolPlatform/notif-middleware/pkg/client/contact"
)

func (h *Handler) GetContacts(ctx context.Context) ([]*npool.Contact, uint32, error) {
	infos, total, err := cli.GetContacts(
		ctx,
		&npool.Conds{
			AppID: &basetypes.StringVal{
				Op:    cruder.EQ,
				Value: h.AppID.String(),
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

func (h *Handler) GetContact(ctx context.Context) (*npool.Contact, error) {
	if h.ID == nil {
		return nil, fmt.Errorf("invalid contact id")
	}

	info, err := cli.GetContact(ctx, h.ID.String())
	if err != nil {
		return nil, err
	}

	return info, nil
}
