package user

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	npool "github.com/NpoolPlatform/message/npool/notif/gw/v1/announcement/user"

	handler1 "github.com/NpoolPlatform/notif-gateway/pkg/announcement/handler"
	amtuser1 "github.com/NpoolPlatform/notif-gateway/pkg/announcement/user"
)

//nolint:dupl
func (s *Server) DeleteAnnouncementUser(
	ctx context.Context,
	in *npool.DeleteAnnouncementUserRequest,
) (
	*npool.DeleteAnnouncementUserResponse,
	error,
) {
	handler, err := amtuser1.NewHandler(
		ctx,
		handler1.WithID(&in.ID, true),
		handler1.WithEntID(&in.EntID, true),
		handler1.WithAppID(&in.AppID, true),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"DeleteAnnouncementUser",
			"In", in,
			"Error", err,
		)
		return &npool.DeleteAnnouncementUserResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	info, err := handler.DeleteAnnouncementUser(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"DeleteAnnouncementUser",
			"In", in,
			"Error", err,
		)
		return &npool.DeleteAnnouncementUserResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.DeleteAnnouncementUserResponse{
		Info: info,
	}, nil
}

//nolint:dupl
func (s *Server) DeleteAppAnnouncementUser(
	ctx context.Context,
	in *npool.DeleteAppAnnouncementUserRequest,
) (
	*npool.DeleteAppAnnouncementUserResponse,
	error,
) {
	handler, err := amtuser1.NewHandler(
		ctx,
		handler1.WithID(&in.ID, true),
		handler1.WithEntID(&in.EntID, true),
		handler1.WithAppID(&in.TargetAppID, true),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"DeleteAppAnnouncementUser",
			"In", in,
			"Error", err,
		)
		return &npool.DeleteAppAnnouncementUserResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	info, err := handler.DeleteAnnouncementUser(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"DeleteAppAnnouncementUser",
			"In", in,
			"Error", err,
		)
		return &npool.DeleteAppAnnouncementUserResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.DeleteAppAnnouncementUserResponse{
		Info: info,
	}, nil
}
