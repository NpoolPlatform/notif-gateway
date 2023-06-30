//nolint:nolintlint,dupl
package frontend

import (
	"context"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	npool "github.com/NpoolPlatform/message/npool/notif/gw/v1/template/frontend"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	frontendtemplate1 "github.com/NpoolPlatform/notif-gateway/pkg/template/frontend"
)

func (s *Server) GetFrontendTemplate(
	ctx context.Context,
	in *npool.GetFrontendTemplateRequest,
) (
	*npool.GetFrontendTemplateResponse,
	error,
) {
	handler, err := frontendtemplate1.NewHandler(
		ctx,
		frontendtemplate1.WithID(&in.ID),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"GetFrontendTemplate",
			"In", in,
			"Error", err,
		)
		return &npool.GetFrontendTemplateResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	info, err := handler.GetFrontendTemplate(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"GetFrontendTemplate",
			"In", in,
			"Error", err,
		)
		return &npool.GetFrontendTemplateResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.GetFrontendTemplateResponse{
		Info: info,
	}, nil
}

func (s *Server) GetFrontendTemplates(
	ctx context.Context,
	in *npool.GetFrontendTemplatesRequest,
) (
	*npool.GetFrontendTemplatesResponse,
	error,
) {
	handler, err := frontendtemplate1.NewHandler(
		ctx,
		frontendtemplate1.WithAppID(&in.AppID),
		frontendtemplate1.WithOffset(in.GetOffset()),
		frontendtemplate1.WithLimit(in.GetLimit()),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"GetFrontendTemplates",
			"In", in,
			"Error", err,
		)
		return &npool.GetFrontendTemplatesResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	infos, total, err := handler.GetFrontendTemplates(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"GetFrontendTemplates",
			"In", in,
			"Error", err,
		)
		return &npool.GetFrontendTemplatesResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.GetFrontendTemplatesResponse{
		Infos: infos,
		Total: total,
	}, nil
}

func (s *Server) GetAppFrontendTemplates(
	ctx context.Context,
	in *npool.GetAppFrontendTemplatesRequest,
) (
	*npool.GetAppFrontendTemplatesResponse,
	error,
) {
	handler, err := frontendtemplate1.NewHandler(
		ctx,
		frontendtemplate1.WithAppID(&in.TargetAppID),
		frontendtemplate1.WithOffset(in.GetOffset()),
		frontendtemplate1.WithLimit(in.GetLimit()),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"GetAppFrontendTemplates",
			"In", in,
			"Error", err,
		)
		return &npool.GetAppFrontendTemplatesResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	infos, total, err := handler.GetFrontendTemplates(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"GetAppFrontendTemplates",
			"In", in,
			"Error", err,
		)
		return &npool.GetAppFrontendTemplatesResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.GetAppFrontendTemplatesResponse{
		Infos: infos,
		Total: total,
	}, nil
}
