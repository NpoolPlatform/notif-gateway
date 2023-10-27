package frontend

import (
	"context"
	"fmt"

	"github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"
	frontendtemplatemwpb "github.com/NpoolPlatform/message/npool/notif/mw/v1/template/frontend"
	frontendtemplatemwcli "github.com/NpoolPlatform/notif-middleware/pkg/client/template/frontend"
)

func (h *Handler) UpdateFrontendTemplate(ctx context.Context) (*frontendtemplatemwpb.FrontendTemplate, error) {
	exist, err := frontendtemplatemwcli.ExistFrontendTemplateConds(ctx, &frontendtemplatemwpb.Conds{
		ID:    &basetypes.Uint32Val{Op: cruder.EQ, Value: *h.ID},
		AppID: &basetypes.StringVal{Op: cruder.EQ, Value: *h.AppID},
	})
	if err != nil {
		return nil, err
	}
	if !exist {
		return nil, fmt.Errorf("frontend template not exist")
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

	h.EntID = &info.EntID

	return h.GetFrontendTemplate(ctx)
}
