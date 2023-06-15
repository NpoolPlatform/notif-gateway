// //nolint:dupl
// package user

// import (
// 	"context"

// 	appmwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/app"
// 	usermwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/user"

// 	mgrcli "github.com/NpoolPlatform/notif-manager/pkg/client/announcement/user"

// 	constant "github.com/NpoolPlatform/notif-gateway/pkg/message/const"

// 	"go.opentelemetry.io/otel"
// 	scodes "go.opentelemetry.io/otel/codes"
// 	"google.golang.org/grpc/codes"
// 	"google.golang.org/grpc/status"

// 	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
// 	usermwpb "github.com/NpoolPlatform/message/npool/appuser/mw/v1/user"
// 	npool "github.com/NpoolPlatform/message/npool/notif/gw/v1/announcement/user/user"

// 	"github.com/google/uuid"

// 	cruder "github.com/NpoolPlatform/libent-cruder/pkg/cruder"
// 	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"
// 	user1 "github.com/NpoolPlatform/notif-gateway/pkg/announcement/user"

// 	announcementmgrcli "github.com/NpoolPlatform/notif-manager/pkg/client/announcement"
// )

// const Limit = 1000

// //nolint
// func (s *Server) CreateAnnouncementUsers(
// 	ctx context.Context,
// 	in *npool.CreateAnnouncementUsersRequest,
// ) (
// 	*npool.CreateAnnouncementUsersResponse,
// 	error,
// ) {
// 	var err error

// 	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "CreateAnnouncementUser")
// 	defer span.End()

// 	defer func() {
// 		if err != nil {
// 			span.SetStatus(scodes.Error, err.Error())
// 			span.RecordError(err)
// 		}
// 	}()

// 	_, err = uuid.Parse(in.GetAppID())
// 	if err != nil {
// 		logger.Sugar().Errorw("CreateAnnouncementUser", "AppID", in.GetAppID(), "error", err)
// 		return &npool.CreateAnnouncementUsersResponse{}, status.Error(codes.InvalidArgument, err.Error())
// 	}

// 	if len(in.GetUserIDs()) == 0 {
// 		logger.Sugar().Errorw("CreateAnnouncementUser", "UserIDs", in.GetUserIDs(), "error", err)
// 		return &npool.CreateAnnouncementUsersResponse{}, status.Error(codes.InvalidArgument, "UserIDs is empty")
// 	}

// 	if len(in.GetUserIDs()) > Limit {
// 		logger.Sugar().Errorw("CreateAnnouncementUser", "UserIDs", in.GetUserIDs(), "error", err)
// 		return &npool.CreateAnnouncementUsersResponse{}, status.Error(codes.InvalidArgument, "UserIDs is too many")
// 	}

// 	idMap := map[string]struct{}{}
// 	for _, id := range in.GetUserIDs() {
// 		if _, ok := idMap[id]; ok {
// 			logger.Sugar().Errorw("CreateAnnouncementUser", "UserID", id, "error", "userID repeat")
// 			return &npool.CreateAnnouncementUsersResponse{}, status.Error(codes.InvalidArgument, "userID repeat")
// 		}
// 		idMap[id] = struct{}{}

// 		_, err = uuid.Parse(id)
// 		if err != nil {
// 			logger.Sugar().Errorw("CreateAnnouncementUser", "UserID", id, "error", err)
// 			return &npool.CreateAnnouncementUsersResponse{}, status.Error(codes.InvalidArgument, err.Error())
// 		}

// 	}

// 	_, err = uuid.Parse(in.GetAnnouncementID())
// 	if err != nil {
// 		logger.Sugar().Errorw("CreateAnnouncementUser", "AnnouncementID", in.GetAnnouncementID(), "error", err)
// 		return &npool.CreateAnnouncementUsersResponse{}, status.Error(codes.InvalidArgument, err.Error())
// 	}

// 	exist, err := announcementmgrcli.ExistAnnouncement(ctx, in.GetAnnouncementID())
// 	if err != nil {
// 		logger.Sugar().Errorw("CreateAnnouncementUser", "AnnouncementID", in.GetAnnouncementID(), "error", err)
// 		return &npool.CreateAnnouncementUsersResponse{}, status.Error(codes.InvalidArgument, err.Error())
// 	}
// 	if !exist {
// 		logger.Sugar().Errorw("CreateAnnouncementUser", "AnnouncementID", in.GetAnnouncementID(), "error", "Announcement not exist")
// 		return &npool.CreateAnnouncementUsersResponse{}, status.Error(codes.InvalidArgument, "Announcement not exist")
// 	}

// 	userInfos, _, err := usermwcli.GetUsers(ctx, &usermwpb.Conds{
// 		IDs: &basetypes.StringSliceVal{Op: cruder.IN, Value: in.GetUserIDs()},
// 	}, 0, int32(len(in.GetUserIDs())))
// 	if err != nil {
// 		logger.Sugar().Errorw("CreateAnnouncementUser", "AnnouncementID", in.GetAnnouncementID(), "error", err)
// 		return &npool.CreateAnnouncementUsersResponse{}, status.Error(codes.InvalidArgument, err.Error())
// 	}
// 	if len(userInfos) != len(in.GetUserIDs()) {
// 		logger.Sugar().Errorw("CreateAnnouncementUser", "AnnouncementID", in.GetAnnouncementID(), "error", "User not exist")
// 		return &npool.CreateAnnouncementUsersResponse{}, status.Error(codes.InvalidArgument, "User not exist")
// 	}

// 	appInfo, err := appmwcli.GetApp(ctx, in.GetAppID())
// 	if err != nil {
// 		return nil, err
// 	}
// 	if err != nil {
// 		logger.Sugar().Errorw("CreateAnnouncementUser", "AnnouncementID", in.GetAnnouncementID(), "error", err)
// 		return &npool.CreateAnnouncementUsersResponse{}, status.Error(codes.InvalidArgument, err.Error())
// 	}
// 	if appInfo == nil {
// 		logger.Sugar().Errorw("CreateAnnouncementUser", "AnnouncementID", in.GetAnnouncementID(), "error", "Announcement not exist")
// 		return &npool.CreateAnnouncementUsersResponse{}, status.Error(codes.InvalidArgument, "App not exist")
// 	}

// 	infos, _, err := user1.CreateAnnouncementUsers(
// 		ctx,
// 		in.GetAppID(),
// 		in.GetUserIDs(),
// 		in.GetAnnouncementID(),
// 	)
// 	if err != nil {
// 		logger.Sugar().Errorw("CreateAnnouncementUser", "error", err)
// 		return &npool.CreateAnnouncementUsersResponse{}, status.Error(codes.Internal, err.Error())
// 	}

// 	return &npool.CreateAnnouncementUsersResponse{
// 		Infos: infos,
// 	}, nil
// }

// func (s *Server) DeleteAnnouncementUser(
// 	ctx context.Context,
// 	in *npool.DeleteAnnouncementUserRequest,
// ) (
// 	*npool.DeleteAnnouncementUserResponse,
// 	error,
// ) {
// 	var err error

// 	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "CreateAnnouncementUser")
// 	defer span.End()

// 	defer func() {
// 		if err != nil {
// 			span.SetStatus(scodes.Error, err.Error())
// 			span.RecordError(err)
// 		}
// 	}()

// 	_, err = uuid.Parse(in.GetID())
// 	if err != nil {
// 		logger.Sugar().Errorw("CreateAnnouncementUser", "ID", in.GetID(), "error", err)
// 		return &npool.DeleteAnnouncementUserResponse{}, status.Error(codes.InvalidArgument, err.Error())
// 	}

// 	_, err = uuid.Parse(in.GetAppID())
// 	if err != nil {
// 		logger.Sugar().Errorw("CreateAnnouncementUser", "AppID", in.GetAppID(), "error", err)
// 		return &npool.DeleteAnnouncementUserResponse{}, status.Error(codes.InvalidArgument, err.Error())
// 	}

// 	info, err := mgrcli.GetUser(ctx, in.GetID())
// 	if err != nil {
// 		logger.Sugar().Errorw("CreateAnnouncementUser", "ID", in.GetID(), "error", err)
// 		return &npool.DeleteAnnouncementUserResponse{}, status.Error(codes.InvalidArgument, err.Error())
// 	}

// 	if info.AppID != in.GetAppID() {
// 		logger.Sugar().Errorw("CreateAnnouncementUser", "AppID", info.AppID, "error", err)
// 		return &npool.DeleteAnnouncementUserResponse{}, status.Error(codes.PermissionDenied, "permission denied")
// 	}

// 	info1, err := user1.DeleteAnnouncementUser(
// 		ctx,
// 		in.GetID(),
// 	)
// 	if err != nil {
// 		logger.Sugar().Errorw("CreateAnnouncementUser", "error", err)
// 		return &npool.DeleteAnnouncementUserResponse{}, status.Error(codes.Internal, err.Error())
// 	}

// 	return &npool.DeleteAnnouncementUserResponse{
// 		Info: info1,
// 	}, nil
// }

// func (s *Server) GetAnnouncementUsers(
// 	ctx context.Context,
// 	in *npool.GetAnnouncementUsersRequest,
// ) (
// 	*npool.GetAnnouncementUsersResponse,
// 	error,
// ) {
// 	var err error

// 	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "GetAnnouncementUsers")
// 	defer span.End()

// 	defer func() {
// 		if err != nil {
// 			span.SetStatus(scodes.Error, err.Error())
// 			span.RecordError(err)
// 		}
// 	}()

// 	_, err = uuid.Parse(in.GetAppID())
// 	if err != nil {
// 		logger.Sugar().Errorw("GetAnnouncementUsers", "AppID", in.GetAppID(), "error", err)
// 		return &npool.GetAnnouncementUsersResponse{}, status.Error(codes.InvalidArgument, err.Error())
// 	}

// 	_, err = uuid.Parse(in.GetAnnouncementID())
// 	if err != nil {
// 		logger.Sugar().Errorw("GetAnnouncementUsers", "AnnouncementID", in.GetAnnouncementID(), "error", err)
// 		return &npool.GetAnnouncementUsersResponse{}, status.Error(codes.InvalidArgument, err.Error())
// 	}

// 	infos, total, err := user1.GetAnnouncementUsers(ctx, in.GetAppID(), &in.AnnouncementID, nil, in.GetOffset(), in.GetLimit())
// 	if err != nil {
// 		logger.Sugar().Errorw("GetAnnouncementUsers", "error", err)
// 		return &npool.GetAnnouncementUsersResponse{}, status.Error(codes.Internal, err.Error())
// 	}

// 	return &npool.GetAnnouncementUsersResponse{
// 		Infos: infos,
// 		Total: total,
// 	}, nil
// }

// func (s *Server) GetAppAnnouncementUsers(
// 	ctx context.Context,
// 	in *npool.GetAppAnnouncementUsersRequest,
// ) (
// 	*npool.GetAppAnnouncementUsersResponse,
// 	error,
// ) {
// 	var err error

// 	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "GetAppAnnouncementUsers")
// 	defer span.End()

// 	defer func() {
// 		if err != nil {
// 			span.SetStatus(scodes.Error, err.Error())
// 			span.RecordError(err)
// 		}
// 	}()

// 	_, err = uuid.Parse(in.GetAppID())
// 	if err != nil {
// 		logger.Sugar().Errorw("GetAppAnnouncementUsers", "AppID", in.GetAppID(), "error", err)
// 		return &npool.GetAppAnnouncementUsersResponse{}, status.Error(codes.InvalidArgument, err.Error())
// 	}

// 	infos, total, err := user1.GetAnnouncementUsers(ctx, in.GetAppID(), nil, nil, in.GetOffset(), in.GetLimit())
// 	if err != nil {
// 		logger.Sugar().Errorw("GetAppAnnouncementUsers", "error", err)
// 		return &npool.GetAppAnnouncementUsersResponse{}, status.Error(codes.Internal, err.Error())
// 	}

// 	return &npool.GetAppAnnouncementUsersResponse{
// 		Infos: infos,
// 		Total: total,
// 	}, nil
// }

// func (s *Server) GetNAppAnnouncementUsers(
// 	ctx context.Context,
// 	in *npool.GetNAppAnnouncementUsersRequest,
// ) (
// 	*npool.GetNAppAnnouncementUsersResponse,
// 	error,
// ) {
// 	var err error

// 	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "GetNAppAnnouncementUsers")
// 	defer span.End()

// 	defer func() {
// 		if err != nil {
// 			span.SetStatus(scodes.Error, err.Error())
// 			span.RecordError(err)
// 		}
// 	}()

// 	_, err = uuid.Parse(in.GetTargetAppID())
// 	if err != nil {
// 		logger.Sugar().Errorw("GetNAppAnnouncementUsers", "TargetAppID", in.GetTargetAppID(), "error", err)
// 		return &npool.GetNAppAnnouncementUsersResponse{}, status.Error(codes.InvalidArgument, err.Error())
// 	}

// 	infos, total, err := user1.GetAnnouncementUsers(ctx, in.GetTargetAppID(), nil, nil, in.GetOffset(), in.GetLimit())
// 	if err != nil {
// 		logger.Sugar().Errorw("GetNAppAnnouncementUsers", "error", err)
// 		return &npool.GetNAppAnnouncementUsersResponse{}, status.Error(codes.Internal, err.Error())
// 	}

// 	return &npool.GetNAppAnnouncementUsersResponse{
// 		Infos: infos,
// 		Total: total,
// 	}, nil
// }
