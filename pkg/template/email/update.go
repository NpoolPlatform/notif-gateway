package email

import (
	"context"
	"fmt"

	cruder "github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"

	emailtemplatemwpb "github.com/NpoolPlatform/message/npool/notif/mw/v1/template/email"
	emailtemplatemwcli "github.com/NpoolPlatform/notif-middleware/pkg/client/template/email"

	applangmwcli "github.com/NpoolPlatform/g11n-middleware/pkg/client/applang"
	applangmwpb "github.com/NpoolPlatform/message/npool/g11n/mw/v1/applang"
)

func (h *Handler) UpdateEmailTemplate(ctx context.Context) (*emailtemplatemwpb.EmailTemplate, error) {
	if h.ID == nil || *h.ID == "" {
		return nil, fmt.Errorf("id invalid")
	}
	emailInfo, err := h.GetEmailTemplate(ctx)
	if err != nil {
		return nil, err
	}
	if emailInfo == nil {
		return nil, fmt.Errorf("email template not exist")
	}
	if emailInfo.AppID != *h.AppID {
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

	info, err := emailtemplatemwcli.UpdateEmailTemplate(ctx, &emailtemplatemwpb.EmailTemplateReq{
		ID:                h.ID,
		LangID:            h.LangID,
		Sender:            h.Sender,
		ReplyTos:          h.ReplyTos,
		CCTos:             h.CCTos,
		Subject:           h.Subject,
		Body:              h.Body,
		DefaultToUsername: h.DefaultToUsername,
	})
	if err != nil {
		return nil, err
	}

	return info, nil
}
