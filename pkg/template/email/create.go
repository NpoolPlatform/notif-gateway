package email

import (
	"context"
	"fmt"

	cruder "github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"

	applangmwcli "github.com/NpoolPlatform/g11n-middleware/pkg/client/applang"
	applangpb "github.com/NpoolPlatform/message/npool/g11n/mw/v1/applang"

	emailtemplatemwpb "github.com/NpoolPlatform/message/npool/notif/mw/v1/template/email"
	emailtemplatemwcli "github.com/NpoolPlatform/notif-middleware/pkg/client/template/email"
)

func (h *Handler) CreateEmailTemplate(ctx context.Context) (*emailtemplatemwpb.EmailTemplate, error) {
	existLang, err := applangmwcli.ExistAppLangConds(ctx, &applangpb.Conds{
		AppID:  &basetypes.StringVal{Op: cruder.EQ, Value: *h.AppID},
		LangID: &basetypes.StringVal{Op: cruder.EQ, Value: *h.LangID},
	})
	if err != nil {
		return nil, err
	}
	if !existLang {
		return nil, fmt.Errorf("invalid applang")
	}

	return emailtemplatemwcli.CreateEmailTemplate(ctx, &emailtemplatemwpb.EmailTemplateReq{
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
