package sms

import (
	"context"

	"github.com/google/uuid"

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

	applangmwcli "github.com/NpoolPlatform/g11n-middleware/pkg/client/applang"
	applangmwpb "github.com/NpoolPlatform/message/npool/g11n/mw/v1/applang"

	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"

	cruder "github.com/NpoolPlatform/libent-cruder/pkg/cruder"
)

func (s *Server) UpdateSMSTemplate(
	ctx context.Context,
	in *npool.UpdateSMSTemplateRequest,
) (
	*npool.UpdateSMSTemplateResponse,
	error,
) {
	var err error

	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "UpdateSMSTemplate")
	defer span.End()

	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	if _, err := uuid.Parse(in.GetID()); err != nil {
		logger.Sugar().Errorw("validate", "ID", in.GetID())
		return &npool.UpdateSMSTemplateResponse{}, status.Error(codes.InvalidArgument, "ID is invalid")
	}

	info, err := mgrcli.GetSMSTemplate(ctx, in.GetID())
	if err != nil {
		logger.Sugar().Errorw("validate", "err", err)
		return &npool.UpdateSMSTemplateResponse{}, status.Error(codes.Internal, err.Error())
	}

	if info.GetAppID() != in.GetAppID() {
		logger.Sugar().Errorw("validate", "err", err)
		return &npool.UpdateSMSTemplateResponse{}, status.Error(codes.PermissionDenied, "permission denied")
	}

	if _, err := uuid.Parse(in.GetTargetLangID()); err != nil {
		logger.Sugar().Errorw("validate", "ID", in.GetID())
		return &npool.UpdateSMSTemplateResponse{}, status.Error(codes.InvalidArgument, "LangID is invalid")
	}

	if in.GetSubject() == "" {
		logger.Sugar().Errorw("validate", "Subject", in.GetSubject())
		return &npool.UpdateSMSTemplateResponse{}, status.Error(codes.InvalidArgument, "Subject is empty")
	}

	appLang, err := applangmwcli.GetLangOnly(ctx, &applangmwpb.Conds{
		AppID: &basetypes.StringVal{
			Op:    cruder.EQ,
			Value: in.GetAppID(),
		},
		LangID: &basetypes.StringVal{
			Op:    cruder.EQ,
			Value: in.GetTargetLangID(),
		},
	})
	if err != nil {
		logger.Sugar().Errorw("validate", "err", err)
		return &npool.UpdateSMSTemplateResponse{}, status.Error(codes.Internal, err.Error())
	}
	if appLang == nil {
		logger.Sugar().Errorw("validate", "err", err)
		return &npool.UpdateSMSTemplateResponse{}, status.Error(codes.Internal, "Lang not exist")
	}

	span = commontracer.TraceInvoker(span, "contact", "manager", "UpdateSMSTemplate")

	info, err = mgrcli.UpdateSMSTemplate(ctx, &mgrpb.SMSTemplateReq{
		ID:      &in.ID,
		LangID:  &in.TargetLangID,
		Subject: in.Subject,
		Message: in.Message,
	})

	if err != nil {
		logger.Sugar().Errorw("validate", "err", err)
		return &npool.UpdateSMSTemplateResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.UpdateSMSTemplateResponse{
		Info: info,
	}, nil
}

func (s *Server) UpdateAppSMSTemplate(
	ctx context.Context,
	in *npool.UpdateAppSMSTemplateRequest,
) (
	*npool.UpdateAppSMSTemplateResponse,
	error,
) {
	var err error

	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "UpdateAppSMSTemplate")
	defer span.End()

	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	if _, err := uuid.Parse(in.GetID()); err != nil {
		logger.Sugar().Errorw("validate", "ID", in.GetID())
		return &npool.UpdateAppSMSTemplateResponse{}, status.Error(codes.InvalidArgument, "ID is invalid")
	}

	if _, err := uuid.Parse(in.GetTargetLangID()); err != nil {
		logger.Sugar().Errorw("validate", "ID", in.GetID())
		return &npool.UpdateAppSMSTemplateResponse{}, status.Error(codes.InvalidArgument, "LangID is invalid")
	}

	if in.GetSubject() == "" {
		logger.Sugar().Errorw("validate", "Subject", in.GetSubject())
		return &npool.UpdateAppSMSTemplateResponse{}, status.Error(codes.InvalidArgument, "Subject is empty")
	}

	appLang, err := applangmwcli.GetLangOnly(ctx, &applangmwpb.Conds{
		AppID: &basetypes.StringVal{
			Op:    cruder.EQ,
			Value: in.GetTargetAppID(),
		},
		LangID: &basetypes.StringVal{
			Op:    cruder.EQ,
			Value: in.GetTargetLangID(),
		},
	})
	if err != nil {
		logger.Sugar().Errorw("validate", "err", err)
		return &npool.UpdateAppSMSTemplateResponse{}, status.Error(codes.Internal, err.Error())
	}
	if appLang == nil {
		logger.Sugar().Errorw("validate", "err", err)
		return &npool.UpdateAppSMSTemplateResponse{}, status.Error(codes.Internal, "AppLang not exist")
	}

	span = commontracer.TraceInvoker(span, "contact", "manager", "UpdateSMSTemplate")

	info, err := mgrcli.UpdateSMSTemplate(ctx, &mgrpb.SMSTemplateReq{
		ID:      &in.ID,
		LangID:  &in.TargetLangID,
		Subject: in.Subject,
		Message: in.Message,
	})

	if err != nil {
		logger.Sugar().Errorw("validate", "err", err)
		return &npool.UpdateAppSMSTemplateResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.UpdateAppSMSTemplateResponse{
		Info: info,
	}, nil
}
