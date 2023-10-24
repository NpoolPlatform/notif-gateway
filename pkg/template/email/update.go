package email

import (
	"context"
	"fmt"

	emailtemplatemwpb "github.com/NpoolPlatform/message/npool/notif/mw/v1/template/email"
	emailtemplatemwcli "github.com/NpoolPlatform/notif-middleware/pkg/client/template/email"
)

func (h *Handler) UpdateEmailTemplate(ctx context.Context) (*emailtemplatemwpb.EmailTemplate, error) {
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

	_, err = emailtemplatemwcli.UpdateEmailTemplate(ctx, &emailtemplatemwpb.EmailTemplateReq{
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

	h.EntID = &emailInfo.EntID

	return h.GetEmailTemplate(ctx)
}
