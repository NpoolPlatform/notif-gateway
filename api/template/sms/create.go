package sms

import (
	"context"

	tracer "github.com/NpoolPlatform/notif-manager/pkg/tracer/template/sms"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	npool "github.com/NpoolPlatform/message/npool/notif/gw/v1/template/sms"
	constant "github.com/NpoolPlatform/notif-gateway/pkg/message/const"
	commontracer "github.com/NpoolPlatform/notif-gateway/pkg/tracer"
	"go.opentelemetry.io/otel"
	scodes "go.opentelemetry.io/otel/codes"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	mgrpb "github.com/NpoolPlatform/message/npool/notif/mgr/v1/template/sms"
	mgrcli "github.com/NpoolPlatform/notif-manager/pkg/client/template/sms"
)

func (s *Server) CreateSMSTemplate(
	ctx context.Context,
	in *npool.CreateSMSTemplateRequest,
) (
	*npool.CreateSMSTemplateResponse,
	error,
) {
	var err error

	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "CreateSMSTemplate")
	defer span.End()

	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	contactInfo := &mgrpb.SMSTemplateReq{
		AppID:   &in.AppID,
		LangID:  &in.TargetLangID,
		UsedFor: &in.UsedFor,
		Subject: &in.Subject,
		Message: &in.Message,
	}

	tracer.Trace(span, contactInfo)

	err = validate(ctx, &npool.CreateSMSTemplateRequest{
		AppID:        in.AppID,
		TargetLangID: in.TargetLangID,
		UsedFor:      in.UsedFor,
		Subject:      in.Subject,
		Message:      in.Message,
	})
	if err != nil {
		return nil, err
	}

	span = commontracer.TraceInvoker(span, "contact", "manager", "CreateSMSTemplate")

	info, err := mgrcli.CreateSMSTemplate(ctx, contactInfo)

	if err != nil {
		logger.Sugar().Errorw("validate", "err", err)
		return &npool.CreateSMSTemplateResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.CreateSMSTemplateResponse{
		Info: info,
	}, nil
}

func (s *Server) CreateAppSMSTemplate(
	ctx context.Context,
	in *npool.CreateAppSMSTemplateRequest,
) (
	*npool.CreateAppSMSTemplateResponse,
	error,
) {
	var err error

	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "CreateAppSMSTemplate")
	defer span.End()

	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	contactInfo := &mgrpb.SMSTemplateReq{
		AppID:   &in.TargetAppID,
		LangID:  &in.TargetLangID,
		UsedFor: &in.UsedFor,
		Subject: &in.Subject,
		Message: &in.Message,
	}

	tracer.Trace(span, contactInfo)

	span = commontracer.TraceInvoker(span, "contact", "manager", "CreateAppSMSTemplate")

	err = validate(ctx, &npool.CreateSMSTemplateRequest{
		AppID:        in.TargetAppID,
		TargetLangID: in.TargetLangID,
		UsedFor:      in.UsedFor,
		Subject:      in.Subject,
		Message:      in.Message,
	})
	if err != nil {
		return nil, err
	}

	info, err := mgrcli.CreateSMSTemplate(ctx, contactInfo)

	if err != nil {
		logger.Sugar().Errorw("validate", "err", err)
		return &npool.CreateAppSMSTemplateResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.CreateAppSMSTemplateResponse{
		Info: info,
	}, nil
}
