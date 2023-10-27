package sms

import (
	"context"
	"fmt"

	"github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"
	smstemplatemwpb "github.com/NpoolPlatform/message/npool/notif/mw/v1/template/sms"
	smstemplatemwcli "github.com/NpoolPlatform/notif-middleware/pkg/client/template/sms"
)

func (h *Handler) UpdateSMSTemplate(ctx context.Context) (*smstemplatemwpb.SMSTemplate, error) {
	exist, err := smstemplatemwcli.ExistSMSTemplateConds(ctx, &smstemplatemwpb.Conds{
		ID:    &basetypes.Uint32Val{Op: cruder.EQ, Value: *h.ID},
		AppID: &basetypes.StringVal{Op: cruder.EQ, Value: *h.AppID},
	})
	if err != nil {
		return nil, err
	}
	if !exist {
		return nil, fmt.Errorf("sms template not exist")
	}

	info, err := smstemplatemwcli.UpdateSMSTemplate(ctx, &smstemplatemwpb.SMSTemplateReq{
		ID:      h.ID,
		AppID:   h.AppID,
		Subject: h.Subject,
		Message: h.Message,
	})
	if err != nil {
		return nil, err
	}

	h.EntID = &info.EntID

	return h.GetSMSTemplate(ctx)
}
