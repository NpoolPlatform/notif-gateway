//nolint:dupl
package user

import (
	"context"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	npool "github.com/NpoolPlatform/message/npool/notif/gw/v1/notif/user"
	notifuser1 "github.com/NpoolPlatform/notif-gateway/pkg/notif/user"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) CreateNotifUser(ctx context.Context, in *npool.CreateNotifUserRequest) (*npool.CreateNotifUserResponse, error) {
	handler, err := notifuser1.NewHandler(
		ctx,
		notifuser1.WithAppID(&in.AppID, true),
		notifuser1.WithUserID(&in.TargetUserID, true),
		notifuser1.WithEventType(&in.EventType, true),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"CreateNotifUser",
			"In", in,
			"Error", err,
		)
		return &npool.CreateNotifUserResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	info, err := handler.CreateNotifUser(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"CreateNotifUser",
			"In", in,
			"Error", err,
		)
		return &npool.CreateNotifUserResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.CreateNotifUserResponse{
		Info: info,
	}, nil
}
