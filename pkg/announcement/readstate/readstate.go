package readstate

import (
	"context"
	"fmt"

	appcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/app"
	usercli "github.com/NpoolPlatform/appuser-middleware/pkg/client/user"
	"github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	apppb "github.com/NpoolPlatform/message/npool/appuser/mw/v1/app"
	userpb "github.com/NpoolPlatform/message/npool/appuser/mw/v1/user"
	npool "github.com/NpoolPlatform/message/npool/notif/gw/v1/announcement/readstate"

	npoolpb "github.com/NpoolPlatform/message/npool"

	mgrpb "github.com/NpoolPlatform/message/npool/notif/mgr/v1/announcement/readstate"
	mgrcli "github.com/NpoolPlatform/notif-manager/pkg/client/announcement/readstate"
	mwcli "github.com/NpoolPlatform/notif-middleware/pkg/client/announcement/readstate"
)

func CreateReadState(ctx context.Context, appID, userID, announcementID string) (*npool.ReadState, error) {
	exist, err := mgrcli.ExistReadStateConds(ctx, &mgrpb.Conds{
		AppID: &npoolpb.StringVal{
			Op:    cruder.EQ,
			Value: appID,
		},
		UserID: &npoolpb.StringVal{
			Op:    cruder.EQ,
			Value: userID,
		},
		AnnouncementID: &npoolpb.StringVal{
			Op:    cruder.EQ,
			Value: announcementID,
		},
	})
	if err != nil {
		return nil, err
	}
	if exist {
		return GetReadState(ctx, appID, userID, announcementID)
	}

	info, err := mgrcli.CreateReadState(ctx, &mgrpb.ReadStateReq{
		AppID:          &appID,
		UserID:         &userID,
		AnnouncementID: &announcementID,
	})
	if err != nil {
		return nil, err
	}

	return GetReadState(ctx, info.AppID, info.UserID, info.AnnouncementID)
}

func GetReadState(ctx context.Context, appID, userID, announcementID string) (*npool.ReadState, error) {
	info, err := mwcli.GetReadState(ctx, announcementID, userID)
	if err != nil {
		return nil, err
	}
	if info == nil {
		return nil, fmt.Errorf("announcement not exist")
	}

	if info.AppID != appID {
		return nil, fmt.Errorf("permission denied")
	}

	appInfo, err := appcli.GetApp(ctx, appID)
	if err != nil {
		return nil, err
	}

	if appInfo == nil {
		return nil, fmt.Errorf("user not exist")
	}

	userInfo, err := usercli.GetUser(ctx, appID, userID)
	if err != nil {
		return nil, err
	}

	if userInfo == nil {
		return nil, fmt.Errorf("user not exist")
	}

	return &npool.ReadState{
		AnnouncementID: info.AnnouncementID,
		AppID:          info.AppID,
		AppName:        appInfo.Name,
		UserID:         info.UserID,
		EmailAddress:   userInfo.EmailAddress,
		PhoneNO:        userInfo.PhoneNO,
		Username:       userInfo.Username,
		Title:          info.Title,
		Content:        info.Content,
		CreatedAt:      info.CreatedAt,
		UpdatedAt:      info.UpdatedAt,
	}, nil
}

func GetReadStates(ctx context.Context, appID string, userID *string, offset, limit uint32) ([]*npool.ReadState, uint32, error) {
	if limit == 0 {
		limit = 100
	}
	conds := &mgrpb.Conds{
		AppID: &npoolpb.StringVal{
			Op:    cruder.EQ,
			Value: appID,
		},
	}
	if userID != nil {
		conds.UserID = &npoolpb.StringVal{
			Op:    cruder.EQ,
			Value: *userID,
		}
	}
	rows, total, err := mwcli.GetReadStates(ctx, conds, int32(offset), int32(limit))
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

	infos := []*npool.ReadState{}
	for _, val := range rows {
		app, ok := appMap[val.AppID]
		if !ok {
			continue
		}
		user, ok := userMap[val.UserID]
		if !ok {
			continue
		}
		infos = append(infos, &npool.ReadState{
			AnnouncementID: val.AnnouncementID,
			AppID:          val.AppID,
			AppName:        app.Name,
			UserID:         val.UserID,
			EmailAddress:   user.EmailAddress,
			PhoneNO:        user.PhoneNO,
			Username:       user.Username,
			Title:          val.Title,
			Content:        val.Content,
			CreatedAt:      val.CreatedAt,
			UpdatedAt:      val.UpdatedAt,
		})
	}
	return infos, total, nil
}
