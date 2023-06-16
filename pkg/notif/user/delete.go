package user

import (
	"context"

	npool "github.com/NpoolPlatform/message/npool/notif/gw/v1/notif/user"
	notifusermw "github.com/NpoolPlatform/message/npool/notif/mw/v1/notif/user"
	cli "github.com/NpoolPlatform/notif-middleware/pkg/client/notif/user"
)

func (h *Handler) DeleteNotifUser(ctx context.Context) (*npool.NotifUser, error) {
	info, err := h.GetNotifUser(ctx)
	if err != nil {
		return nil, err
	}

	_, err = cli.DeleteUser(ctx, &notifusermw.UserNotifReq{
		ID: h.ID,
	})
	if err != nil {
		return nil, err
	}

	return info, nil
}
