package announcement

import (
	"context"
	"fmt"

	usermwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/user"
	"github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"
	npool "github.com/NpoolPlatform/message/npool/notif/gw/v1/announcement"
	mwpb "github.com/NpoolPlatform/message/npool/notif/mw/v1/announcement"
	mwcli "github.com/NpoolPlatform/notif-middleware/pkg/client/announcement"
)

func (h *Handler) GetAnnouncements(ctx context.Context) ([]*npool.Announcement, uint32, error) {
	conds := &mwpb.Conds{
		AppID:  &basetypes.StringVal{Op: cruder.EQ, Value: *h.AppID},
		UserID: &basetypes.StringVal{Op: cruder.EQ, Value: *h.UserID},
		LangID: &basetypes.StringVal{Op: cruder.EQ, Value: *h.LangID},
	}

	infos, total, err := mwcli.GetAnnouncements(ctx, conds, h.Offset, h.Limit)
	if err != nil {
		return nil, 0, err
	}
	if len(infos) == 0 {
		return nil, total, nil
	}

	user, err := usermwcli.GetUser(ctx, *h.AppID, *h.UserID)
    fmt.Println("=============User-----------", err)
	if err != nil {
		return nil, 0, err
	}

	announcements := []*npool.Announcement{}
	for _, amt := range infos {
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
			StartAt:          amt.StartAt,
			Notified:         amt.Notified,
			Channel:          amt.Channel,
			AnnouncementType: amt.AnnouncementType,
			CreatedAt:        amt.CreatedAt,
			UpdatedAt:        amt.UpdatedAt,
		})
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

	return info, nil
}
