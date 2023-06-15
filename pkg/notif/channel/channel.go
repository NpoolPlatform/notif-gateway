package channel

// import (
// 	"context"
// 	"fmt"

// 	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"
// 	"github.com/NpoolPlatform/message/npool/notif/mgr/v1/channel"

// 	appmwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/app"
// 	appmwpb "github.com/NpoolPlatform/message/npool/appuser/mw/v1/app"

// 	mgrpb "github.com/NpoolPlatform/message/npool/notif/mgr/v1/notif/channel"
// 	mgrcli "github.com/NpoolPlatform/notif-manager/pkg/client/notif/channel"

// 	npool "github.com/NpoolPlatform/message/npool/notif/gw/v1/notif/channel"

// 	cruder "github.com/NpoolPlatform/libent-cruder/pkg/cruder"
// 	commonpb "github.com/NpoolPlatform/message/npool"
// )

// func DeleteChannel(ctx context.Context, id string) (*npool.Channel, error) {
// 	info, err := mgrcli.DeleteChannel(ctx, id)
// 	if err != nil {
// 		return nil, err
// 	}

// 	appInfo, err := appmwcli.GetApp(ctx, info.AppID)
// 	if err != nil {
// 		return nil, err
// 	}

// 	if appInfo == nil {
// 		return nil, fmt.Errorf("app %s not found", info.AppID)
// 	}

// 	return &npool.Channel{
// 		ID:        info.ID,
// 		AppID:     info.AppID,
// 		AppName:   appInfo.Name,
// 		EventType: info.EventType,
// 		Channel:   info.Channel,
// 		CreatedAt: info.CreatedAt,
// 	}, nil
// }

// func CreateChannels(
// 	ctx context.Context,
// 	appID string,
// 	eventTypes []basetypes.UsedFor,
// 	channel1 channel.NotifChannel,
// ) (
// 	[]*npool.Channel,
// 	error,
// ) {
// 	types := []uint32{}
// 	for _, typ := range eventTypes {
// 		types = append(types, uint32(typ))
// 	}

// 	ncs, _, err := mgrcli.GetChannels(ctx, &mgrpb.Conds{
// 		AppID: &commonpb.StringVal{
// 			Op:    cruder.EQ,
// 			Value: appID,
// 		},
// 		EventTypes: &commonpb.Uint32SliceVal{
// 			Op:    cruder.IN,
// 			Value: types,
// 		},
// 		Channel: &commonpb.Uint32Val{
// 			Op:    cruder.EQ,
// 			Value: uint32(channel1),
// 		},
// 	}, 0, int32(len(eventTypes)))
// 	if err != nil {
// 		return nil, err
// 	}

// 	evTypes := []basetypes.UsedFor{}

// nextType:
// 	for _, typ := range eventTypes {
// 		for _, nc := range ncs {
// 			if nc.EventType == typ {
// 				continue nextType
// 			}
// 		}
// 		evTypes = append(evTypes, typ)
// 	}

// 	if len(evTypes) == 0 {
// 		return nil, nil
// 	}

// 	var req []*mgrpb.ChannelReq
// 	for key := range evTypes {
// 		req = append(req, &mgrpb.ChannelReq{
// 			AppID:     &appID,
// 			EventType: &evTypes[key],
// 			Channel:   &channel1,
// 		})
// 	}

// 	rows, err := mgrcli.CreateChannels(ctx, req)
// 	if err != nil {
// 		return nil, err
// 	}

// 	if len(rows) == 0 {
// 		return nil, nil
// 	}

// 	infos, err := extends(ctx, rows)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return infos, nil
// }

// func GetChannels(
// 	ctx context.Context,
// 	appID string,
// 	offset, limit uint32,
// ) (
// 	[]*npool.Channel,
// 	uint32,
// 	error,
// ) {
// 	rows, total, err := mgrcli.GetChannels(ctx, &mgrpb.Conds{
// 		AppID: &commonpb.StringVal{
// 			Op:    cruder.EQ,
// 			Value: appID,
// 		},
// 	}, int32(offset), int32(limit))
// 	if err != nil {
// 		return nil, 0, err
// 	}

// 	if len(rows) == 0 {
// 		return nil, total, nil
// 	}

// 	infos, err := extends(ctx, rows)
// 	if err != nil {
// 		return nil, 0, err
// 	}
// 	return infos, total, nil
// }

// func extends(
// 	ctx context.Context,
// 	rows []*mgrpb.Channel,
// ) (
// 	[]*npool.Channel,
// 	error,
// ) {
// 	appIDs := []string{}

// 	for _, val := range rows {
// 		appIDs = append(appIDs, val.AppID)
// 	}
// 	appInfos, _, err := appmwcli.GetApps(ctx, &appmwpb.Conds{
// 		IDs: &basetypes.StringSliceVal{Op: cruder.IN, Value: appIDs},
// 	}, 0, int32(len(appIDs)))
// 	if err != nil {
// 		return nil, err
// 	}
// 	appMap := map[string]*appmwpb.App{}
// 	for _, val := range appInfos {
// 		appMap[val.ID] = val
// 	}

// 	infos := []*npool.Channel{}
// 	for _, val := range rows {
// 		app, ok := appMap[val.AppID]
// 		if !ok {
// 			continue
// 		}

// 		infos = append(infos, &npool.Channel{
// 			ID:        val.ID,
// 			AppID:     val.AppID,
// 			AppName:   app.Name,
// 			EventType: val.EventType,
// 			Channel:   val.Channel,
// 			CreatedAt: val.CreatedAt,
// 		})
// 	}
// 	return infos, nil
// }
