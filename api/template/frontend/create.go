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

func (s *Server) CreateFrontendTemplate(
	ctx context.Context,
	in *npool.CreateFrontendTemplateRequest,
) (
	*npool.CreateFrontendTemplateResponse,
	error,
) {
	handler, err := frontendtemplate1.NewHandler(
		ctx,
		frontendtemplate1.WithAppID(&in.AppID),
		frontendtemplate1.WithLangID(&in.TargetLangID),
		frontendtemplate1.WithUsedFor(&in.UsedFor),
		frontendtemplate1.WithTitle(&in.Title),
		frontendtemplate1.WithContent(&in.Content),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"CreateFrontendTemplate",
			"In", in,
			"Error", err,
		)
		return &npool.CreateFrontendTemplateResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	info, err := handler.CreateFrontendTemplate(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"CreateFrontendTemplate",
			"In", in,
			"Error", err,
		)
		return &npool.CreateFrontendTemplateResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.CreateFrontendTemplateResponse{
		Info: info,
	}, nil
}

func (s *Server) CreateAppFrontendTemplate(
	ctx context.Context,
	in *npool.CreateAppFrontendTemplateRequest,
) (
	*npool.CreateAppFrontendTemplateResponse,
	error,
) {
	handler, err := frontendtemplate1.NewHandler(
		ctx,
		frontendtemplate1.WithAppID(&in.TargetAppID),
		frontendtemplate1.WithLangID(&in.TargetLangID),
		frontendtemplate1.WithUsedFor(&in.UsedFor),
		frontendtemplate1.WithTitle(&in.Title),
		frontendtemplate1.WithContent(&in.Content),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"CreateAppFrontendTemplate",
			"In", in,
			"Error", err,
		)
		return &npool.CreateAppFrontendTemplateResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	info, err := handler.CreateFrontendTemplate(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"CreateAppFrontendTemplate",
			"In", in,
			"Error", err,
		)
		return &npool.CreateAppFrontendTemplateResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.CreateAppFrontendTemplateResponse{
		Info: info,
	}, nil
}
