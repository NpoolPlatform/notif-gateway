package sms

import (
	"context"
	"fmt"

	smstemplatemwpb "github.com/NpoolPlatform/message/npool/notif/mw/v1/template/sms"
	smstemplatemwcli "github.com/NpoolPlatform/notif-middleware/pkg/client/template/sms"
)

func (h *Handler) UpdateSMSTemplate(ctx context.Context) (*smstemplatemwpb.SMSTemplate, error) {
	smsInfo, err := h.GetSMSTemplate(ctx)
	if err != nil {
		return nil, err
	}
	if smsInfo == nil {
		return nil, fmt.Errorf("sms template not exist")
	}
	if smsInfo.AppID != *h.AppID {
		return nil, fmt.Errorf("permission denied")
	}

	_, err = smstemplatemwcli.UpdateSMSTemplate(ctx, &smstemplatemwpb.SMSTemplateReq{
		ID:      h.ID,
		AppID:   h.AppID,
		Subject: h.Subject,
		Message: h.Message,
	})
	if err != nil {
		return nil, err
	}

	h.EntID = &smsInfo.EntID

	return h.GetSMSTemplate(ctx)
}
