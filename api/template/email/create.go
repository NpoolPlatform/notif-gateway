//nolint:nolintlint,dupl
package email

import (
	"context"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	npool "github.com/NpoolPlatform/message/npool/notif/gw/v1/template/email"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	emailtemplate1 "github.com/NpoolPlatform/notif-gateway/pkg/template/email"
)

func (s *Server) CreateEmailTemplate(
	ctx context.Context,
	in *npool.CreateEmailTemplateRequest,
) (
	*npool.CreateEmailTemplateResponse,
	error,
) {
	handler, err := emailtemplate1.NewHandler(
		ctx,
		emailtemplate1.WithAppID(&in.AppID, true),
		emailtemplate1.WithLangID(&in.TargetLangID, true),
		emailtemplate1.WithUsedFor(&in.UsedFor, true),
		emailtemplate1.WithSubject(&in.Subject, true),
		emailtemplate1.WithDefaultToUsername(&in.DefaultToUsername, true),
		emailtemplate1.WithUsedFor(&in.UsedFor, true),
		emailtemplate1.WithSender(&in.Sender, true),
		emailtemplate1.WithReplyTos(in.ReplyTos, true),
		emailtemplate1.WithCCTos(in.CCTos, true),
		emailtemplate1.WithBody(&in.Body, true),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"CreateEmailTemplate",
			"In", in,
			"Error", err,
		)
		return &npool.CreateEmailTemplateResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	info, err := handler.CreateEmailTemplate(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"CreateEmailTemplate",
			"In", in,
			"Error", err,
		)
		return &npool.CreateEmailTemplateResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.CreateEmailTemplateResponse{
		Info: info,
	}, nil
}

func (s *Server) CreateAppEmailTemplate(
	ctx context.Context,
	in *npool.CreateAppEmailTemplateRequest,
) (
	*npool.CreateAppEmailTemplateResponse,
	error,
) {
	handler, err := emailtemplate1.NewHandler(
		ctx,
		emailtemplate1.WithAppID(&in.TargetAppID, true),
		emailtemplate1.WithLangID(&in.TargetLangID, true),
		emailtemplate1.WithUsedFor(&in.UsedFor, true),
		emailtemplate1.WithSubject(&in.Subject, true),
		emailtemplate1.WithDefaultToUsername(&in.DefaultToUsername, true),
		emailtemplate1.WithUsedFor(&in.UsedFor, true),
		emailtemplate1.WithSender(&in.Sender, true),
		emailtemplate1.WithReplyTos(in.ReplyTos, false),
		emailtemplate1.WithCCTos(in.CCTos, false),
		emailtemplate1.WithBody(&in.Body, true),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"CreateAppEmailTemplate",
			"In", in,
			"Error", err,
		)
		return &npool.CreateAppEmailTemplateResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	info, err := handler.CreateEmailTemplate(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"CreateAppEmailTemplate",
			"In", in,
			"Error", err,
		)
		return &npool.CreateAppEmailTemplateResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.CreateAppEmailTemplateResponse{
		Info: info,
	}, nil
}
