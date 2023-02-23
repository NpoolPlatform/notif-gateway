package email

import (
	"context"

	"github.com/google/uuid"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	npool "github.com/NpoolPlatform/message/npool/notif/gw/v1/template/email"
	constant "github.com/NpoolPlatform/notif-gateway/pkg/message/const"
	commontracer "github.com/NpoolPlatform/notif-gateway/pkg/tracer"

	"go.opentelemetry.io/otel"
	scodes "go.opentelemetry.io/otel/codes"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	mgrpb "github.com/NpoolPlatform/message/npool/notif/mgr/v1/template/email"
	mgrcli "github.com/NpoolPlatform/notif-manager/pkg/client/template/email"

	applangmwcli "github.com/NpoolPlatform/g11n-middleware/pkg/client/applang"
	applangmgrpb "github.com/NpoolPlatform/message/npool/g11n/mgr/v1/applang"

	cruder "github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	commonpb "github.com/NpoolPlatform/message/npool"
)

func (s *Server) UpdateEmailTemplate(
	ctx context.Context,
	in *npool.UpdateEmailTemplateRequest,
) (
	*npool.UpdateEmailTemplateResponse,
	error,
) {
	var err error

	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "UpdateEmailTemplate")
	defer span.End()

	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	if _, err := uuid.Parse(in.GetID()); err != nil {
		logger.Sugar().Errorw("validate", "ID", in.GetID())
		return &npool.UpdateEmailTemplateResponse{}, status.Error(codes.InvalidArgument, "ID is invalid")
	}

	if _, err := uuid.Parse(in.GetTargetLangID()); err != nil {
		logger.Sugar().Errorw("validate", "ID", in.GetID())
		return &npool.UpdateEmailTemplateResponse{}, status.Error(codes.InvalidArgument, "LangID is invalid")
	}

	if in.GetSubject() == "" {
		logger.Sugar().Errorw("validate", "Subject", in.GetSubject())
		return &npool.UpdateEmailTemplateResponse{}, status.Error(codes.InvalidArgument, "Subject is empty")
	}

	info, err := mgrcli.GetEmailTemplate(ctx, in.GetID())
	if err != nil {
		logger.Sugar().Errorw("validate", "err", err)
		return &npool.UpdateEmailTemplateResponse{}, status.Error(codes.Internal, err.Error())
	}

	if info.GetAppID() != in.GetAppID() {
		logger.Sugar().Errorw("validate", "err", err)
		return &npool.UpdateEmailTemplateResponse{}, status.Error(codes.PermissionDenied, "permission denied")
	}

	appLang, err := applangmwcli.GetLangOnly(ctx, &applangmgrpb.Conds{
		AppID: &commonpb.StringVal{
			Op:    cruder.EQ,
			Value: in.GetAppID(),
		},
		LangID: &commonpb.StringVal{
			Op:    cruder.EQ,
			Value: in.GetTargetLangID(),
		},
	})
	if err != nil {
		logger.Sugar().Errorw("validate", "err", err)
		return &npool.UpdateEmailTemplateResponse{}, status.Error(codes.Internal, err.Error())
	}
	if appLang == nil {
		logger.Sugar().Errorw("validate", "err", err)
		return &npool.UpdateEmailTemplateResponse{}, status.Error(codes.Internal, "AppLang not exist")
	}

	span = commontracer.TraceInvoker(span, "contact", "manager", "UpdateEmailTemplate")

	info, err = mgrcli.UpdateEmailTemplate(ctx, &mgrpb.EmailTemplateReq{
		ID:                &in.ID,
		LangID:            &in.TargetLangID,
		Sender:            in.Sender,
		ReplyTos:          in.ReplyTos,
		CCTos:             in.CCTos,
		Subject:           in.Subject,
		Body:              in.Body,
		DefaultToUsername: in.DefaultToUsername,
	})

	if err != nil {
		logger.Sugar().Errorw("validate", "err", err)
		return &npool.UpdateEmailTemplateResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.UpdateEmailTemplateResponse{
		Info: info,
	}, nil
}

func (s *Server) UpdateAppEmailTemplate(
	ctx context.Context,
	in *npool.UpdateAppEmailTemplateRequest,
) (
	*npool.UpdateAppEmailTemplateResponse,
	error,
) {
	var err error

	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "UpdateEmailTemplate")
	defer span.End()

	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	span = commontracer.TraceInvoker(span, "contact", "manager", "UpdateEmailTemplate")

	if _, err := uuid.Parse(in.GetID()); err != nil {
		logger.Sugar().Errorw("validate", "ID", in.GetID())
		return &npool.UpdateAppEmailTemplateResponse{}, status.Error(codes.InvalidArgument, "ID is invalid")
	}

	if _, err := uuid.Parse(in.GetTargetLangID()); err != nil {
		logger.Sugar().Errorw("validate", "ID", in.GetID())
		return &npool.UpdateAppEmailTemplateResponse{}, status.Error(codes.InvalidArgument, "LangID is invalid")
	}

	if in.GetSubject() == "" {
		logger.Sugar().Errorw("validate", "Subject", in.GetSubject())
		return &npool.UpdateAppEmailTemplateResponse{}, status.Error(codes.InvalidArgument, "Subject is empty")
	}

	appLang, err := applangmwcli.GetLangOnly(ctx, &applangmgrpb.Conds{
		AppID: &commonpb.StringVal{
			Op:    cruder.EQ,
			Value: in.GetTargetAppID(),
		},
		LangID: &commonpb.StringVal{
			Op:    cruder.EQ,
			Value: in.GetTargetLangID(),
		},
	})
	if err != nil {
		logger.Sugar().Errorw("validate", "err", err)
		return &npool.UpdateAppEmailTemplateResponse{}, status.Error(codes.Internal, err.Error())
	}
	if appLang == nil {
		logger.Sugar().Errorw("validate", "err", err)
		return &npool.UpdateAppEmailTemplateResponse{}, status.Error(codes.Internal, "AppLang not exist")
	}

	info, err := mgrcli.UpdateEmailTemplate(ctx, &mgrpb.EmailTemplateReq{
		ID:                &in.ID,
		LangID:            &in.TargetLangID,
		Sender:            in.Sender,
		ReplyTos:          in.ReplyTos,
		CCTos:             in.CCTos,
		Subject:           in.Subject,
		Body:              in.Body,
		DefaultToUsername: in.DefaultToUsername,
	})

	if err != nil {
		logger.Sugar().Errorw("validate", "err", err)
		return &npool.UpdateAppEmailTemplateResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.UpdateAppEmailTemplateResponse{
		Info: info,
	}, nil
}
