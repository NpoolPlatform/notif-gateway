package email

import (
	"context"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	"github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	npoolpb "github.com/NpoolPlatform/message/npool"
	npool "github.com/NpoolPlatform/message/npool/notif/gw/v1/template/email"
	constant "github.com/NpoolPlatform/notif-gateway/pkg/message/const"
	commontracer "github.com/NpoolPlatform/notif-gateway/pkg/tracer"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	scodes "go.opentelemetry.io/otel/codes"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	mgrpb "github.com/NpoolPlatform/message/npool/notif/mgr/v1/template/email"
	"github.com/NpoolPlatform/notif-manager/pkg/client/template/email"
)

func (s *Server) GetEmailTemplate(ctx context.Context, in *npool.GetEmailTemplateRequest) (*npool.GetEmailTemplateResponse, error) {
	var err error

	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "GetEmailTemplate")
	defer span.End()

	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	span = commontracer.TraceInvoker(span, "contact", "manager", "GetEmailTemplate")
	commontracer.TraceID(span, in.GetID())

	if _, err := uuid.Parse(in.GetID()); err != nil {
		logger.Sugar().Errorw("validate", "ID", in.GetID())
		return &npool.GetEmailTemplateResponse{}, status.Error(codes.InvalidArgument, "ID is invalid")
	}

	info, err := email.GetEmailTemplate(ctx, in.GetID())
	if err != nil {
		logger.Sugar().Errorw("validate", "err", err)
		return &npool.GetEmailTemplateResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.GetEmailTemplateResponse{
		Info: info,
	}, nil
}

func (s *Server) GetEmailTemplates(ctx context.Context, in *npool.GetEmailTemplatesRequest) (*npool.GetEmailTemplatesResponse, error) {
	var err error

	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "GetEmailTemplates")
	defer span.End()

	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	span = commontracer.TraceInvoker(span, "contact", "manager", "GetEmailTemplates")

	if _, err := uuid.Parse(in.GetAppID()); err != nil {
		logger.Sugar().Errorw("validate", "AppID", in.GetAppID())
		return &npool.GetEmailTemplatesResponse{}, status.Error(codes.InvalidArgument, "AppID is invalid")
	}

	infos, total, err := email.GetEmailTemplates(ctx, &mgrpb.Conds{
		AppID: &npoolpb.StringVal{
			Op:    cruder.EQ,
			Value: in.GetAppID(),
		},
	}, in.GetOffset(), in.GetLimit())
	if err != nil {
		logger.Sugar().Errorw("validate", "err", err)
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
	var err error

	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "GetAppEmailTemplates")
	defer span.End()

	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	span = commontracer.TraceInvoker(span, "contact", "manager", "GetAppEmailTemplates")
	commontracer.TraceOffsetLimit(span, int(in.GetOffset()), int(in.GetLimit()))

	if _, err := uuid.Parse(in.GetTargetAppID()); err != nil {
		logger.Sugar().Errorw("validate", "TargetAppID", in.GetTargetAppID())
		return &npool.GetAppEmailTemplatesResponse{}, status.Error(codes.InvalidArgument, "TargetAppID is invalid")
	}

	infos, total, err := email.GetEmailTemplates(ctx, &mgrpb.Conds{
		AppID: &npoolpb.StringVal{
			Op:    cruder.EQ,
			Value: in.GetTargetAppID(),
		},
	}, in.GetOffset(), in.GetLimit())
	if err != nil {
		logger.Sugar().Errorw("validate", "err", err)
		return &npool.GetAppEmailTemplatesResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.GetAppEmailTemplatesResponse{
		Infos: infos,
		Total: total,
	}, nil
}
