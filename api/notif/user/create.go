package user

import (
	"context"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	npool "github.com/NpoolPlatform/message/npool/notif/gw/v1/notif/user"
	notifuser1 "github.com/NpoolPlatform/notif-gateway/pkg/notif/user"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// TODO: Need Refactor
func (s *Server) CreateNotifUser(ctx context.Context, in *npool.CreateNotifUserRequest) (*npool.CreateNotifUserResponse, error) {
	handler, err := notifuser1.NewHandler(
		ctx,
		notifuser1.WithAppID(&in.AppID),
		notifuser1.WithUserID(&in.AppID, &in.UserID),
		notifuser1.WithNotifID(&in.AppID, &in.NotifID),
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

func (s *Server) CreateNotifUsers(ctx context.Context, in *npool.CreateNotifUsersRequest) (*npool.CreateNotifUsersResponse, error) {
	notifUsers := []*npool.NotifUser{}
	for _, userID := range in.GetUserIDs() {
		resp, err := s.CreateNotifUser(ctx, &npool.CreateNotifUserRequest{
			AppID:   in.AppID,
			UserID:  userID,
			NotifID: in.NotifID,
		})
		if err != nil {
			logger.Sugar().Errorw(
				"CreateNotifUsers",
				"In", in,
				"Error", err,
			)
			return &npool.CreateNotifUsersResponse{}, status.Error(codes.InvalidArgument, err.Error())
		}
		notifUsers = append(notifUsers, resp.Info)
	}

	return &npool.CreateNotifUsersResponse{
		Infos: notifUsers,
	}, nil
}
