package frontend

import (
	"context"

	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"

	"github.com/google/uuid"

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

	applangmwcli "github.com/NpoolPlatform/g11n-middleware/pkg/client/applang"
	applangmgrpb "github.com/NpoolPlatform/message/npool/g11n/mgr/v1/applang"

	cruder "github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	commonpb "github.com/NpoolPlatform/message/npool"
)

//nolint:funlen,gocyclo
func (s *Server) UpdateFrontendTemplate(
	ctx context.Context,
	in *npool.UpdateFrontendTemplateRequest,
) (
	*npool.UpdateFrontendTemplateResponse,
	error,
) {
	var err error

	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "UpdateFrontendTemplate")
	defer span.End()

	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	if _, err := uuid.Parse(in.GetID()); err != nil {
		logger.Sugar().Errorw("validate", "ID", in.GetID())
		return &npool.UpdateFrontendTemplateResponse{}, status.Error(codes.InvalidArgument, "ID is invalid")
	}

	if _, err := uuid.Parse(in.GetTargetLangID()); err != nil {
		logger.Sugar().Errorw("validate", "ID", in.GetID())
		return &npool.UpdateFrontendTemplateResponse{}, status.Error(codes.InvalidArgument, "LangID is invalid")
	}

	if _, err := uuid.Parse(in.GetAppID()); err != nil {
		logger.Sugar().Errorw("validate", "ID", in.GetID())
		return &npool.UpdateFrontendTemplateResponse{}, status.Error(codes.InvalidArgument, "LangID is invalid")
	}

	if in.GetTitle() == "" && in.Title != nil {
		logger.Sugar().Errorw("validate", "Title", in.GetTitle())
		return &npool.UpdateFrontendTemplateResponse{}, status.Error(codes.InvalidArgument, "Title is empty")
	}

	if in.GetContent() == "" && in.Content != nil {
		logger.Sugar().Errorw("validate", "Content", in.GetContent())
		return &npool.UpdateFrontendTemplateResponse{}, status.Error(codes.InvalidArgument, "Content is empty")
	}

	if in.GetSender() == "" && in.Sender != nil {
		logger.Sugar().Errorw("validate", "Sender", in.GetSender())
		return &npool.UpdateFrontendTemplateResponse{}, status.Error(codes.InvalidArgument, "Sender is empty")
	}

	info, err := mgrcli.GetFrontendTemplate(ctx, in.GetID())
	if err != nil {
		logger.Sugar().Errorw("validate", "err", err)
		return &npool.UpdateFrontendTemplateResponse{}, status.Error(codes.Internal, err.Error())
	}

	if info.GetAppID() != in.GetAppID() {
		logger.Sugar().Errorw("validate", "err", err)
		return &npool.UpdateFrontendTemplateResponse{}, status.Error(codes.PermissionDenied, "permission denied")
	}

	if in.UsedFor != nil {
		switch in.GetUsedFor() {
		case basetypes.UsedFor_WithdrawalRequest:
		case basetypes.UsedFor_WithdrawalCompleted:
		case basetypes.UsedFor_DepositReceived:
		case basetypes.UsedFor_KYCApproved:
		case basetypes.UsedFor_KYCRejected:
		case basetypes.UsedFor_Announcement:
		default:
			logger.Sugar().Errorw("validate", "err", err)
			return &npool.UpdateFrontendTemplateResponse{}, status.Error(codes.InvalidArgument, "UsedFor is invalid")
		}
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
		return &npool.UpdateFrontendTemplateResponse{}, status.Error(codes.Internal, err.Error())
	}
	if appLang == nil {
		logger.Sugar().Errorw("validate", "err", err)
		return &npool.UpdateFrontendTemplateResponse{}, status.Error(codes.Internal, "AppLang not exist")
	}

	span = commontracer.TraceInvoker(span, "contact", "manager", "UpdateFrontendTemplate")

	info, err = mgrcli.UpdateFrontendTemplate(ctx, &mgrpb.FrontendTemplateReq{
		ID:      &in.ID,
		AppID:   &in.AppID,
		LangID:  in.TargetLangID,
		UsedFor: in.UsedFor,
		Title:   in.Title,
		Content: in.Content,
		Sender:  in.Sender,
	})

	if err != nil {
		logger.Sugar().Errorw("validate", "err", err)
		return &npool.UpdateFrontendTemplateResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.UpdateFrontendTemplateResponse{
		Info: info,
	}, nil
}

//nolint:gocyclo
func (s *Server) UpdateAppFrontendTemplate(
	ctx context.Context,
	in *npool.UpdateAppFrontendTemplateRequest,
) (
	*npool.UpdateAppFrontendTemplateResponse,
	error,
) {
	var err error

	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "UpdateFrontendTemplate")
	defer span.End()

	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	span = commontracer.TraceInvoker(span, "contact", "manager", "UpdateFrontendTemplate")

	if _, err := uuid.Parse(in.GetID()); err != nil {
		logger.Sugar().Errorw("validate", "ID", in.GetID())
		return &npool.UpdateAppFrontendTemplateResponse{}, status.Error(codes.InvalidArgument, "ID is invalid")
	}

	if _, err := uuid.Parse(in.GetTargetLangID()); err != nil {
		logger.Sugar().Errorw("validate", "ID", in.GetID())
		return &npool.UpdateAppFrontendTemplateResponse{}, status.Error(codes.InvalidArgument, "LangID is invalid")
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
		return &npool.UpdateAppFrontendTemplateResponse{}, status.Error(codes.Internal, err.Error())
	}
	if appLang == nil {
		logger.Sugar().Errorw("validate", "err", err)
		return &npool.UpdateAppFrontendTemplateResponse{}, status.Error(codes.Internal, "AppLang not exist")
	}

	info, err := mgrcli.GetFrontendTemplate(ctx, in.GetID())
	if err != nil {
		logger.Sugar().Errorw("validate", "err", err)
		return &npool.UpdateAppFrontendTemplateResponse{}, status.Error(codes.Internal, err.Error())
	}

	if info.GetAppID() != in.GetTargetAppID() {
		logger.Sugar().Errorw("validate", "err", err)
		return &npool.UpdateAppFrontendTemplateResponse{}, status.Error(codes.PermissionDenied, "permission denied")
	}

	if in.UsedFor != nil {
		switch in.GetUsedFor() {
		case basetypes.UsedFor_WithdrawalRequest:
		case basetypes.UsedFor_WithdrawalCompleted:
		case basetypes.UsedFor_DepositReceived:
		case basetypes.UsedFor_KYCApproved:
		case basetypes.UsedFor_KYCRejected:
		case basetypes.UsedFor_Announcement:
		default:
			logger.Sugar().Errorw("validate", "err", err)
			return &npool.UpdateAppFrontendTemplateResponse{}, status.Error(codes.InvalidArgument, "UsedFor is invalid")
		}
	}

	info, err = mgrcli.UpdateFrontendTemplate(ctx, &mgrpb.FrontendTemplateReq{
		ID:      &in.ID,
		AppID:   &in.TargetAppID,
		LangID:  in.TargetLangID,
		UsedFor: in.UsedFor,
		Title:   in.Title,
		Content: in.Content,
		Sender:  in.Sender,
	})

	if err != nil {
		logger.Sugar().Errorw("validate", "err", err)
		return &npool.UpdateAppFrontendTemplateResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.UpdateAppFrontendTemplateResponse{
		Info: info,
	}, nil
}
