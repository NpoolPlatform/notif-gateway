package contact

import (
	"context"

	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"
	commontracer "github.com/NpoolPlatform/notif-gateway/pkg/tracer"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"

	npool "github.com/NpoolPlatform/message/npool/notif/gw/v1/contact"
	contactmwpb "github.com/NpoolPlatform/message/npool/notif/mw/v1/contact"
	sendmwpb "github.com/NpoolPlatform/message/npool/third/mw/v1/send"

	constant "github.com/NpoolPlatform/notif-gateway/pkg/message/const"
	contactmwcli "github.com/NpoolPlatform/notif-middleware/pkg/client/contact"
	sendmwcli "github.com/NpoolPlatform/third-middleware/pkg/client/send"

	"go.opentelemetry.io/otel"
	scodes "go.opentelemetry.io/otel/codes"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) ContactViaEmail(ctx context.Context, in *npool.ContactViaEmailRequest) (*npool.ContactViaEmailResponse, error) {
	var err error

	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "ContactViaEmail")
	defer span.End()

	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	if _, err := uuid.Parse(in.GetAppID()); err != nil {
		logger.Sugar().Errorw("ContactViaEmail", "AppID", in.GetAppID())
		return &npool.ContactViaEmailResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	switch in.GetUsedFor() {
	case basetypes.UsedFor_Contact:
	default:
		logger.Sugar().Errorw("ContactViaEmail", "UsedFor", in.GetUsedFor())
		return &npool.ContactViaEmailResponse{}, status.Error(codes.InvalidArgument, "UsedFor is invalid")
	}

	if in.GetSender() == "" {
		logger.Sugar().Errorw("ContactViaEmail", "Sender", in.GetSender())
		return &npool.ContactViaEmailResponse{}, status.Error(codes.InvalidArgument, "Sender is empty")
	}
	if in.GetSubject() == "" {
		logger.Sugar().Errorw("ContactViaEmail", "Subject", in.GetSubject())
		return &npool.ContactViaEmailResponse{}, status.Error(codes.InvalidArgument, "Subject is empty")
	}
	if in.GetBody() == "" {
		logger.Sugar().Errorw("ContactViaEmail", "Body", in.GetBody())
		return &npool.ContactViaEmailResponse{}, status.Error(codes.InvalidArgument, "Body is empty")
	}

	span = commontracer.TraceInvoker(span, "contact", "middleware", "ContactViaEmail")

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
		logger.Sugar().Errorw("ContactViaEmail", "err", err)
		return &npool.ContactViaEmailResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.ContactViaEmailResponse{}, nil
}
