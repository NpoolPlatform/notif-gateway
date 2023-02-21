//nolint:dupl
package sms

import (
	"context"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	"github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	npoolpb "github.com/NpoolPlatform/message/npool"
	npool "github.com/NpoolPlatform/message/npool/notif/gw/v1/template/sms"
	constant "github.com/NpoolPlatform/notif-gateway/pkg/message/const"
	commontracer "github.com/NpoolPlatform/notif-gateway/pkg/tracer"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	scodes "go.opentelemetry.io/otel/codes"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	mgrpb "github.com/NpoolPlatform/message/npool/notif/mgr/v1/template/sms"
	"github.com/NpoolPlatform/notif-manager/pkg/client/template/sms"
)

func (s *Server) GetSMSTemplate(ctx context.Context, in *npool.GetSMSTemplateRequest) (*npool.GetSMSTemplateResponse, error) {
	var err error

	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "GetContact")
	defer span.End()

	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	span = commontracer.TraceInvoker(span, "contact", "manager", "GetContact")

	if _, err := uuid.Parse(in.GetID()); err != nil {
		logger.Sugar().Errorw("validate", "ID", in.GetID())
		return &npool.GetSMSTemplateResponse{}, status.Error(codes.InvalidArgument, "ID is invalid")
	}

	info, err := sms.GetSMSTemplate(ctx, in.GetID())
	if err != nil {
		logger.Sugar().Errorw("validate", "err", err)
		return &npool.GetSMSTemplateResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.GetSMSTemplateResponse{
		Info: info,
	}, nil
}

func (s *Server) GetSMSTemplates(ctx context.Context, in *npool.GetSMSTemplatesRequest) (*npool.GetSMSTemplatesResponse, error) {
	var err error

	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "GetContact")
	defer span.End()

	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	span = commontracer.TraceInvoker(span, "contact", "manager", "GetContact")

	if _, err := uuid.Parse(in.GetAppID()); err != nil {
		logger.Sugar().Errorw("validate", "AppID", in.GetAppID())
		return &npool.GetSMSTemplatesResponse{}, status.Error(codes.InvalidArgument, "AppID is invalid")
	}

	infos, total, err := sms.GetSMSTemplates(ctx, &mgrpb.Conds{
		AppID: &npoolpb.StringVal{
			Op:    cruder.EQ,
			Value: in.GetAppID(),
		},
	}, in.GetOffset(), in.GetLimit())
	if err != nil {
		logger.Sugar().Errorw("validate", "err", err)
		return &npool.GetSMSTemplatesResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.GetSMSTemplatesResponse{
		Infos: infos,
		Total: total,
	}, nil
}

func (s *Server) GetAppSMSTemplates(ctx context.Context, in *npool.GetAppSMSTemplatesRequest) (*npool.GetAppSMSTemplatesResponse, error) {
	var err error

	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "GetContact")
	defer span.End()

	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	span = commontracer.TraceInvoker(span, "contact", "manager", "GetContact")

	if _, err := uuid.Parse(in.GetTargetAppID()); err != nil {
		logger.Sugar().Errorw("validate", "TargetAppID", in.GetTargetAppID())
		return &npool.GetAppSMSTemplatesResponse{}, status.Error(codes.InvalidArgument, "TargetAppID is invalid")
	}

	infos, total, err := sms.GetSMSTemplates(ctx, &mgrpb.Conds{
		AppID: &npoolpb.StringVal{
			Op:    cruder.EQ,
			Value: in.GetTargetAppID(),
		},
	}, in.GetOffset(), in.GetLimit())
	if err != nil {
		logger.Sugar().Errorw("validate", "err", err)
		return &npool.GetAppSMSTemplatesResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.GetAppSMSTemplatesResponse{
		Infos: infos,
		Total: total,
	}, nil
}
