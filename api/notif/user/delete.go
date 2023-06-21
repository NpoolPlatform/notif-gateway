package user

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	npool "github.com/NpoolPlatform/message/npool/notif/gw/v1/notif/user"

	notifuser1 "github.com/NpoolPlatform/notif-gateway/pkg/notif/user"
)

func (s *Server) DeleteNotifUser(ctx context.Context, in *npool.DeleteNotifUserRequest) (*npool.DeleteNotifUserResponse, error) {
	handler, err := notifuser1.NewHandler(
		ctx,
		notifuser1.WithID(&in.ID),
		notifuser1.WithAppID(&in.AppID),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"DeleteNotifUser",
			"In", in,
			"Error", err,
		)
		return &npool.DeleteNotifUserResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	info, err := handler.DeleteNotifUser(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"DeleteNotifUser",
			"In", in,
			"Error", err,
		)
		return &npool.DeleteNotifUserResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.DeleteNotifUserResponse{
		Info: info,
	}, nil
}
