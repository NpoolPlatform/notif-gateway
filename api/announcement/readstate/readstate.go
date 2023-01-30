//nolint:dupl
package readstate

import (
	"context"

	usercli "github.com/NpoolPlatform/appuser-middleware/pkg/client/user"
	mgrcli "github.com/NpoolPlatform/notif-manager/pkg/client/announcement"

	commontracer "github.com/NpoolPlatform/notif-gateway/pkg/tracer"

	constant "github.com/NpoolPlatform/notif-gateway/pkg/message/const"

	"go.opentelemetry.io/otel"
	scodes "go.opentelemetry.io/otel/codes"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	npool "github.com/NpoolPlatform/message/npool/notif/gw/v1/announcement/readstate"

	"github.com/google/uuid"

	readstate1 "github.com/NpoolPlatform/notif-gateway/pkg/announcement/readstate"
)

func (s *Server) CreateReadState(
	ctx context.Context,
	in *npool.CreateReadStateRequest,
) (
	*npool.CreateReadStateResponse,
	error,
) {
	var err error

	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "CreateReadState")
	defer span.End()

	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	_, err = uuid.Parse(in.GetAppID())
	if err != nil {
		logger.Sugar().Errorw("CreateReadState", "AppID", in.GetAppID(), "error", err)
		return &npool.CreateReadStateResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}
	_, err = uuid.Parse(in.GetUserID())
	if err != nil {
		logger.Sugar().Errorw("CreateReadState", "UserID", in.GetUserID(), "error", err)
		return &npool.CreateReadStateResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}
	_, err = uuid.Parse(in.GetAnnouncementID())
	if err != nil {
		logger.Sugar().Errorw("CreateReadState", "AnnouncementID", in.GetAnnouncementID(), "error", err)
		return &npool.CreateReadStateResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	user, err := usercli.GetUser(ctx, in.GetAppID(), in.GetUserID())
	if err != nil {
		logger.Sugar().Errorw("CreateReadState", "error", err)
		return &npool.CreateReadStateResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}
	if user == nil {
		logger.Sugar().Errorw(
			"CreateReadState",
			"AppID",
			in.GetAppID(),
			"User",
			in.GetUserID(),
			"error",
			"app user not exist",
		)
		return &npool.CreateReadStateResponse{}, status.Error(codes.InvalidArgument, "app user not exist")
	}

	exist, err := mgrcli.ExistAnnouncement(ctx, in.GetAnnouncementID())
	if err != nil {
		logger.Sugar().Errorw("CreateReadState", "error", err)
		return &npool.CreateReadStateResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	if !exist {
		logger.Sugar().Errorw(
			"CreateReadState",
			"AppID",
			in.GetAppID(),
			"AnnouncementID",
			in.GetAnnouncementID(),
			"error",
			"announcement not exist",
		)
		return &npool.CreateReadStateResponse{}, status.Error(codes.InvalidArgument, "announcement not exist")
	}

	info, err := readstate1.CreateReadState(
		ctx,
		in.GetAppID(),
		in.GetUserID(),
		in.GetAnnouncementID(),
	)
	if err != nil {
		logger.Sugar().Errorw("CreateReadState", "error", err)
		return &npool.CreateReadStateResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.CreateReadStateResponse{
		Info: info,
	}, nil
}

func (s *Server) GetReadState(
	ctx context.Context,
	in *npool.GetReadStateRequest,
) (
	*npool.GetReadStateResponse,
	error,
) {
	var err error

	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "GetReadState")
	defer span.End()

	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	_, err = uuid.Parse(in.GetAppID())
	if err != nil {
		logger.Sugar().Errorw("GetReadState", "AppID", in.GetAppID(), "error", err)
		return &npool.GetReadStateResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}
	_, err = uuid.Parse(in.GetUserID())
	if err != nil {
		logger.Sugar().Errorw("GetReadState", "UserID", in.GetUserID(), "error", err)
		return &npool.GetReadStateResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}
	_, err = uuid.Parse(in.GetAnnouncementID())
	if err != nil {
		logger.Sugar().Errorw("GetReadState", "AnnouncementID", in.GetAnnouncementID(), "error", err)
		return &npool.GetReadStateResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	span = commontracer.TraceInvoker(span, "announcement/readstate", "crud", "Row")

	info, err := readstate1.GetReadState(ctx, in.GetAppID(), in.GetUserID(), in.GetAnnouncementID())
	if err != nil {
		logger.Sugar().Errorw("GetReadState", "error", err)
		return &npool.GetReadStateResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.GetReadStateResponse{
		Info: info,
	}, nil
}

func (s *Server) GetReadStates(
	ctx context.Context,
	in *npool.GetReadStatesRequest,
) (
	*npool.GetReadStatesResponse,
	error,
) {
	var err error

	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "GetReadStates")
	defer span.End()

	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	_, err = uuid.Parse(in.GetAppID())
	if err != nil {
		logger.Sugar().Errorw("GetReadStates", "AppID", in.GetAppID(), "error", err)
		return &npool.GetReadStatesResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}
	_, err = uuid.Parse(in.GetUserID())
	if err != nil {
		logger.Sugar().Errorw("GetReadStates", "UserID", in.GetUserID(), "error", err)
		return &npool.GetReadStatesResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	infos, total, err := readstate1.GetReadStates(ctx, in.GetAppID(), &in.UserID, in.GetOffset(), in.GetLimit())
	if err != nil {
		logger.Sugar().Errorw("GetReadStates", "error", err)
		return &npool.GetReadStatesResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.GetReadStatesResponse{
		Infos: infos,
		Total: total,
	}, nil
}

func (s *Server) GetAppUserReadStates(
	ctx context.Context,
	in *npool.GetAppUserReadStatesRequest,
) (
	*npool.GetAppUserReadStatesResponse,
	error,
) {
	var err error

	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "GetAppUserReadStates")
	defer span.End()

	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	_, err = uuid.Parse(in.GetTargetAppID())
	if err != nil {
		logger.Sugar().Errorw("GetAppUserReadStates", "TargetAppID", in.GetTargetAppID(), "error", err)
		return &npool.GetAppUserReadStatesResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}
	_, err = uuid.Parse(in.GetTargetUserID())
	if err != nil {
		logger.Sugar().Errorw("GetAppUserReadStates", "TargetUserID", in.GetTargetUserID(), "error", err)
		return &npool.GetAppUserReadStatesResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	infos, total, err := readstate1.GetReadStates(ctx, in.GetTargetAppID(), &in.TargetUserID, in.GetOffset(), in.GetLimit())
	if err != nil {
		logger.Sugar().Errorw("GetAppUserReadStates", "error", err)
		return &npool.GetAppUserReadStatesResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.GetAppUserReadStatesResponse{
		Infos: infos,
		Total: total,
	}, nil
}

func (s *Server) GetAppReadStates(
	ctx context.Context,
	in *npool.GetAppReadStatesRequest,
) (
	*npool.GetAppReadStatesResponse,
	error,
) {
	var err error

	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "GetAppReadStates")
	defer span.End()

	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	_, err = uuid.Parse(in.GetAppID())
	if err != nil {
		logger.Sugar().Errorw("GetAppReadStates", "AppID", in.GetAppID(), "error", err)
		return &npool.GetAppReadStatesResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	infos, total, err := readstate1.GetReadStates(ctx, in.GetAppID(), nil, in.GetOffset(), in.GetLimit())
	if err != nil {
		logger.Sugar().Errorw("GetAppReadStates", "error", err)
		return &npool.GetAppReadStatesResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.GetAppReadStatesResponse{
		Infos: infos,
		Total: total,
	}, nil
}

func (s *Server) GetNAppReadStates(
	ctx context.Context,
	in *npool.GetNAppReadStatesRequest,
) (
	*npool.GetNAppReadStatesResponse,
	error,
) {
	var err error

	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "GetNAppReadStates")
	defer span.End()

	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	_, err = uuid.Parse(in.GetTargetAppID())
	if err != nil {
		logger.Sugar().Errorw("GetNAppReadStates", "TargetAppID", in.GetTargetAppID(), "error", err)
		return &npool.GetNAppReadStatesResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	infos, total, err := readstate1.GetReadStates(ctx, in.GetTargetAppID(), nil, in.GetOffset(), in.GetLimit())
	if err != nil {
		logger.Sugar().Errorw("GetNAppReadStates", "error", err)
		return &npool.GetNAppReadStatesResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.GetNAppReadStatesResponse{
		Infos: infos,
		Total: total,
	}, nil
}
