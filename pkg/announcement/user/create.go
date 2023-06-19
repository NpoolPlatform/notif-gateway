package user

import (
	"context"
	"fmt"

	"github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"
	npool "github.com/NpoolPlatform/message/npool/notif/gw/v1/announcement/user"
	mwpb "github.com/NpoolPlatform/message/npool/notif/mw/v1/announcement/user"
	cli "github.com/NpoolPlatform/notif-middleware/pkg/client/announcement/user"
)

type createHandler struct {
	*Handler
}

func (h *createHandler) validate() error {
	if h.AppID == nil {
		return fmt.Errorf("app id is empty")
	}
	if h.UserID == nil {
		return fmt.Errorf("user id is empty")
	}
	if h.AnnouncementID == nil {
		return fmt.Errorf("announcement id is empty")
	}
	return nil
}

func (h *Handler) CreateAnnouncementUser(ctx context.Context) (*npool.AnnouncementUser, error) {
	handler := &createHandler{
		Handler: h,
	}

	if err := handler.validate(); err != nil {
		return nil, err
	}

	exist, err := cli.ExistAnnouncementUserConds(
		ctx,
		&mwpb.ExistAnnouncementUserCondsRequest{
			Conds: &mwpb.Conds{
				AppID: &basetypes.StringVal{
					Op:    cruder.EQ,
					Value: *h.AppID,
				},
				UserID: &basetypes.StringVal{
					Op:    cruder.EQ,
					Value: *h.UserID,
				},
				AnnouncementID: &basetypes.StringVal{
					Op:    cruder.EQ,
					Value: *h.AnnouncementID,
				},
			},
		},
	)
	if err != nil {
		return nil, err
	}
	if exist {
		return nil, fmt.Errorf("user exist")
	}

	info, err := cli.CreateAnnouncementUser(
		ctx,
		&mwpb.AnnouncementUserReq{
			AppID:          h.AppID,
			UserID:         h.UserID,
			AnnouncementID: h.AnnouncementID,
		},
	)
	if err != nil {
		return nil, err
	}

	h.ID = &info.ID
	return h.GetAnnouncementUser(ctx)
}
