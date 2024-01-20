package announcement

import (
	"context"
	"fmt"
	"time"

	usermwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/user"
	applangmwcli "github.com/NpoolPlatform/g11n-middleware/pkg/client/applang"
	"github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	usermwpb "github.com/NpoolPlatform/message/npool/appuser/mw/v1/user"
	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"
	"github.com/NpoolPlatform/message/npool/g11n/mw/v1/applang"
	npool "github.com/NpoolPlatform/message/npool/notif/gw/v1/announcement"
	mwpb "github.com/NpoolPlatform/message/npool/notif/mw/v1/announcement"
	mwcli "github.com/NpoolPlatform/notif-middleware/pkg/client/announcement"
)

func (h *Handler) GetAnnouncements(ctx context.Context) ([]*npool.Announcement, uint32, error) {
	if h.UserID != nil {
		existUser, err := usermwcli.ExistUser(ctx, *h.AppID, *h.UserID)
		if err != nil {
			return nil, 0, err
		}
		if !existUser {
			return nil, 0, fmt.Errorf("invalid user")
		}
	}

	existLang, err := applangmwcli.ExistAppLangConds(ctx, &applang.Conds{
		AppID:  &basetypes.StringVal{Op: cruder.EQ, Value: *h.AppID},
		LangID: &basetypes.StringVal{Op: cruder.EQ, Value: *h.LangID},
	})
	if err != nil {
		return nil, 0, err
	}
	if !existLang {
		return nil, 0, fmt.Errorf("invalid applang")
	}

	conds := &mwpb.Conds{
		AppID:  &basetypes.StringVal{Op: cruder.EQ, Value: *h.AppID},
		LangID: &basetypes.StringVal{Op: cruder.EQ, Value: *h.LangID},
	}
	if h.UserID != nil {
		conds.UserID = &basetypes.StringVal{Op: cruder.EQ, Value: *h.UserID}
	}

	infos, total, err := mwcli.GetAnnouncements(ctx, conds, h.Offset, h.Limit)
	if err != nil {
		return nil, 0, err
	}
	if len(infos) == 0 {
		return nil, 0, nil
	}

	var user *usermwpb.User
	if h.UserID != nil {
		user, err = usermwcli.GetUser(ctx, *h.AppID, *h.UserID)
		if err != nil {
			return nil, 0, err
		}
	}

	announcements := []*npool.Announcement{}
	now := uint32(time.Now().Unix())
	for _, amt := range infos {
		if amt.StartAt >= now {
			continue
		}
		announcement := &npool.Announcement{
			ID:               amt.ID,
			EntID:            amt.EntID,
			AppID:            amt.AppID,
			LangID:           amt.LangID,
			Title:            amt.Title,
			Content:          amt.Content,
			EndAt:            amt.EndAt,
			StartAt:          amt.StartAt,
			Notified:         amt.Notified,
			Channel:          amt.Channel,
			AnnouncementType: amt.AnnouncementType,
			CreatedAt:        amt.CreatedAt,
			UpdatedAt:        amt.UpdatedAt,
		}
		if user != nil {
			announcement.UserID = user.EntID
			announcement.EmailAddress = user.EmailAddress
			announcement.PhoneNO = user.PhoneNO
			announcement.Username = user.Username
		}
		announcements = append(announcements, announcement)
	}

	return announcements, total, nil
}

func (h *Handler) GetAppAnnouncements(ctx context.Context) ([]*npool.Announcement, uint32, error) {
	conds := &mwpb.Conds{AppID: &basetypes.StringVal{Op: cruder.EQ, Value: *h.AppID}}
	infos, total, err := mwcli.GetAnnouncements(ctx, conds, h.Offset, h.Limit)
	if err != nil {
		return nil, 0, err
	}
	if len(infos) == 0 {
		return nil, total, nil
	}

	rows := []*npool.Announcement{}
	for _, amt := range infos {
		row := &npool.Announcement{
			ID:               amt.ID,
			EntID:            amt.EntID,
			AppID:            amt.AppID,
			LangID:           amt.LangID,
			Title:            amt.Title,
			Content:          amt.Content,
			StartAt:          amt.StartAt,
			EndAt:            amt.EndAt,
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
	amt, err := mwcli.GetAnnouncement(ctx, *h.EntID)
	if err != nil {
		return nil, err
	}
	if amt == nil {
		return nil, fmt.Errorf("announcement not exist")
	}

	info := &npool.Announcement{
		ID:               amt.ID,
		EntID:            amt.EntID,
		AppID:            amt.AppID,
		LangID:           amt.LangID,
		Title:            amt.Title,
		Content:          amt.Content,
		StartAt:          amt.StartAt,
		EndAt:            amt.EndAt,
		Channel:          amt.Channel,
		AnnouncementType: amt.AnnouncementType,
		CreatedAt:        amt.CreatedAt,
		UpdatedAt:        amt.UpdatedAt,
	}

	return info, nil
}

func (h *Handler) GetAnnouncementExt(amt *mwpb.Announcement) (*npool.Announcement, error) {
	info := &npool.Announcement{
		ID:               amt.ID,
		EntID:            amt.EntID,
		AppID:            amt.AppID,
		LangID:           amt.LangID,
		Title:            amt.Title,
		Content:          amt.Content,
		StartAt:          amt.StartAt,
		EndAt:            amt.EndAt,
		Channel:          amt.Channel,
		AnnouncementType: amt.AnnouncementType,
		CreatedAt:        amt.CreatedAt,
		UpdatedAt:        amt.UpdatedAt,
	}

	return info, nil
}
