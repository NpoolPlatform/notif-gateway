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

func (s *Server) UpdateFrontendTemplate(
	ctx context.Context,
	in *npool.UpdateFrontendTemplateRequest,
) (
	*npool.UpdateFrontendTemplateResponse,
	error,
) {
	handler, err := frontendtemplate1.NewHandler(
		ctx,
		frontendtemplate1.WithID(&in.ID),
		frontendtemplate1.WithAppID(&in.AppID),
		frontendtemplate1.WithLangID(in.TargetLangID),
		frontendtemplate1.WithTitle(in.Title),
		frontendtemplate1.WithContent(in.Content),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"UpdateFrontendTemplate",
			"In", in,
			"Error", err,
		)
		return &npool.UpdateFrontendTemplateResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	info, err := handler.UpdateFrontendTemplate(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"UpdateFrontendTemplate",
			"In", in,
			"Error", err,
		)
		return &npool.UpdateFrontendTemplateResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.UpdateFrontendTemplateResponse{
		Info: info,
	}, nil
}

func (s *Server) UpdateAppFrontendTemplate(
	ctx context.Context,
	in *npool.UpdateAppFrontendTemplateRequest,
) (
	*npool.UpdateAppFrontendTemplateResponse,
	error,
) {
	handler, err := frontendtemplate1.NewHandler(
		ctx,
		frontendtemplate1.WithID(&in.ID),
		frontendtemplate1.WithAppID(&in.TargetAppID),
		frontendtemplate1.WithLangID(in.TargetLangID),
		frontendtemplate1.WithTitle(in.Title),
		frontendtemplate1.WithContent(in.Content),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"UpdateAppFrontendTemplate",
			"In", in,
			"Error", err,
		)
		return &npool.UpdateAppFrontendTemplateResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	info, err := handler.UpdateFrontendTemplate(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"UpdateAppFrontendTemplate",
			"In", in,
			"Error", err,
		)
		return &npool.UpdateAppFrontendTemplateResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.UpdateAppFrontendTemplateResponse{
		Info: info,
	}, nil
}
