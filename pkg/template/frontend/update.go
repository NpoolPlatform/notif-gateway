package frontend

import (
	"context"
	"fmt"

	frontendtemplatemwpb "github.com/NpoolPlatform/message/npool/notif/mw/v1/template/frontend"
	frontendtemplatemwcli "github.com/NpoolPlatform/notif-middleware/pkg/client/template/frontend"
)

func (h *Handler) UpdateFrontendTemplate(ctx context.Context) (*frontendtemplatemwpb.FrontendTemplate, error) {
	if h.ID == nil || *h.ID == "" {
		return nil, fmt.Errorf("id invalid")
	}
	frontendInfo, err := h.GetFrontendTemplate(ctx)
	if err != nil {
		return nil, err
	}
	if frontendInfo == nil {
		return nil, fmt.Errorf("frontend template not exist")
	}
	if h.AppID == nil || *h.AppID == "" {
		return nil, fmt.Errorf("appid invalid")
	}
	if frontendInfo.AppID != *h.AppID {
		return nil, fmt.Errorf("permission denied")
	}

	info, err := frontendtemplatemwcli.UpdateFrontendTemplate(ctx, &frontendtemplatemwpb.FrontendTemplateReq{
		ID:      h.ID,
		AppID:   h.AppID,
		Title:   h.Title,
		Content: h.Content,
	})
	if err != nil {
		return nil, err
	}

	return info, nil
}
