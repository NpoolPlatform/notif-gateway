package email

import (
	"context"
	"fmt"

	"github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"
	emailtemplatemwpb "github.com/NpoolPlatform/message/npool/notif/mw/v1/template/email"
	emailtemplatemwcli "github.com/NpoolPlatform/notif-middleware/pkg/client/template/email"
)

func (h *Handler) UpdateEmailTemplate(ctx context.Context) (*emailtemplatemwpb.EmailTemplate, error) {
	exist, err := emailtemplatemwcli.ExistEmailTemplateConds(ctx, &emailtemplatemwpb.Conds{
		ID:    &basetypes.Uint32Val{Op: cruder.EQ, Value: *h.ID},
		AppID: &basetypes.StringVal{Op: cruder.EQ, Value: *h.AppID},
	})
	if err != nil {
		return nil, err
	}
	if !exist {
		return nil, fmt.Errorf("email template not exist")
	}

	info, err := emailtemplatemwcli.UpdateEmailTemplate(ctx, &emailtemplatemwpb.EmailTemplateReq{
		ID:                h.ID,
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

	h.EntID = &info.EntID

	return h.GetEmailTemplate(ctx)
}
