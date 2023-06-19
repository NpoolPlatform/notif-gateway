package readstate

import (
	"context"

	npool "github.com/NpoolPlatform/message/npool/notif/gw/v1/announcement/readstate"
	mwpb "github.com/NpoolPlatform/message/npool/notif/mw/v1/announcement/readstate"
	cli "github.com/NpoolPlatform/notif-middleware/pkg/client/announcement/readstate"
)

func (h *Handler) CreateReadState(ctx context.Context) (*npool.ReadState, error) {
	info, err := cli.CreateReadState(
		ctx,
		&mwpb.ReadStateReq{
			AppID:          h.AppID,
			UserID:         h.UserID,
			AnnouncementID: h.AnnouncementID,
		},
	)
	if err != nil {
		return nil, err
	}

	h.ID = &info.ID

	return h.GetReadState(ctx)
}
