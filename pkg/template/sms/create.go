package sms

import (
	"context"
	"fmt"

	cruder "github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"

	applangmwcli "github.com/NpoolPlatform/g11n-middleware/pkg/client/applang"
	applangmwpb "github.com/NpoolPlatform/message/npool/g11n/mw/v1/applang"

	smstemplatemwpb "github.com/NpoolPlatform/message/npool/notif/mw/v1/template/sms"
	smstemplatemwcli "github.com/NpoolPlatform/notif-middleware/pkg/client/template/sms"
)

func (h *Handler) CreateSMSTemplate(ctx context.Context) (*smstemplatemwpb.SMSTemplate, error) {
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

	return smstemplatemwcli.CreateSMSTemplate(ctx, &smstemplatemwpb.SMSTemplateReq{
		ID:      h.ID,
		AppID:   h.AppID,
		LangID:  h.LangID,
		UsedFor: h.UsedFor,
		Subject: h.Subject,
		Message: h.Message,
	})
}
