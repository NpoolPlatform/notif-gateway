package frontend

import (
	"context"
	"fmt"

	cruder "github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"

	applangmwcli "github.com/NpoolPlatform/g11n-middleware/pkg/client/applang"
	applangpb "github.com/NpoolPlatform/message/npool/g11n/mw/v1/applang"

	frontendtemplatemwpb "github.com/NpoolPlatform/message/npool/notif/mw/v1/template/frontend"
	frontendtemplatemwcli "github.com/NpoolPlatform/notif-middleware/pkg/client/template/frontend"
)

func (h *Handler) CreateFrontendTemplate(ctx context.Context) (*frontendtemplatemwpb.FrontendTemplate, error) {
	existLang, err := applangmwcli.ExistAppLangConds(ctx, &applangpb.Conds{
		AppID:  &basetypes.StringVal{Op: cruder.EQ, Value: *h.AppID},
		LangID: &basetypes.StringVal{Op: cruder.EQ, Value: *h.LangID},
	})
	if err != nil {
		return nil, err
	}
	if !existLang {
		return nil, fmt.Errorf("invalid applang")
	}

	return frontendtemplatemwcli.CreateFrontendTemplate(ctx, &frontendtemplatemwpb.FrontendTemplateReq{
		AppID:   h.AppID,
		LangID:  h.LangID,
		UsedFor: h.UsedFor,
		Title:   h.Title,
		Content: h.Content,
	})
}
