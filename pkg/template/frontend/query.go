package frontend

import (
	"context"

	cruder "github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"
	frontendtemplatemwpb "github.com/NpoolPlatform/message/npool/notif/mw/v1/template/frontend"
	frontendtemplatemwcli "github.com/NpoolPlatform/notif-middleware/pkg/client/template/frontend"
)

func (h *Handler) GetFrontendTemplate(ctx context.Context) (*frontendtemplatemwpb.FrontendTemplate, error) {
	return frontendtemplatemwcli.GetFrontendTemplate(ctx, *h.EntID)
}

func (h *Handler) GetFrontendTemplates(ctx context.Context) ([]*frontendtemplatemwpb.FrontendTemplate, uint32, error) {
	return frontendtemplatemwcli.GetFrontendTemplates(ctx, &frontendtemplatemwpb.Conds{
		AppID: &basetypes.StringVal{Op: cruder.EQ, Value: *h.AppID}},
		h.Offset,
		h.Limit,
	)
}
