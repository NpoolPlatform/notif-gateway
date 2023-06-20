package notif

import (
	"context"
	"fmt"

	appmwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/app"
	usermwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/user"
	applangmwcli "github.com/NpoolPlatform/g11n-middleware/pkg/client/applang"
	cruder "github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	appmwpb "github.com/NpoolPlatform/message/npool/appuser/mw/v1/app"
	usermwpb "github.com/NpoolPlatform/message/npool/appuser/mw/v1/user"
	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"
	applangmwpb "github.com/NpoolPlatform/message/npool/g11n/mw/v1/applang"
	npool "github.com/NpoolPlatform/message/npool/notif/gw/v1/notif"
	notifmwpb "github.com/NpoolPlatform/message/npool/notif/mw/v1/notif"
	notifmwcli "github.com/NpoolPlatform/notif-middleware/pkg/client/notif"
)

func (h *Handler) GetNotif(ctx context.Context) (*npool.Notif, error) {
	info, err := notifmwcli.GetNotifOnly(ctx, &notifmwpb.Conds{
		ID: &basetypes.StringVal{Op: cruder.EQ, Value: *h.ID},
	})
	if err != nil {
		return nil, err
	}
	if info == nil {
		return nil, nil
	}

	app, err := appmwcli.GetApp(ctx, info.AppID)
	if err != nil {
		return nil, err
	}
	if app == nil {
		return nil, fmt.Errorf("app %s not found", info.AppID)
	}
	user, err := usermwcli.GetUser(ctx, info.AppID, info.UserID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, fmt.Errorf("user %s not found", info.UserID)
	}
	return &npool.Notif{
		ID:           info.ID,
		AppID:        info.AppID,
		AppName:      app.Name,
		UserID:       info.UserID,
		EmailAddress: user.EmailAddress,
		PhoneNO:      user.PhoneNO,
		Username:     user.Username,
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

func (h *Handler) GetNotifs(ctx context.Context) ([]*npool.Notif, uint32, error) {
	conds := &notifmwpb.Conds{}
	if h.AppID != nil {
		conds.AppID = &basetypes.StringVal{Op: cruder.EQ, Value: *h.AppID}
	}
	if h.UserID != nil {
		conds.UserID = &basetypes.StringVal{Op: cruder.EQ, Value: *h.UserID}
	}
	if h.LangID != nil {
		conds.LangID = &basetypes.StringVal{Op: cruder.EQ, Value: *h.LangID}
	}
	if h.Channel != nil {
		conds.Channel = &basetypes.Uint32Val{Op: cruder.EQ, Value: uint32(*h.Channel)}
	}
	if len(h.IDs) > 0 {
		conds.IDs = &basetypes.StringSliceVal{Op: cruder.IN, Value: h.IDs}
	}
	rows, total, err := notifmwcli.GetNotifs(ctx, conds, h.Offset, h.Limit)
	if err != nil {
		return nil, 0, err
	}

	if len(rows) == 0 {
		return nil, 0, nil
	}

	appIDs := []string{}
	userIDs := []string{}
	langIDs := []string{}

	for _, val := range rows {
		appIDs = append(appIDs, val.AppID)
		userIDs = append(userIDs, val.UserID)
		langIDs = append(langIDs, val.LangID)
	}
	appInfos, _, err := appmwcli.GetApps(ctx, &appmwpb.Conds{
		IDs: &basetypes.StringSliceVal{Op: cruder.IN, Value: appIDs},
	}, 0, int32(len(appIDs)))
	if err != nil {
		return nil, 0, err
	}
	appMap := map[string]*appmwpb.App{}
	for _, val := range appInfos {
		appMap[val.ID] = val
	}

	userInfos, _, err := usermwcli.GetUsers(ctx, &usermwpb.Conds{
		IDs: &basetypes.StringSliceVal{Op: cruder.IN, Value: userIDs},
	}, 0, int32(len(userIDs)))
	if err != nil {
		return nil, 0, err
	}
	userMap := map[string]*usermwpb.User{}
	for _, val := range userInfos {
		userMap[val.ID] = val
	}

	langs, _, err := applangmwcli.GetLangs(ctx, &applangmwpb.Conds{
		AppID:   conds.AppID,
		LangIDs: &basetypes.StringSliceVal{Op: cruder.IN, Value: langIDs},
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
