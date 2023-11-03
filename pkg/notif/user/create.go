package user

import (
	"context"
	"fmt"

	usermwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/user"
	npool "github.com/NpoolPlatform/message/npool/notif/gw/v1/notif/user"
	mwpb "github.com/NpoolPlatform/message/npool/notif/mw/v1/notif/user"
	cli "github.com/NpoolPlatform/notif-middleware/pkg/client/notif/user"
)

func (h *Handler) CreateNotifUser(ctx context.Context) (*npool.NotifUser, error) {
	existUser, err := usermwcli.ExistUser(ctx, *h.AppID, *h.UserID)
	if err != nil {
		return nil, err
	}
	if !existUser {
		return nil, fmt.Errorf("invalid user")
	}

	info, err := cli.CreateNotifUser(
		ctx,
		&mwpb.NotifUserReq{
			AppID:     h.AppID,
			UserID:    h.UserID,
			EventType: h.EventType,
		},
	)
	if err != nil {
		return nil, err
	}

	h.ID = &info.ID
	h.EntID = &info.EntID
	return h.GetNotifUser(ctx)
}
