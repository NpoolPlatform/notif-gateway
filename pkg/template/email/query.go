package email

import (
	"context"

	cruder "github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"
	emailtemplatemwpb "github.com/NpoolPlatform/message/npool/notif/mw/v1/template/email"
	emailtemplatemwcli "github.com/NpoolPlatform/notif-middleware/pkg/client/template/email"
)

func (h *Handler) GetEmailTemplate(ctx context.Context) (*emailtemplatemwpb.EmailTemplate, error) {
	return emailtemplatemwcli.GetEmailTemplate(ctx, *h.ID)
}

func (h *Handler) GetEmailTemplates(ctx context.Context) ([]*emailtemplatemwpb.EmailTemplate, uint32, error) {
	return emailtemplatemwcli.GetEmailTemplates(ctx, &emailtemplatemwpb.Conds{
		AppID: &basetypes.StringVal{Op: cruder.EQ, Value: *h.AppID}},
		h.Offset,
		h.Limit,
	)
}
