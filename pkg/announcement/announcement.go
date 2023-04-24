package announcement

import (
	"context"
	"fmt"

	usercli "github.com/NpoolPlatform/appuser-middleware/pkg/client/user"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	"github.com/NpoolPlatform/libent-cruder/pkg/cruder"

	channelpb "github.com/NpoolPlatform/message/npool/notif/mgr/v1/channel"

	npoolpb "github.com/NpoolPlatform/message/npool"
	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"

	npool "github.com/NpoolPlatform/message/npool/notif/gw/v1/announcement"
	mgrpb "github.com/NpoolPlatform/message/npool/notif/mgr/v1/announcement"
	mgrcli "github.com/NpoolPlatform/notif-manager/pkg/client/announcement"

	mwpb "github.com/NpoolPlatform/message/npool/notif/mw/v1/announcement"
	mwcli "github.com/NpoolPlatform/notif-middleware/pkg/client/announcement"

	appmwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/app"
	appmwpb "github.com/NpoolPlatform/message/npool/appuser/mw/v1/app"
)

func CreateAnnouncement(
	ctx context.Context,
	appID, langID, title, content string,
	channel channelpb.NotifChannel,
	endAt uint32,
	announcementType mgrpb.AnnouncementType,
) (*npool.Announcement, error) {
	info, err := mgrcli.CreateAnnouncement(ctx, &mgrpb.AnnouncementReq{
		AppID:            &appID,
		LangID:           &langID,
		Title:            &title,
		Content:          &content,
		Channel:          &channel,
		EndAt:            &endAt,
		AnnouncementType: &announcementType,
	})
	if err != nil {
		return nil, err
	}

	return expand(ctx, info)
}

func UpdateAnnouncement(
	ctx context.Context,
	id string,
	title, content *string,
	endAt *uint32,
	announcementType *mgrpb.AnnouncementType,
) (
	*npool.Announcement,
	error,
) {
	info, err := mgrcli.UpdateAnnouncement(ctx, &mgrpb.AnnouncementReq{
		ID:               &id,
		Title:            title,
		Content:          content,
		EndAt:            endAt,
		AnnouncementType: announcementType,
	})
	if err != nil {
		return nil, err
	}

	return expand(ctx, info)
}

func DeleteAnnouncement(
	ctx context.Context,
	id string,
) (
	*npool.Announcement,
	error,
) {
	_, err := GetAnnouncement(ctx, id)
	if err != nil {
		return nil, err
	}
	info, err := mgrcli.DeleteAnnouncement(ctx, id)
	if err != nil {
		return nil, err
	}

	return expand(ctx, info)
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

	return expand(ctx, info)
}

func GetAppAnnouncements(
	ctx context.Context,
	appID string,
	offset, limit uint32,
) (
	[]*npool.Announcement,
	uint32,
	error,
) {
	if limit == 0 {
		limit = 100
	}
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

	infos := []*npool.Announcement{}
	for _, r := range rows {
		app, ok := appMap[r.AppID]
		if !ok {
			continue
		}
		infos = append(infos, &npool.Announcement{
			ID:               r.ID,
			AppID:            r.AppID,
			AppName:          app.Name,
			LangID:           r.LangID,
			Title:            r.Title,
			Content:          r.Content,
			CreatedAt:        r.CreatedAt,
			UpdatedAt:        r.UpdatedAt,
			EndAt:            r.EndAt,
			Channel:          r.Channel,
			AnnouncementType: r.AnnouncementType,
		})
	}

	return infos, total, nil
}

func GetAnnouncements(
	ctx context.Context,
	appID, userID, langID string,
	offset, limit uint32,
) (
	[]*npool.Announcement,
	uint32,
	error,
) {
	if limit == 0 {
		limit = 100
	}
	rows, total, err := mwcli.GetAnnouncementStates(ctx, &mwpb.Conds{
		AppID: &npoolpb.StringVal{
			Op:    cruder.EQ,
			Value: appID,
		},
		UserID: &npoolpb.StringVal{
			Op:    cruder.EQ,
			Value: userID,
		},
		LangID: &npoolpb.StringVal{
			Op:    cruder.EQ,
			Value: langID,
		},
	}, int32(offset), int32(limit))
	if err != nil {
		return nil, 0, err
	}

	if len(rows) == 0 {
		return nil, total, nil
	}
	appIDs := []string{}
	for _, r := range rows {
		appIDs = append(appIDs, r.AppID)
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

	userInfo, err := usercli.GetUser(ctx, appID, userID)
	if err != nil {
		return nil, 0, err
	}
	if userInfo == nil {
		return nil, 0, nil
	}

	infos := []*npool.Announcement{}
	for _, r := range rows {
		app, ok := appMap[r.AppID]
		if !ok {
			continue
		}

		infos = append(infos, &npool.Announcement{
			ID:               r.AnnouncementID,
			AppID:            r.AppID,
			AppName:          app.Name,
			UserID:           userInfo.ID,
			EmailAddress:     userInfo.EmailAddress,
			PhoneNO:          userInfo.PhoneNO,
			Username:         userInfo.Username,
			LangID:           r.LangID,
			Title:            r.Title,
			Content:          r.Content,
			Read:             r.Read,
			CreatedAt:        r.CreatedAt,
			UpdatedAt:        r.UpdatedAt,
			EndAt:            r.EndAt,
			Channel:          r.Channel,
			AnnouncementType: r.AnnouncementType,
		})
	}

	return infos, total, nil
}

func expand(
	ctx context.Context,
	info *mgrpb.Announcement,
) (
	*npool.Announcement,
	error,
) {
	appInfo, err := appmwcli.GetApp(ctx, info.AppID)
	if err != nil {
		return nil, err
	}

	appName := ""
	if appInfo == nil {
		appName = appInfo.Name
	}
	return &npool.Announcement{
		ID:               info.ID,
		AppID:            info.AppID,
		AppName:          appName,
		LangID:           info.LangID,
		Title:            info.Title,
		Content:          info.Content,
		CreatedAt:        info.CreatedAt,
		UpdatedAt:        info.UpdatedAt,
		EndAt:            info.EndAt,
		Channel:          info.Channel,
		AnnouncementType: info.AnnouncementType,
	}, nil
}
