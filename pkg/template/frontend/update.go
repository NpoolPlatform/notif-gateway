package frontend

import (
	"context"
	"fmt"

	cruder "github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"

	frontendtemplatemwpb "github.com/NpoolPlatform/message/npool/notif/mw/v1/template/frontend"
	frontendtemplatemwcli "github.com/NpoolPlatform/notif-middleware/pkg/client/template/frontend"

	applangmwcli "github.com/NpoolPlatform/g11n-middleware/pkg/client/applang"
	applangmwpb "github.com/NpoolPlatform/message/npool/g11n/mw/v1/applang"
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
	if frontendInfo.AppID != *h.AppID {
		return nil, fmt.Errorf("permission denied")
	}

	if h.AppID == nil || *h.AppID == "" {
		return nil, fmt.Errorf("appid invalid")
	}
	if h.LangID == nil || *h.LangID == "" {
		return nil, fmt.Errorf("langid invalid")
	}
	applangInfo, err := applangmwcli.GetLangOnly(ctx, &applangmwpb.Conds{
		AppID: &basetypes.StringVal{
			Op:    cruder.EQ,
			Value: *h.AppID,
		},
		LangID: &basetypes.StringVal{
			Op:    cruder.EQ,
			Value: *h.LangID,
		},
	})
	if err != nil {
		return nil, err
	}
	if applangInfo == nil {
		return nil, fmt.Errorf("applang not exist")
	}

	info, err := frontendtemplatemwcli.UpdateFrontendTemplate(ctx, &frontendtemplatemwpb.FrontendTemplateReq{
		ID:      h.ID,
		AppID:   h.AppID,
		LangID:  h.LangID,
		Title:   h.Title,
		Content: h.Content,
	})
	if err != nil {
		return nil, err
	}

	return info, nil
}
