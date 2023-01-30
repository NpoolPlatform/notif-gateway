package announcement

import (
	"context"
	"fmt"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	"github.com/NpoolPlatform/libent-cruder/pkg/cruder"

	channelpb "github.com/NpoolPlatform/message/npool/notif/mgr/v1/channel"

	npoolpb "github.com/NpoolPlatform/message/npool"

	npool "github.com/NpoolPlatform/message/npool/notif/gw/v1/announcement"
	mgrpb "github.com/NpoolPlatform/message/npool/notif/mgr/v1/announcement"
	mgrcli "github.com/NpoolPlatform/notif-manager/pkg/client/announcement"

	appcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/app"
	apppb "github.com/NpoolPlatform/message/npool/appuser/mw/v1/app"
)

func CreateAnnouncement(
	ctx context.Context,
	appID, title, content string,
	channel []channelpb.NotifChannel,
) (*npool.Announcement, error) {
	info, err := mgrcli.CreateAnnouncement(ctx, &mgrpb.AnnouncementReq{
		AppID:    &appID,
		Title:    &title,
		Content:  &content,
		Channels: channel,
	})
	if err != nil {
		return nil, err
	}

	return expend(ctx, info)
}

func UpdateAnnouncement(
	ctx context.Context,
	id string,
	title, content *string,
	channel []channelpb.NotifChannel,
) (
	*npool.Announcement,
	error,
) {
	info, err := mgrcli.UpdateAnnouncement(ctx, &mgrpb.AnnouncementReq{
		ID:       &id,
		Title:    title,
		Content:  content,
		Channels: channel,
	})
	if err != nil {
		return nil, err
	}

	return expend(ctx, info)
}

func DeleteAnnouncement(
	ctx context.Context,
	id string,
) (
	*npool.Announcement,
	error,
) {
	info, err := mgrcli.DeleteAnnouncement(ctx, id)
	if err != nil {
		return nil, err
	}

	return expend(ctx, info)
}

func GetAnnouncement(
	ctx context.Context,
	id string,
) (
	*npool.Announcement,
	error,
) {
	info, err := mgrcli.GetAnnouncement(ctx, id)
	if err != nil {
		return nil, err
	}
	if info == nil {
		logger.Sugar().Errorw("GetAnnouncement", "id", id, "error", "announcement not exist")
		return nil, fmt.Errorf("announcement not exist")
	}

	return expend(ctx, info)
}

func GetAnnouncements(
	ctx context.Context,
	appID string,
	offset, limit uint32,
) (
	[]*npool.Announcement,
	uint32,
	error,
) {
	rows, total, err := mgrcli.GetAnnouncements(ctx, &mgrpb.Conds{
		AppID: &npoolpb.StringVal{
			Op:    cruder.EQ,
			Value: appID,
		},
	}, int32(offset), int32(limit))
	if err != nil {
		return nil, 0, err
	}
	if len(rows) == 0 {
		return nil, 0, nil
	}
	appIDs := []string{}
	for _, r := range rows {
		appIDs = append(appIDs, r.AppID)
	}

	appInfos, _, err := appcli.GetManyApps(ctx, appIDs)
	if err != nil {
		return nil, 0, err
	}
	appMap := map[string]*apppb.App{}
	for _, val := range appInfos {
		appMap[val.ID] = val
	}

	infos := []*npool.Announcement{}
	for _, r := range rows {
		app, ok := appMap[r.AppID]
		if !ok {
			continue
		}
		infos = append(infos, &npool.Announcement{
			ID:       r.ID,
			AppID:    r.AppID,
			AppName:  app.Name,
			Title:    r.Title,
			Content:  r.Content,
			Channels: r.Channels,
		})
	}

	return infos, total, nil
}

func expend(
	ctx context.Context,
	info *mgrpb.Announcement,
) (
	*npool.Announcement,
	error,
) {
	appInfo, err := appcli.GetApp(ctx, info.AppID)
	if err != nil {
		return nil, err
	}

	return &npool.Announcement{
		ID:       info.ID,
		AppID:    info.AppID,
		AppName:  appInfo.Name,
		Title:    info.Title,
		Content:  info.Content,
		Channels: info.Channels,
	}, nil
}
