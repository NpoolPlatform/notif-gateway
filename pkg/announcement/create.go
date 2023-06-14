package announcement

import (
	"context"
	"fmt"

	npool "github.com/NpoolPlatform/message/npool/notif/mw/v1/announcement"
	cli "github.com/NpoolPlatform/notif-middleware/pkg/client/announcement"
)

type createHandler struct {
	*Handler
}

func (h *createHandler) validate() error {
	if h.Title == nil {
		return fmt.Errorf("title is empty")
	}
	if h.Content == nil {
		return fmt.Errorf("content is empty")
	}
	if h.Type == nil {
		return fmt.Errorf("type is empty")
	}
	if &h.EndAt == nil {
		return fmt.Errorf("endat is empty")
	}
	return nil
}

func (h *Handler) CreateAnnouncement(ctx context.Context) (*npool.Announcement, error) {
	handler := &createHandler{
		Handler: h,
	}

	if err := handler.validate(); err != nil {
		return nil, err
	}

	info, err := cli.CreateAnnouncement(
		ctx,
		&npool.AnnouncementReq{
			AppID:            h.AppID,
			Title:            h.Title,
			Content:          h.Content,
			LangID:           h.LangID,
			Channel:          h.Channel,
			AnnouncementType: h.Type,
			EndAt:            &h.EndAt,
		},
	)
	if err != nil {
		return nil, err
	}

	h.ID = &info.ID
	return h.GetAnnouncement(ctx)
}
