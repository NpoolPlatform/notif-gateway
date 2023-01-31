package notif

import (
	"context"
	"fmt"

	appcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/app"
	usercli "github.com/NpoolPlatform/appuser-middleware/pkg/client/user"
	"github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	npoolpb "github.com/NpoolPlatform/message/npool"
	apppb "github.com/NpoolPlatform/message/npool/appuser/mw/v1/app"
	userpb "github.com/NpoolPlatform/message/npool/appuser/mw/v1/user"

	mgrpb "github.com/NpoolPlatform/message/npool/notif/mgr/v1/notif"
	mgrcli "github.com/NpoolPlatform/notif-manager/pkg/client/notif"

	npool "github.com/NpoolPlatform/message/npool/notif/gw/v1/notif"
)

func GetNotif(ctx context.Context, id string) (*npool.Notif, error) {
	info, err := mgrcli.GetNotif(ctx, id)
	if err != nil {
		return nil, err
	}

	appInfo, err := appcli.GetApp(ctx, info.AppID)
	if err != nil {
		return nil, err
	}

	if appInfo == nil {
		return nil, fmt.Errorf("app %s not found", info.AppID)
	}

	userInfo, err := usercli.GetUser(ctx, info.AppID, info.UserID)
	if err != nil {
		return nil, err
	}

	if userInfo == nil {
		return nil, fmt.Errorf("user %s not found", info.UserID)
	}

	return &npool.Notif{
		ID:           info.ID,
		AppID:        info.AppID,
		AppName:      appInfo.Name,
		UserID:       info.UserID,
		EmailAddress: userInfo.EmailAddress,
		PhoneNO:      userInfo.PhoneNO,
		Username:     userInfo.Username,
		EventType:    info.EventType,
		UseTemplate:  info.UseTemplate,
		Title:        info.Title,
		Content:      info.Content,
		Channels:     info.Channels,
		AlreadyRead:  info.AlreadyRead,
	}, nil
}

func GetNotifs(ctx context.Context, appID, userID string, offset, limit uint32) ([]*npool.Notif, uint32, error) {
	rows, total, err := mgrcli.GetNotifs(ctx, &mgrpb.Conds{
		AppID: &npoolpb.StringVal{
			Op:    cruder.EQ,
			Value: appID,
		},
		UserID: &npoolpb.StringVal{
			Op:    cruder.EQ,
			Value: userID,
		},
	}, int32(offset), int32(limit))
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
		userIDs = append(userIDs, val.UserID)
	}
	appInfos, _, err := appcli.GetManyApps(ctx, appIDs)
	if err != nil {
		return nil, 0, err
	}
	appMap := map[string]*apppb.App{}
	for _, val := range appInfos {
		appMap[val.ID] = val
	}

	userInfos, _, err := usercli.GetManyUsers(ctx, userIDs)
	if err != nil {
		return nil, 0, err
	}
	userMap := map[string]*userpb.User{}
	for _, val := range userInfos {
		userMap[val.ID] = val
	}

	infos := []*npool.Notif{}
	for _, val := range rows {
		app, ok := appMap[val.AppID]
		if !ok {
			continue
		}
		user, ok := userMap[val.UserID]
		if !ok {
			continue
		}

		infos = append(infos, &npool.Notif{
			ID:           val.ID,
			AppID:        val.AppID,
			AppName:      app.Name,
			UserID:       val.UserID,
			EmailAddress: user.EmailAddress,
			PhoneNO:      user.PhoneNO,
			Username:     user.Username,
			EventType:    val.EventType,
			UseTemplate:  val.UseTemplate,
			Title:        val.Title,
			Content:      val.Content,
			Channels:     val.Channels,
			AlreadyRead:  val.AlreadyRead,
		})
	}
	return infos, total, nil
}
