package sendstate

import (
	"context"

	appcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/app"
	usercli "github.com/NpoolPlatform/appuser-middleware/pkg/client/user"
	"github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	apppb "github.com/NpoolPlatform/message/npool/appuser/mw/v1/app"
	userpb "github.com/NpoolPlatform/message/npool/appuser/mw/v1/user"
	npool "github.com/NpoolPlatform/message/npool/notif/gw/v1/announcement/sendstate"
	channelpb "github.com/NpoolPlatform/message/npool/notif/mgr/v1/channel"

	npoolpb "github.com/NpoolPlatform/message/npool"

	mgrpb "github.com/NpoolPlatform/message/npool/notif/mgr/v1/announcement/sendstate"
	mwcli "github.com/NpoolPlatform/notif-middleware/pkg/client/announcement/sendstate"
)

func GetSendStates(
	ctx context.Context,
	appID string,
	userID *string,
	offset, limit uint32,
	channel *channelpb.NotifChannel,
) (
	[]*npool.SendState,
	uint32,
	error,
) {
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
	if channel != nil {
		conds.Channel = &npoolpb.Uint32Val{
			Op:    cruder.EQ,
			Value: uint32(*channel),
		}
	}
	rows, total, err := mwcli.GetSendStates(ctx, conds, int32(offset), int32(limit))
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
	userInfos, _, err := usercli.GetManyUsers(ctx, userIDs)
	if err != nil {
		return nil, 0, err
	}
	userMap := map[string]*userpb.User{}
	for _, val := range userInfos {
		userMap[val.ID] = val
	}

	infos := []*npool.SendState{}
	for _, val := range rows {
		app, ok := appMap[val.AppID]
		if !ok {
			continue
		}
		user, ok := userMap[val.UserID]
		if !ok {
			continue
		}
		infos = append(infos, &npool.SendState{
			AnnouncementID: val.AnnouncementID,
			AppID:          val.AppID,
			AppName:        app.Name,
			UserID:         val.UserID,
			EmailAddress:   user.EmailAddress,
			PhoneNO:        user.PhoneNO,
			Username:       user.Username,
			Title:          val.Title,
			Content:        val.Content,
			Channel:        val.Channel,
			AlreadySend:    val.AlreadySend,
		})
	}
	return infos, total, nil
}
