package sms

import (
	"context"

	cruder "github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"
	smstemplatemwpb "github.com/NpoolPlatform/message/npool/notif/mw/v1/template/sms"
	smstemplatemwcli "github.com/NpoolPlatform/notif-middleware/pkg/client/template/sms"
)

func (h *Handler) GetSMSTemplate(ctx context.Context) (*smstemplatemwpb.SMSTemplate, error) {
	return smstemplatemwcli.GetSMSTemplate(ctx, *h.EntID)
}

func (h *Handler) GetSMSTemplates(ctx context.Context) ([]*smstemplatemwpb.SMSTemplate, uint32, error) {
	return smstemplatemwcli.GetSMSTemplates(ctx, &smstemplatemwpb.Conds{
		AppID: &basetypes.StringVal{Op: cruder.EQ, Value: *h.AppID}},
		h.Offset,
		h.Limit,
	)
}
