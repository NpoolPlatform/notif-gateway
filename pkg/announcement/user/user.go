package user

import (
	"context"
	"fmt"

	appcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/app"
	usercli "github.com/NpoolPlatform/appuser-middleware/pkg/client/user"
	"github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	apppb "github.com/NpoolPlatform/message/npool/appuser/mw/v1/app"
	userpb "github.com/NpoolPlatform/message/npool/appuser/mw/v1/user"
	npool "github.com/NpoolPlatform/message/npool/notif/gw/v1/announcement/user"

	npoolpb "github.com/NpoolPlatform/message/npool"

	mgrpb "github.com/NpoolPlatform/message/npool/notif/mgr/v1/announcement/user"
	mgrcli "github.com/NpoolPlatform/notif-manager/pkg/client/announcement/user"
	mwcli "github.com/NpoolPlatform/notif-middleware/pkg/client/announcement/user"
)

func CreateAnnouncementUsers(
	ctx context.Context,
	appID string,
	userIDs []string,
	announcementID string,
) ([]*npool.AnnouncementUser, uint32, error) {
	req := []*mgrpb.UserReq{}
	for key := range userIDs {
		exist, err := mgrcli.ExistUserConds(ctx, &mgrpb.Conds{
			AppID: &npoolpb.StringVal{
				Op:    cruder.EQ,
				Value: appID,
			},
			UserID: &npoolpb.StringVal{
				Op:    cruder.EQ,
				Value: userIDs[key],
			},
			AnnouncementID: &npoolpb.StringVal{
				Op:    cruder.EQ,
				Value: announcementID,
			},
		})
		if err != nil {
			return nil, 0, err
		}
		if exist {
			return nil, 0, fmt.Errorf("already exist")
		}
		req = append(req, &mgrpb.UserReq{
			AppID:          &appID,
			UserID:         &userIDs[key],
			AnnouncementID: &announcementID,
		})
	}
	infos, err := mgrcli.CreateUsers(ctx, req)
	if err != nil {
		return nil, 0, err
	}

	ids := []string{}
	for _, val := range infos {
		ids = append(ids, val.ID)
	}
	return GetAnnouncementUsers(ctx, appID, &announcementID, ids, 0, uint32(len(ids)))
}

func DeleteAnnouncementUser(
	ctx context.Context,
	id string,
) (*npool.AnnouncementUser, error) {
	info, err := mgrcli.DeleteUser(ctx, id)
	if err != nil {
		return nil, err
	}

	return &npool.AnnouncementUser{
		ID:             info.ID,
		AnnouncementID: info.AnnouncementID,
		AppID:          info.AppID,
		UserID:         info.UserID,
		CreatedAt:      info.CreatedAt,
		UpdatedAt:      info.UpdatedAt,
	}, nil
}

//nolint:gocyclo
func GetAnnouncementUsers(
	ctx context.Context,
	appID string,
	announcementID *string,
	ids []string,
	offset, limit uint32,
) ([]*npool.AnnouncementUser, uint32, error) {
	if limit == 0 {
		limit = 100
	}
	conds := &mgrpb.Conds{
		AppID: &npoolpb.StringVal{
			Op:    cruder.EQ,
			Value: appID,
		},
	}
	if announcementID != nil {
		conds.AnnouncementID = &npoolpb.StringVal{
			Op:    cruder.EQ,
			Value: *announcementID,
		}
	}
	if len(ids) > 0 {
		conds.IDs = &npoolpb.StringSliceVal{
			Op:    cruder.IN,
			Value: ids,
		}
	}
	rows, total, err := mwcli.GetUsers(ctx, conds, int32(offset), int32(limit))
	if err != nil {
		return nil, 0, err
	}

	if len(rows) == 0 {
		return nil, total, nil
	}

	appIDs := []string{}
	userIDs := []string{}

	for _, val := range rows {
		appIDs = append(appIDs, val.AppID)
		if val.UserID != "" {
			userIDs = append(userIDs, val.UserID)
		}
	}
	appInfos, _, err := appcli.GetManyApps(ctx, appIDs)
	if err != nil {
		return nil, 0, err
	}
	appMap := map[string]*apppb.App{}
	for _, val := range appInfos {
		appMap[val.ID] = val
	}

	userMap := map[string]*userpb.User{}
	if len(userIDs) > 0 {
		userInfos, _, err := usercli.GetManyUsers(ctx, userIDs)
		if err != nil {
			return nil, 0, err
		}

		for _, val := range userInfos {
			userMap[val.ID] = val
		}
	}

	infos := []*npool.AnnouncementUser{}
	for _, val := range rows {
		app, ok := appMap[val.AppID]
		if !ok {
			continue
		}
		user, ok := userMap[val.UserID]
		if !ok {
			continue
		}
		infos = append(infos, &npool.AnnouncementUser{
			ID:               val.ID,
			AnnouncementID:   val.AnnouncementID,
			AppID:            val.AppID,
			AppName:          app.Name,
			UserID:           val.UserID,
			EmailAddress:     user.EmailAddress,
			PhoneNO:          user.PhoneNO,
			Username:         user.Username,
			Title:            val.Title,
			Content:          val.Content,
			AnnouncementType: val.AnnouncementType,
			CreatedAt:        val.CreatedAt,
			UpdatedAt:        val.UpdatedAt,
		})
	}
	return infos, total, nil
}
