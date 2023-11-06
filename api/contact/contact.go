package contact

import (
	"context"

	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"

	npool "github.com/NpoolPlatform/message/npool/notif/gw/v1/contact"
	contactmwpb "github.com/NpoolPlatform/message/npool/notif/mw/v1/contact"
	sendmwpb "github.com/NpoolPlatform/message/npool/third/mw/v1/send"

	contact "github.com/NpoolPlatform/notif-gateway/pkg/contact"
	contact1 "github.com/NpoolPlatform/notif-gateway/pkg/contact/generate"
	contactmwcli "github.com/NpoolPlatform/notif-middleware/pkg/client/contact"
	sendmwcli "github.com/NpoolPlatform/third-middleware/pkg/client/send"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) ContactViaEmail(ctx context.Context, in *npool.ContactViaEmailRequest) (*npool.ContactViaEmailResponse, error) {
	_, err := contact1.NewHandler(
		ctx,
		contact1.WithSubject(&in.Subject, true),
		contact1.WithBody(&in.Body, true),
		contact.WithUsedFor(&in.UsedFor, true),
		contact.WithAppID(&in.AppID, true),
		contact.WithSender(&in.Sender, true),
		contact1.WithSenderName(&in.SenderName, true),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"ContactViaEmail",
			"Req", in,
			"Error", err,
		)
		return &npool.ContactViaEmailResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	info, err := contactmwcli.GenerateContact(ctx, &contactmwpb.GenerateContactRequest{
		AppID:      in.GetAppID(),
		UsedFor:    in.GetUsedFor(),
		Sender:     in.GetSender(),
		Subject:    in.GetSubject(),
		Body:       in.GetBody(),
		SenderName: in.GetSenderName(),
	})
	if err != nil {
		return &npool.ContactViaEmailResponse{}, status.Error(codes.Internal, err.Error())
	}

	err = sendmwcli.SendMessage(ctx, &sendmwpb.SendMessageRequest{
		Subject:     info.Subject,
		Content:     info.Content,
		From:        info.From,
		To:          info.To,
		ToCCs:       info.ToCCs,
		ReplyTos:    info.ReplyTos,
		AccountType: basetypes.SignMethod_Email,
	})
	if err != nil {
		logger.Sugar().Errorw(
			"ContactViaEmail",
			"Req", in,
			"Error", err,
		)
		return &npool.ContactViaEmailResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.ContactViaEmailResponse{}, nil
}
