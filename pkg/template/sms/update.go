package sms

import (
	"context"
	"fmt"

	cruder "github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"

	smstemplatemwpb "github.com/NpoolPlatform/message/npool/notif/mw/v1/template/sms"
	smstemplatemwcli "github.com/NpoolPlatform/notif-middleware/pkg/client/template/sms"

	applangmwcli "github.com/NpoolPlatform/g11n-middleware/pkg/client/applang"
	applangmwpb "github.com/NpoolPlatform/message/npool/g11n/mw/v1/applang"
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
	if smsInfo.AppID != *h.AppID {
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

	info, err := smstemplatemwcli.UpdateSMSTemplate(ctx, &smstemplatemwpb.SMSTemplateReq{
		ID:      h.ID,
		AppID:   h.AppID,
		LangID:  h.LangID,
		Subject: h.Subject,
		Message: h.Message,
	})
	if err != nil {
		return nil, err
	}

	return info, nil
}
