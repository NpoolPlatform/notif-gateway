package frontend

import (
	"context"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	npool "github.com/NpoolPlatform/message/npool/notif/gw/v1/template/frontend"
	constant "github.com/NpoolPlatform/notif-gateway/pkg/message/const"
	commontracer "github.com/NpoolPlatform/notif-gateway/pkg/tracer"
	"go.opentelemetry.io/otel"
	scodes "go.opentelemetry.io/otel/codes"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	mgrpb "github.com/NpoolPlatform/message/npool/notif/mgr/v1/template/frontend"
	mgrcli "github.com/NpoolPlatform/notif-manager/pkg/client/template/frontend"
)

func (s *Server) CreateFrontendTemplate(
	ctx context.Context,
	in *npool.CreateFrontendTemplateRequest,
) (
	*npool.CreateFrontendTemplateResponse,
	error,
) {
	var err error

	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "CreateContact")
	defer span.End()

	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	err = validate(ctx, in)
	if err != nil {
		return nil, err
	}

	span = commontracer.TraceInvoker(span, "contact", "manager", "CreateFrontendTemplate")

	info, err := mgrcli.CreateFrontendTemplate(ctx, &mgrpb.FrontendTemplateReq{
		AppID:   &in.AppID,
		LangID:  &in.TargetLangID,
		UsedFor: &in.UsedFor,
		Title:   &in.Title,
		Content: &in.Content,
	})

	if err != nil {
		logger.Sugar().Errorw("validate", "err", err)
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
	var err error

	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "CreateContact")
	defer span.End()

	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	span = commontracer.TraceInvoker(span, "contact", "manager", "CreateFrontendTemplate")

	err = validate(ctx, &npool.CreateFrontendTemplateRequest{
		AppID:        in.TargetAppID,
		TargetLangID: in.TargetLangID,
		UsedFor:      in.UsedFor,
		Title:        in.Title,
		Content:      in.Content,
	})
	if err != nil {
		return nil, err
	}

	info, err := mgrcli.CreateFrontendTemplate(ctx, &mgrpb.FrontendTemplateReq{
		AppID:   &in.TargetAppID,
		LangID:  &in.TargetLangID,
		UsedFor: &in.UsedFor,
		Title:   &in.Title,
		Content: &in.Content,
	})

	if err != nil {
		logger.Sugar().Errorw("validate", "err", err)
		return &npool.CreateAppFrontendTemplateResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.CreateAppFrontendTemplateResponse{
		Info: info,
	}, nil
}
