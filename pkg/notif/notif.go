package notif

import (
	"context"
	"fmt"

	appmwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/app"
	usermwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/user"
	applangmwcli "github.com/NpoolPlatform/g11n-middleware/pkg/client/applang"

	appmwpb "github.com/NpoolPlatform/message/npool/appuser/mw/v1/app"
	usermwpb "github.com/NpoolPlatform/message/npool/appuser/mw/v1/user"
	applangmgrpb "github.com/NpoolPlatform/message/npool/g11n/mgr/v1/applang"
	applangmwpb "github.com/NpoolPlatform/message/npool/g11n/mw/v1/applang"

	mgrpb "github.com/NpoolPlatform/message/npool/notif/mgr/v1/notif"
	mgrcli "github.com/NpoolPlatform/notif-manager/pkg/client/notif"

	mwcli "github.com/NpoolPlatform/notif-middleware/pkg/client/notif"

	npool "github.com/NpoolPlatform/message/npool/notif/gw/v1/notif"

	cruder "github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	commonpb "github.com/NpoolPlatform/message/npool"
)

func GetNotif(ctx context.Context, id string) (*npool.Notif, error) {
	info, err := mgrcli.GetNotif(ctx, id)
	if err != nil {
		return nil, err
	}

	appInfo, err := appmwcli.GetApp(ctx, info.AppID)
	if err != nil {
		return nil, err
	}

	if appInfo == nil {
		return nil, fmt.Errorf("app %s not found", info.AppID)
	}

	userInfo, err := usermwcli.GetUser(ctx, info.AppID, info.UserID)
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
		Channel:      info.Channel,
		Notified:     info.Notified,
		CreatedAt:    info.CreatedAt,
		UpdatedAt:    info.UpdatedAt,
	}, nil
}

func UpdateNotifs(ctx context.Context, ids []string, notified bool) ([]*npool.Notif, error) {
	rows, err := mwcli.UpdateNotifs(ctx, ids, notified)
	if err != nil {
		return nil, err
	}

	if len(rows) == 0 {
		return nil, nil
	}

	appIDs := []string{}
	userIDs := []string{}

	for _, val := range rows {
		appIDs = append(appIDs, val.AppID)
		userIDs = append(userIDs, val.UserID)
	}
	appInfos, _, err := appmwcli.GetManyApps(ctx, appIDs)
	if err != nil {
		return nil, err
	}
	appMap := map[string]*appmwpb.App{}
	for _, val := range appInfos {
		appMap[val.ID] = val
	}

	userInfos, _, err := usermwcli.GetManyUsers(ctx, userIDs)
	if err != nil {
		return nil, err
	}
	userMap := map[string]*usermwpb.User{}
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
			Channel:      val.Channel,
			Notified:     val.Notified,
			CreatedAt:    val.CreatedAt,
			UpdatedAt:    val.UpdatedAt,
		})
	}
	return infos, nil
}

func GetNotifs(ctx context.Context, conds *mgrpb.Conds, offset, limit int32) ([]*npool.Notif, uint32, error) {
	rows, total, err := mgrcli.GetNotifs(ctx, conds, offset, limit)
	if err != nil {
		return nil, 0, err
	}

	if len(rows) == 0 {
		return nil, total, nil
	}

	appIDs := []string{}
	userIDs := []string{}
	langIDs := []string{}

	for _, val := range rows {
		appIDs = append(appIDs, val.AppID)
		userIDs = append(userIDs, val.UserID)
		langIDs = append(langIDs, val.LangID)
	}
	appInfos, _, err := appmwcli.GetManyApps(ctx, appIDs)
	if err != nil {
		return nil, 0, err
	}
	appMap := map[string]*appmwpb.App{}
	for _, val := range appInfos {
		appMap[val.ID] = val
	}

	userInfos, _, err := usermwcli.GetManyUsers(ctx, userIDs)
	if err != nil {
		return nil, 0, err
	}
	userMap := map[string]*usermwpb.User{}
	for _, val := range userInfos {
		userMap[val.ID] = val
	}

	langs, _, err := applangmwcli.GetLangs(ctx, &applangmgrpb.Conds{
		AppID:   conds.AppID,
		LangIDs: &commonpb.StringSliceVal{Op: cruder.IN, Value: langIDs},
	}, 0, int32(len(langIDs)))
	if err != nil {
		return nil, 0, err
	}

	langMap := map[string]*applangmwpb.Lang{}
	for _, lang := range langs {
		langMap[lang.LangID] = lang
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
		lang, ok := langMap[val.LangID]
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
			Channel:      val.Channel,
			LangID:       lang.LangID,
			Lang:         lang.Lang,
			Notified:     val.Notified,
			CreatedAt:    val.CreatedAt,
			UpdatedAt:    val.UpdatedAt,
		})
	}
	return infos, total, nil
}
