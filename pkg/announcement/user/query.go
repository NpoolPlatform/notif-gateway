package user

import (
	"context"
	"fmt"

	usermwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/user"
	"github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	usermwpb "github.com/NpoolPlatform/message/npool/appuser/mw/v1/user"
	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"
	npool "github.com/NpoolPlatform/message/npool/notif/gw/v1/announcement/user"
	mwpb "github.com/NpoolPlatform/message/npool/notif/mw/v1/announcement/user"
	mwcli "github.com/NpoolPlatform/notif-middleware/pkg/client/announcement/user"
)

func (h *Handler) GetAnnouncementUsers(ctx context.Context) ([]*npool.AnnouncementUser, uint32, error) {
	conds := &mwpb.Conds{AppID: &basetypes.StringVal{Op: cruder.EQ, Value: *h.AppID}}

	rows, total, err := mwcli.GetAnnouncementUsers(ctx, conds, h.Offset, h.Limit)
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

	infos := []*npool.AnnouncementUser{}
	for _, val := range rows {
		user, ok := userMap[val.UserID]
		if !ok {
			continue
		}
		infos = append(infos, &npool.AnnouncementUser{
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

func (h *Handler) GetAnnouncementUser(ctx context.Context) (*npool.AnnouncementUser, error) {
	row, err := mwcli.GetAnnouncementUser(ctx, *h.AppID, *h.EntID)
	if err != nil {
		return nil, err
	}
	if row == nil {
		return nil, fmt.Errorf("announcement not exist")
	}

	user, err := usermwcli.GetUser(ctx, row.AppID, row.UserID)
	if err != nil {
		return nil, err
	}

	info := &npool.AnnouncementUser{
		ID:               row.ID,
		EntID:            row.EntID,
		AnnouncementID:   row.AnnouncementID,
		AppID:            row.AppID,
		UserID:           row.UserID,
		EmailAddress:     user.EmailAddress,
		PhoneNO:          user.PhoneNO,
		Username:         user.Username,
		Title:            row.Title,
		Content:          row.Content,
		AnnouncementType: basetypes.NotifType(basetypes.NotifType_value[row.AnnouncementType]),
		Channel:          basetypes.NotifChannel(basetypes.NotifChannel_value[row.Channel]),
		CreatedAt:        row.CreatedAt,
		UpdatedAt:        row.UpdatedAt,
	}

	return info, nil
}

func (h *Handler) GetAnnouncementUserExt(ctx context.Context, row *mwpb.AnnouncementUser) (*npool.AnnouncementUser, error) {
	user, err := usermwcli.GetUser(ctx, row.AppID, row.UserID)
	if err != nil {
		return nil, err
	}

	info := &npool.AnnouncementUser{
		ID:               row.ID,
		EntID:            row.EntID,
		AnnouncementID:   row.AnnouncementID,
		AppID:            row.AppID,
		UserID:           row.UserID,
		EmailAddress:     user.EmailAddress,
		PhoneNO:          user.PhoneNO,
		Username:         user.Username,
		Title:            row.Title,
		Content:          row.Content,
		AnnouncementType: basetypes.NotifType(basetypes.NotifType_value[row.AnnouncementType]),
		Channel:          basetypes.NotifChannel(basetypes.NotifChannel_value[row.Channel]),
		CreatedAt:        row.CreatedAt,
		UpdatedAt:        row.UpdatedAt,
	}

	return info, nil
}
