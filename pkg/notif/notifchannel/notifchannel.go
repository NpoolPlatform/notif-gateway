package notifchannel

import (
	"context"
	"fmt"

	"github.com/NpoolPlatform/message/npool/notif/mgr/v1/channel"
	"github.com/NpoolPlatform/message/npool/third/mgr/v1/usedfor"

	appcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/app"
	"github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	npoolpb "github.com/NpoolPlatform/message/npool"
	apppb "github.com/NpoolPlatform/message/npool/appuser/mw/v1/app"

	mgrpb "github.com/NpoolPlatform/message/npool/notif/mgr/v1/notif/notifchannel"
	mgrcli "github.com/NpoolPlatform/notif-manager/pkg/client/notif/notifchannel"

	npool "github.com/NpoolPlatform/message/npool/notif/gw/v1/notif/notifchannel"
)

func DeleteNotifChannel(ctx context.Context, id string) (*npool.NotifChannel, error) {
	info, err := mgrcli.DeleteNotifChannel(ctx, id)
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

	return &npool.NotifChannel{
		ID:        info.ID,
		AppID:     info.AppID,
		AppName:   appInfo.Name,
		EventType: info.EventType,
		Channel:   info.Channel,
		CreatedAt: info.CreatedAt,
	}, nil
}

func CreateNotifChannels(
	ctx context.Context,
	appID string,
	eventTypes []usedfor.UsedFor,
	channel1 channel.NotifChannel,
) (
	[]*npool.NotifChannel,
	error,
) {
	var req []*mgrpb.NotifChannelReq
	for key := range eventTypes {
		req = append(req, &mgrpb.NotifChannelReq{
			AppID:     &appID,
			EventType: &eventTypes[key],
			Channel:   &channel1,
		})
	}
	rows, err := mgrcli.CreateNotifChannels(ctx, req)
	if err != nil {
		return nil, err
	}

	if len(rows) == 0 {
		return nil, nil
	}

	infos, err := extends(ctx, rows)
	if err != nil {
		return nil, err
	}
	return infos, nil
}

func GetNotifChannels(
	ctx context.Context,
	appID string,
	offset, limit uint32,
) (
	[]*npool.NotifChannel,
	uint32,
	error,
) {
	rows, total, err := mgrcli.GetNotifChannels(ctx, &mgrpb.Conds{
		AppID: &npoolpb.StringVal{
			Op:    cruder.EQ,
			Value: appID,
		},
	}, int32(offset), int32(limit))
	if err != nil {
		return nil, 0, err
	}

	if len(rows) == 0 {
		return nil, total, nil
	}

	infos, err := extends(ctx, rows)
	if err != nil {
		return nil, 0, err
	}
	return infos, total, nil
}

func extends(
	ctx context.Context,
	rows []*mgrpb.NotifChannel,
) (
	[]*npool.NotifChannel,
	error,
) {
	appIDs := []string{}

	for _, val := range rows {
		appIDs = append(appIDs, val.AppID)
	}
	appInfos, _, err := appcli.GetManyApps(ctx, appIDs)
	if err != nil {
		return nil, err
	}
	appMap := map[string]*apppb.App{}
	for _, val := range appInfos {
		appMap[val.ID] = val
	}

	infos := []*npool.NotifChannel{}
	for _, val := range rows {
		app, ok := appMap[val.AppID]
		if !ok {
			continue
		}

		infos = append(infos, &npool.NotifChannel{
			ID:        val.ID,
			AppID:     val.AppID,
			AppName:   app.Name,
			EventType: val.EventType,
			Channel:   val.Channel,
			CreatedAt: val.CreatedAt,
		})
	}
	return infos, nil
}