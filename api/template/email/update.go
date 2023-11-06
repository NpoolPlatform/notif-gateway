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

func (s *Server) UpdateEmailTemplate(
	ctx context.Context,
	in *npool.UpdateEmailTemplateRequest,
) (
	*npool.UpdateEmailTemplateResponse,
	error,
) {
	handler, err := emailtemplate1.NewHandler(
		ctx,
		emailtemplate1.WithID(&in.ID, true),
		emailtemplate1.WithEntID(&in.EntID, true),
		emailtemplate1.WithAppID(&in.AppID, true),
		emailtemplate1.WithSender(in.Sender, false),
		emailtemplate1.WithReplyTos(in.ReplyTos, false),
		emailtemplate1.WithCCTos(in.CCTos, false),
		emailtemplate1.WithSubject(in.Subject, false),
		emailtemplate1.WithBody(in.Body, false),
		emailtemplate1.WithDefaultToUsername(in.DefaultToUsername, false),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"UpdateEmailTemplate",
			"In", in,
			"Error", err,
		)
		return &npool.UpdateEmailTemplateResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	info, err := handler.UpdateEmailTemplate(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"UpdateEmailTemplate",
			"In", in,
			"Error", err,
		)
		return &npool.UpdateEmailTemplateResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.UpdateEmailTemplateResponse{
		Info: info,
	}, nil
}

func (s *Server) UpdateAppEmailTemplate(
	ctx context.Context,
	in *npool.UpdateAppEmailTemplateRequest,
) (
	*npool.UpdateAppEmailTemplateResponse,
	error,
) {
	handler, err := emailtemplate1.NewHandler(
		ctx,
		emailtemplate1.WithID(&in.ID, true),
		emailtemplate1.WithEntID(&in.EntID, true),
		emailtemplate1.WithAppID(&in.TargetAppID, true),
		emailtemplate1.WithSender(in.Sender, false),
		emailtemplate1.WithReplyTos(in.ReplyTos, false),
		emailtemplate1.WithCCTos(in.CCTos, false),
		emailtemplate1.WithSubject(in.Subject, false),
		emailtemplate1.WithBody(in.Body, false),
		emailtemplate1.WithDefaultToUsername(in.DefaultToUsername, false),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"UpdateAppEmailTemplate",
			"In", in,
			"Error", err,
		)
		return &npool.UpdateAppEmailTemplateResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	info, err := handler.UpdateEmailTemplate(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"UpdateAppEmailTemplate",
			"In", in,
			"Error", err,
		)
		return &npool.UpdateAppEmailTemplateResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.UpdateAppEmailTemplateResponse{
		Info: info,
	}, nil
}
