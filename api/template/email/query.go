package email

import (
	"context"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	npool "github.com/NpoolPlatform/message/npool/notif/gw/v1/template/email"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	emailtemplate1 "github.com/NpoolPlatform/notif-gateway/pkg/template/email"
)

func (s *Server) GetEmailTemplate(ctx context.Context, in *npool.GetEmailTemplateRequest) (*npool.GetEmailTemplateResponse, error) {
	handler, err := emailtemplate1.NewHandler(
		ctx,
		emailtemplate1.WithID(&in.ID),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"GetEmailTemplate",
			"In", in,
			"Error", err,
		)
		return &npool.GetEmailTemplateResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	info, err := handler.GetEmailTemplate(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"GetEmailTemplate",
			"In", in,
			"Error", err,
		)
		return &npool.GetEmailTemplateResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.GetEmailTemplateResponse{
		Info: info,
	}, nil
}

func (s *Server) GetEmailTemplates(ctx context.Context, in *npool.GetEmailTemplatesRequest) (*npool.GetEmailTemplatesResponse, error) {
	handler, err := emailtemplate1.NewHandler(
		ctx,
		emailtemplate1.WithAppID(&in.AppID),
		emailtemplate1.WithOffset(in.GetOffset()),
		emailtemplate1.WithLimit(in.GetLimit()),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"GetEmailTemplates",
			"In", in,
			"Error", err,
		)
		return &npool.GetEmailTemplatesResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	infos, total, err := handler.GetEmailTemplates(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"GetEmailTemplates",
			"In", in,
			"Error", err,
		)
		return &npool.GetEmailTemplatesResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.GetEmailTemplatesResponse{
		Infos: infos,
		Total: total,
	}, nil
}

func (s *Server) GetAppEmailTemplates(
	ctx context.Context,
	in *npool.GetAppEmailTemplatesRequest,
) (
	*npool.GetAppEmailTemplatesResponse,
	error,
) {
	handler, err := emailtemplate1.NewHandler(
		ctx,
		emailtemplate1.WithAppID(&in.TargetAppID),
		emailtemplate1.WithOffset(in.GetOffset()),
		emailtemplate1.WithLimit(in.GetLimit()),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"GetAppEmailTemplates",
			"In", in,
			"Error", err,
		)
		return &npool.GetAppEmailTemplatesResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	infos, total, err := handler.GetEmailTemplates(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"GetAppEmailTemplates",
			"In", in,
			"Error", err,
		)
		return &npool.GetAppEmailTemplatesResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.GetAppEmailTemplatesResponse{
		Infos: infos,
		Total: total,
	}, nil
}
