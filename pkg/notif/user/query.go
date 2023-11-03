package user

import (
	"context"
	"fmt"

	usermwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/user"
	"github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	usermwpb "github.com/NpoolPlatform/message/npool/appuser/mw/v1/user"
	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"
	npool "github.com/NpoolPlatform/message/npool/notif/gw/v1/notif/user"
	mwpb "github.com/NpoolPlatform/message/npool/notif/mw/v1/notif/user"
	mwcli "github.com/NpoolPlatform/notif-middleware/pkg/client/notif/user"
)

func (h *Handler) GetNotifUsers(ctx context.Context) ([]*npool.NotifUser, uint32, error) {
	conds := &mwpb.Conds{
		AppID: &basetypes.StringVal{Op: cruder.EQ, Value: *h.AppID},
	}
	if h.EventType != nil {
		conds.EventType = &basetypes.Uint32Val{Op: cruder.EQ, Value: uint32(*h.EventType)}
	}
	if h.UserID != nil {
		conds.UserID = &basetypes.StringVal{Op: cruder.EQ, Value: *h.UserID}
	}

	rows, total, err := mwcli.GetNotifUsers(ctx, conds, h.Offset, h.Limit)
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
			IDs: &basetypes.StringSliceVal{
				Op: cruder.IN, Value: userIDs,
			},
		}, 0, int32(len(userIDs)))
		if err != nil {
			return nil, 0, err
		}

		for _, val := range userInfos {
			userMap[val.ID] = val
		}
	}

	infos := []*npool.NotifUser{}
	for _, val := range rows {
		user, ok := userMap[val.UserID]
		if !ok {
			continue
		}
		infos = append(infos, &npool.NotifUser{
			ID:           val.ID,
			EntID:        val.EntID,
			EventType:    val.EventType,
			AppID:        val.AppID,
			UserID:       val.UserID,
			EmailAddress: user.EmailAddress,
			PhoneNO:      user.PhoneNO,
			Username:     user.Username,
			// Title:            val.Title, // TODO
			// Content:          val.Content,
			// NotifType: val.NotifType,
			CreatedAt: val.CreatedAt,
			UpdatedAt: val.UpdatedAt,
		})
	}
	return infos, total, nil
}

func (h *Handler) GetNotifUser(ctx context.Context) (*npool.NotifUser, error) {
	row, err := mwcli.GetNotifUser(ctx, *h.EntID)
	if err != nil {
		return nil, err
	}
	if row == nil {
		return nil, fmt.Errorf("notif not exist")
	}

	user, err := usermwcli.GetUser(ctx, row.AppID, row.UserID)
	if err != nil {
		return nil, err
	}

	info := &npool.NotifUser{
		ID:           row.ID,
		EntID:        row.EntID,
		EventType:    row.EventType,
		AppID:        row.AppID,
		UserID:       row.UserID,
		EmailAddress: user.EmailAddress,
		PhoneNO:      user.PhoneNO,
		Username:     user.Username,
		CreatedAt:    row.CreatedAt,
		UpdatedAt:    row.UpdatedAt,
	}

	return info, nil
}

func (h *Handler) GetNotifUserExt(ctx context.Context, row *mwpb.NotifUser) (*npool.NotifUser, error) {
	user, err := usermwcli.GetUser(ctx, row.AppID, row.UserID)
	if err != nil {
		return nil, err
	}

	info := &npool.NotifUser{
		ID:           row.ID,
		EntID:        row.EntID,
		EventType:    row.EventType,
		AppID:        row.AppID,
		UserID:       row.UserID,
		EmailAddress: user.EmailAddress,
		PhoneNO:      user.PhoneNO,
		Username:     user.Username,
		CreatedAt:    row.CreatedAt,
		UpdatedAt:    row.UpdatedAt,
	}

	return info, nil
}
