package readstate

import (
	"context"
	"fmt"

	usermwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/user"
	"github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	usermwpb "github.com/NpoolPlatform/message/npool/appuser/mw/v1/user"
	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"
	npool "github.com/NpoolPlatform/message/npool/notif/gw/v1/announcement/readstate"
	mwpb "github.com/NpoolPlatform/message/npool/notif/mw/v1/announcement/readstate"
	mwcli "github.com/NpoolPlatform/notif-middleware/pkg/client/announcement/readstate"
)

func (h *Handler) GetReadStates(ctx context.Context) ([]*npool.ReadState, uint32, error) {
	if h.UserID != nil {
		exist, err := usermwcli.ExistUser(ctx, *h.AppID, *h.UserID)
		if err != nil {
			return nil, 0, err
		}
		if !exist {
			return nil, 0, fmt.Errorf("invalid user")
		}
	}

	conds := &mwpb.Conds{
		AppID: &basetypes.StringVal{Op: cruder.EQ, Value: *h.AppID},
	}
	if h.UserID != nil {
		conds.UserID = &basetypes.StringVal{Op: cruder.EQ, Value: *h.UserID}
	}
	if h.AnnouncementID != nil {
		conds.AnnouncementID = &basetypes.StringVal{Op: cruder.EQ, Value: *h.AnnouncementID}
	}

	rows, total, err := mwcli.GetReadStates(ctx, conds, h.Offset, h.Limit)
	if err != nil {
		return nil, 0, err
	}
	if len(rows) == 0 {
		return nil, total, nil
	}

	userIDs := []string{}

	for _, val := range rows {
		if val.UserID != "" {
			userIDs = append(userIDs, val.UserID)
		}
	}

	userMap := map[string]*usermwpb.User{}
	if len(userIDs) > 0 {
		userInfos, _, err := usermwcli.GetUsers(ctx, &usermwpb.Conds{
			EntIDs: &basetypes.StringSliceVal{
				Op: cruder.IN, Value: userIDs,
			},
		}, 0, int32(len(userIDs)))
		if err != nil {
			return nil, 0, err
		}

		for _, val := range userInfos {
			userMap[val.EntID] = val
		}
	}

	infos := []*npool.ReadState{}
	for _, val := range rows {
		user, ok := userMap[val.UserID]
		if !ok {
			continue
		}
		infos = append(infos, &npool.ReadState{
			ID:               val.ID,
			EntID:            val.EntID,
			AnnouncementID:   val.AnnouncementID,
			AppID:            val.AppID,
			UserID:           val.UserID,
			EmailAddress:     user.EmailAddress,
			PhoneNO:          user.PhoneNO,
			Username:         user.Username,
			Title:            val.Title,
			Content:          val.Content,
			AnnouncementType: basetypes.NotifType(basetypes.NotifType_value[val.AnnouncementType]),
			Channel:          basetypes.NotifChannel(basetypes.NotifChannel_value[val.Channel]),
			CreatedAt:        val.CreatedAt,
			UpdatedAt:        val.UpdatedAt,
		})
	}
	return infos, total, nil
}

func (h *Handler) GetReadState(ctx context.Context) (*npool.ReadState, error) {
	h.Offset = 0
	h.Limit = 1
	infos, _, err := h.GetReadStates(ctx)
	if err != nil {
		return nil, err
	}

	if len(infos) == 0 {
		return nil, fmt.Errorf("read state not found")
	}

	return infos[0], nil
}
