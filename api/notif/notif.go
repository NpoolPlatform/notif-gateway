package notif

import (
	"context"

	mgrpb "github.com/NpoolPlatform/message/npool/notif/mgr/v1/notif"
	mgrcli "github.com/NpoolPlatform/notif-manager/pkg/client/notif"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"

	npool "github.com/NpoolPlatform/message/npool/notif/gw/v1/notif"
	constant "github.com/NpoolPlatform/notif-gateway/pkg/message/const"

	commontracer "github.com/NpoolPlatform/notif-gateway/pkg/tracer"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	scodes "go.opentelemetry.io/otel/codes"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	notif1 "github.com/NpoolPlatform/notif-gateway/pkg/notif"

	cruder "github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	commonpb "github.com/NpoolPlatform/message/npool"

	chanmgrpb "github.com/NpoolPlatform/message/npool/notif/mgr/v1/channel"
)

func (s *Server) GetNotif(ctx context.Context, in *npool.GetNotifRequest) (*npool.GetNotifResponse, error) {
	var err error

	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "GetNotif")
	defer span.End()

	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	span = commontracer.TraceID(span, in.GetID())

	_, err = uuid.Parse(in.GetID())
	if err != nil {
		logger.Sugar().Errorw("GetNotif", "ID", in.GetID(), "error", err)
		return &npool.GetNotifResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	info, err := notif1.GetNotif(ctx, in.GetID())
	if err != nil {
		logger.Sugar().Errorw("GetNotif", "error", err)
		return &npool.GetNotifResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.GetNotifResponse{
		Info: info,
	}, nil
}

func (s *Server) UpdateNotifs(ctx context.Context, in *npool.UpdateNotifsRequest) (*npool.UpdateNotifsResponse, error) {
	var err error

	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "GetNotif")
	defer span.End()

	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	for _, id := range in.GetIDs() {
		_, err = uuid.Parse(id)
		if err != nil {
			logger.Sugar().Errorw("GetNotif", "ID", id, "error", err)
			return &npool.UpdateNotifsResponse{}, status.Error(codes.InvalidArgument, err.Error())
		}
	}

	_, err = uuid.Parse(in.GetAppID())
	if err != nil {
		logger.Sugar().Errorw("GetNotif", "AppID", in.GetAppID(), "error", err)
		return &npool.UpdateNotifsResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	_, err = uuid.Parse(in.GetUserID())
	if err != nil {
		logger.Sugar().Errorw("GetNotif", "AppID", in.GetAppID(), "error", err)
		return &npool.UpdateNotifsResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	rows, _, err := mgrcli.GetNotifs(ctx, &mgrpb.Conds{
		IDs: &commonpb.StringSliceVal{Op: cruder.IN, Value: in.GetIDs()},
	}, 0, int32(len(in.GetIDs())))
	if err != nil {
		logger.Sugar().Errorw("GetNotif", "error", err)
		return &npool.UpdateNotifsResponse{}, status.Error(codes.Internal, err.Error())
	}

	for _, val := range rows {
		if val.AppID != in.GetAppID() || val.UserID != in.GetUserID() {
			logger.Sugar().Errorw("GetNotif", "error", err)
			return &npool.UpdateNotifsResponse{}, status.Error(codes.PermissionDenied, "permission denied")
		}
	}

	infos, err := notif1.UpdateNotifs(ctx, in.GetIDs(), in.GetNotified())
	if err != nil {
		logger.Sugar().Errorw("GetNotif", "error", err)
		return &npool.UpdateNotifsResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.UpdateNotifsResponse{
		Infos: infos,
	}, nil
}

func (s *Server) GetNotifs(ctx context.Context, in *npool.GetNotifsRequest) (*npool.GetNotifsResponse, error) {
	var err error

	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "GetNotifs")
	defer span.End()

	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	span = commontracer.TraceOffsetLimit(span, int(in.GetOffset()), int(in.GetLimit()))

	if _, err := uuid.Parse(in.GetAppID()); err != nil {
		logger.Sugar().Errorw("validate", "AppID", in.GetAppID(), "error", err)
		return &npool.GetNotifsResponse{}, status.Error(codes.Internal, "appID is invalid")
	}

	if _, err := uuid.Parse(in.GetUserID()); err != nil {
		logger.Sugar().Errorw("validate", "UserID", in.GetUserID(), "error", err)
		return &npool.GetNotifsResponse{}, status.Error(codes.Internal, "userID is invalid")
	}

	if _, err := uuid.Parse(in.GetLangID()); err != nil {
		logger.Sugar().Errorw("validate", "LangID", in.GetLangID(), "error", err)
		return &npool.GetNotifsResponse{}, status.Error(codes.Internal, "langID is invalid")
	}

	rows, total, err := notif1.GetNotifs(ctx, &mgrpb.Conds{
		AppID:   &commonpb.StringVal{Op: cruder.EQ, Value: in.GetAppID()},
		UserID:  &commonpb.StringVal{Op: cruder.EQ, Value: in.GetUserID()},
		LangID:  &commonpb.StringVal{Op: cruder.EQ, Value: in.GetLangID()},
		Channel: &commonpb.Uint32Val{Op: cruder.EQ, Value: uint32(chanmgrpb.NotifChannel_ChannelFrontend)},
	}, in.GetOffset(), in.GetLimit())
	if err != nil {
		logger.Sugar().Errorw("GetNotifs", "error", err)
		return &npool.GetNotifsResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.GetNotifsResponse{
		Infos: rows,
		Total: total,
	}, nil
}

func (s *Server) GetAppNotifs(ctx context.Context, in *npool.GetAppNotifsRequest) (*npool.GetAppNotifsResponse, error) {
	var err error

	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "GetAppNotifs")
	defer span.End()

	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	span = commontracer.TraceOffsetLimit(span, int(in.GetOffset()), int(in.GetLimit()))

	if _, err := uuid.Parse(in.GetAppID()); err != nil {
		logger.Sugar().Errorw("validate", "AppID", in.GetAppID(), "error", err)
		return &npool.GetAppNotifsResponse{}, status.Error(codes.Internal, "appID is invalid")
	}

	rows, total, err := notif1.GetNotifs(ctx, &mgrpb.Conds{
		AppID: &commonpb.StringVal{Op: cruder.EQ, Value: in.GetAppID()},
	}, in.GetOffset(), in.GetLimit())

	if err != nil {
		logger.Sugar().Errorw("GetAppNotifs", "error", err)
		return &npool.GetAppNotifsResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.GetAppNotifsResponse{
		Infos: rows,
		Total: total,
	}, nil
}
