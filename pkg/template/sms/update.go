package sms

import (
	"context"
	"fmt"

	smstemplatemwpb "github.com/NpoolPlatform/message/npool/notif/mw/v1/template/sms"
	smstemplatemwcli "github.com/NpoolPlatform/notif-middleware/pkg/client/template/sms"
)

func (h *Handler) UpdateSMSTemplate(ctx context.Context) (*smstemplatemwpb.SMSTemplate, error) {
	if h.ID == nil || *h.ID == "" {
		return nil, fmt.Errorf("id invalid")
	}
	smsInfo, err := h.GetSMSTemplate(ctx)
	if err != nil {
		return nil, err
	}
	if smsInfo == nil {
		return nil, fmt.Errorf("sms template not exist")
	}
	if h.AppID == nil || *h.AppID == "" {
		return nil, fmt.Errorf("appid invalid")
	}
	if smsInfo.AppID != *h.AppID {
		return nil, fmt.Errorf("permission denied")
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

	return info, nil
}
