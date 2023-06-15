package email

import (
	"context"
	"fmt"

	cruder "github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"

	applangmwcli "github.com/NpoolPlatform/g11n-middleware/pkg/client/applang"
	applangmwpb "github.com/NpoolPlatform/message/npool/g11n/mw/v1/applang"

	emailtemplatemwpb "github.com/NpoolPlatform/message/npool/notif/mw/v1/template/email"
	emailtemplatemwcli "github.com/NpoolPlatform/notif-middleware/pkg/client/template/email"
)

func (h *Handler) CreateEmailTemplate(ctx context.Context) (*emailtemplatemwpb.EmailTemplate, error) {
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

	return emailtemplatemwcli.CreateEmailTemplate(ctx, &emailtemplatemwpb.EmailTemplateReq{
		ID:                h.ID,
		AppID:             h.AppID,
		LangID:            h.LangID,
		UsedFor:           h.UsedFor,
		Sender:            h.Sender,
		ReplyTos:          h.ReplyTos,
		CCTos:             h.CCTos,
		Subject:           h.Subject,
		Body:              h.Body,
		DefaultToUsername: h.DefaultToUsername,
	})
}
