package announcement

import (
	"context"
	"fmt"

	usermwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/user"
	"github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	usermwpb "github.com/NpoolPlatform/message/npool/appuser/mw/v1/user"
	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"
	npool "github.com/NpoolPlatform/message/npool/notif/gw/v1/announcement"
	mwpb "github.com/NpoolPlatform/message/npool/notif/mw/v1/announcement"
	readpb "github.com/NpoolPlatform/message/npool/notif/mw/v1/announcement/readstate"
	cli "github.com/NpoolPlatform/notif-middleware/pkg/client/announcement"
	mwcli "github.com/NpoolPlatform/notif-middleware/pkg/client/announcement"
	readcli "github.com/NpoolPlatform/notif-middleware/pkg/client/announcement/readstate"
)

func (h *Handler) GetAnnouncements(ctx context.Context) ([]*npool.Announcement, uint32, error) {
	conds := &mwpb.Conds{
		AppID: &basetypes.StringVal{
			Op:    cruder.EQ,
			Value: *h.AppID,
		},
	}
	if h.UserID != nil {
		conds.UserID = &basetypes.StringVal{
			Op:    cruder.EQ,
			Value: *h.UserID,
		}
	}
	if h.LangID != nil {
		conds.LangID = &basetypes.StringVal{
			Op:    cruder.EQ,
			Value: *h.LangID,
		}
	}

	infos, total, err := mwcli.GetAnnouncements(ctx, conds, h.Offset, h.Limit)
	if err != nil {
		return nil, 0, err
	}
	if len(infos) == 0 {
		return nil, total, nil
	}

	rows, err := formalize(ctx, infos[0].AppID, infos[0].UserID, infos)
	if err != nil {
		return nil, 0, err
	}

	return rows, total, nil
}

func (h *Handler) GetAnnouncement(ctx context.Context) (*npool.Announcement, error) {
	info, err := cli.GetAnnouncement(ctx, *h.ID)
	if err != nil {
		return nil, err
	}
	if info == nil {
		return nil, fmt.Errorf("announcement not exist")
	}

	infos, err := formalize(ctx, info.AppID, info.UserID, []*mwpb.Announcement{info})
	if err != nil {
		return nil, err
	}
	if len(infos) == 0 {
		return nil, nil
	}

	return infos[0], nil
}

func formalize(ctx context.Context, appID, userID string, amts []*mwpb.Announcement) ([]*npool.Announcement, error) {
	if len(amts) == 0 {
		return nil, nil
	}

	amtIDs := []string{}
	userIDs := []string{}
	for _, amt := range amts {
		amtIDs = append(amtIDs, amt.ID)
		userIDs = append(userIDs, amt.UserID)
	}
	if len(amtIDs) == 0 {
		return nil, nil
	}

	infos, _, err := readcli.GetReadStates(ctx, &readpb.Conds{
		AppID: &basetypes.StringVal{
			Op:    cruder.EQ,
			Value: appID,
		},
		UserID: &basetypes.StringVal{
			Op:    cruder.EQ,
			Value: userID,
		},
		AnnouncementIDs: &basetypes.StringSliceVal{
			Op:    cruder.IN,
			Value: amtIDs,
		},
	}, 0, int32(len(amts)))
	if err != nil {
		return nil, err
	}

	readMap := map[string]*readpb.ReadState{}
	for _, info := range infos {
		readMap[info.AnnouncementID] = info
	}

	// get users
	userMap := map[string]*usermwpb.User{}
	if len(userIDs) > 0 {
		userInfos, _, err := usermwcli.GetUsers(ctx, &usermwpb.Conds{
			IDs: &basetypes.StringSliceVal{Op: cruder.IN, Value: userIDs},
		}, 0, int32(len(userIDs)))
		if err != nil {
			return nil, err
		}

		for _, val := range userInfos {
			userMap[val.ID] = val
		}
	}

	announcements := []*npool.Announcement{}
	for _, amt := range amts {
		notified := true
		_, ok := readMap[amt.ID]
		if !ok {
			notified = false
		}
		user, ok := userMap[amt.UserID]
		if !ok {
			continue
		}
		announcements = append(announcements, &npool.Announcement{
			ID:               amt.ID,
			AppID:            amt.AppID,
			UserID:           amt.UserID,
			EmailAddress:     user.EmailAddress,
			PhoneNO:          user.PhoneNO,
			Username:         user.Username,
			LangID:           amt.LangID,
			Title:            amt.Title,
			Content:          amt.Content,
			EndAt:            amt.EndAt,
			Notified:         notified,
			CreatedAt:        amt.CreatedAt,
			UpdatedAt:        amt.UpdatedAt,
			Channel:          amt.Channel,
			AnnouncementType: amt.AnnouncementType,
		})

	}
	return announcements, nil
}
