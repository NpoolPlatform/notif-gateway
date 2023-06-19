package user

import (
	"context"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	npool "github.com/NpoolPlatform/message/npool/notif/gw/v1/announcement/user"
	handler1 "github.com/NpoolPlatform/notif-gateway/pkg/announcement/handler"
	amtuser1 "github.com/NpoolPlatform/notif-gateway/pkg/announcement/user"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) CreateAnnouncementUser(
	ctx context.Context,
	in *npool.CreateAnnouncementUserRequest,
) (
	*npool.CreateAnnouncementUserResponse,
	error,
) {
	handler, err := amtuser1.NewHandler(
		ctx,
		handler1.WithAppID(&in.AppID),
		handler1.WithUserID(&in.AppID, &in.UserID),
		handler1.WithAnnouncementID(&in.AppID, &in.AnnouncementID),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"CreateAnnouncementUser",
			"In", in,
			"Error", err,
		)
		return &npool.CreateAnnouncementUserResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	info, err := handler.CreateAnnouncementUser(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"CreateAnnouncementUser",
			"In", in,
			"Error", err,
		)
		return &npool.CreateAnnouncementUserResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.CreateAnnouncementUserResponse{
		Info: info,
	}, nil
}

func (s *Server) CreateAnnouncementUsers(
	ctx context.Context,
	in *npool.CreateAnnouncementUsersRequest,
) (
	*npool.CreateAnnouncementUsersResponse,
	error,
) {
	announcementUsers := []*npool.AnnouncementUser{}
	for _, userID := range in.GetUserIDs() {
		resp, err := s.CreateAnnouncementUser(ctx, &npool.CreateAnnouncementUserRequest{
			AppID:          in.AppID,
			UserID:         userID,
			AnnouncementID: in.AnnouncementID,
		})
		if err != nil {
			logger.Sugar().Errorw(
				"CreateAnnouncementUsers",
				"In", in,
				"Error", err,
			)
			return &npool.CreateAnnouncementUsersResponse{}, status.Error(codes.InvalidArgument, err.Error())
		}
		announcementUsers = append(announcementUsers, resp.Info)
	}

	return &npool.CreateAnnouncementUsersResponse{
		Infos: announcementUsers,
	}, nil
}
