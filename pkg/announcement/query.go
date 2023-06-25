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
	mwcli "github.com/NpoolPlatform/notif-middleware/pkg/client/announcement"
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

	if h.UserID != nil {
		rows, err := formalize(ctx, *h.UserID, infos)
		if err != nil {
			return nil, 0, err
		}
		return rows, total, nil
	}

	rows := []*npool.Announcement{}
	for _, amt := range infos {
		row := &npool.Announcement{
			ID:               amt.ID,
			AppID:            amt.AppID,
			UserID:           "",
			EmailAddress:     "",
			PhoneNO:          "",
			Username:         "",
			LangID:           amt.LangID,
			Title:            amt.Title,
			Content:          amt.Content,
			StartAt:          amt.StartAt,
			EndAt:            amt.EndAt,
			Notified:         false,
			CreatedAt:        amt.CreatedAt,
			UpdatedAt:        amt.UpdatedAt,
			Channel:          amt.Channel,
			AnnouncementType: amt.AnnouncementType,
		}
		rows = append(rows, row)
	}
	return rows, total, nil
}

func (h *Handler) GetAnnouncement(ctx context.Context) (*npool.Announcement, error) {
	amt, err := mwcli.GetAnnouncement(ctx, *h.ID)
	if err != nil {
		return nil, err
	}
	if amt == nil {
		return nil, fmt.Errorf("announcement not exist")
	}

	info := &npool.Announcement{
		ID:               amt.ID,
		AppID:            amt.AppID,
		UserID:           "",
		EmailAddress:     "",
		PhoneNO:          "",
		Username:         "",
		LangID:           amt.LangID,
		Title:            amt.Title,
		Content:          amt.Content,
		StartAt:          amt.StartAt,
		EndAt:            amt.EndAt,
		Notified:         false,
		CreatedAt:        amt.CreatedAt,
		UpdatedAt:        amt.UpdatedAt,
		Channel:          amt.Channel,
		AnnouncementType: amt.AnnouncementType,
	}

	return info, nil
}

func formalize(ctx context.Context, userID string, amts []*mwpb.Announcement) ([]*npool.Announcement, error) {
	if len(amts) == 0 {
		return nil, nil
	}

	// get users
	userIDs := []string{userID}
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
	user, ok := userMap[userID]
	if !ok {
		return nil, fmt.Errorf("invalid user id")
	}

	for _, amt := range amts {
		announcements = append(announcements, &npool.Announcement{
			ID:               amt.ID,
			AppID:            amt.AppID,
			UserID:           user.ID,
			EmailAddress:     user.EmailAddress,
			PhoneNO:          user.PhoneNO,
			Username:         user.Username,
			LangID:           amt.LangID,
			Title:            amt.Title,
			Content:          amt.Content,
			EndAt:            amt.EndAt,
			Notified:         amt.Notified,
			Channel:          amt.Channel,
			AnnouncementType: amt.AnnouncementType,
			CreatedAt:        amt.CreatedAt,
			UpdatedAt:        amt.UpdatedAt,
		})
	}
	return announcements, nil
}
