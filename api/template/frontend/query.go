package frontend

import (
	"context"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	"github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	npoolpb "github.com/NpoolPlatform/message/npool"
	npool "github.com/NpoolPlatform/message/npool/notif/gw/v1/template/frontend"
	constant "github.com/NpoolPlatform/notif-gateway/pkg/message/const"
	commontracer "github.com/NpoolPlatform/notif-gateway/pkg/tracer"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	scodes "go.opentelemetry.io/otel/codes"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	mgrpb "github.com/NpoolPlatform/message/npool/notif/mgr/v1/template/frontend"
	"github.com/NpoolPlatform/notif-manager/pkg/client/template/frontend"
)

func (s *Server) GetFrontendTemplate(
	ctx context.Context,
	in *npool.GetFrontendTemplateRequest,
) (
	*npool.GetFrontendTemplateResponse,
	error,
) {
	var err error

	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "GetFrontendTemplate")
	defer span.End()

	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	span = commontracer.TraceInvoker(span, "contact", "manager", "GetFrontendTemplate")
	commontracer.TraceID(span, in.GetID())

	if _, err := uuid.Parse(in.GetID()); err != nil {
		logger.Sugar().Errorw("validate", "ID", in.GetID())
		return &npool.GetFrontendTemplateResponse{}, status.Error(codes.InvalidArgument, "ID is invalid")
	}

	info, err := frontend.GetFrontendTemplate(ctx, in.GetID())
	if err != nil {
		logger.Sugar().Errorw("validate", "err", err)
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
	var err error

	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "GetFrontendTemplates")
	defer span.End()

	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	span = commontracer.TraceInvoker(span, "contact", "manager", "GetFrontendTemplates")

	if _, err := uuid.Parse(in.GetAppID()); err != nil {
		logger.Sugar().Errorw("validate", "AppID", in.GetAppID())
		return &npool.GetFrontendTemplatesResponse{}, status.Error(codes.InvalidArgument, "AppID is invalid")
	}

	infos, total, err := frontend.GetFrontendTemplates(ctx, &mgrpb.Conds{
		AppID: &npoolpb.StringVal{
			Op:    cruder.EQ,
			Value: in.GetAppID(),
		},
	}, in.GetOffset(), in.GetLimit())
	if err != nil {
		logger.Sugar().Errorw("validate", "err", err)
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
	var err error

	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "GetAppFrontendTemplates")
	defer span.End()

	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	span = commontracer.TraceInvoker(span, "contact", "manager", "GetAppFrontendTemplates")
	commontracer.TraceOffsetLimit(span, int(in.GetOffset()), int(in.GetLimit()))

	if _, err := uuid.Parse(in.GetTargetAppID()); err != nil {
		logger.Sugar().Errorw("validate", "TargetAppID", in.GetTargetAppID())
		return &npool.GetAppFrontendTemplatesResponse{}, status.Error(codes.InvalidArgument, "TargetAppID is invalid")
	}

	infos, total, err := frontend.GetFrontendTemplates(ctx, &mgrpb.Conds{
		AppID: &npoolpb.StringVal{
			Op:    cruder.EQ,
			Value: in.GetTargetAppID(),
		},
	}, in.GetOffset(), in.GetLimit())
	if err != nil {
		logger.Sugar().Errorw("validate", "err", err)
		return &npool.GetAppFrontendTemplatesResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.GetAppFrontendTemplatesResponse{
		Infos: infos,
		Total: total,
	}, nil
}