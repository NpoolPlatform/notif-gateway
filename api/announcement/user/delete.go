package user

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	npool "github.com/NpoolPlatform/message/npool/notif/gw/v1/announcement/user"

	"github.com/NpoolPlatform/notif-gateway/pkg/announcement/handler"
	amtuser1 "github.com/NpoolPlatform/notif-gateway/pkg/announcement/user"
)

func (s *Server) DeleteAnnouncementUser(ctx context.Context, in *npool.DeleteAnnouncementUserRequest) (*npool.DeleteAnnouncementUserResponse, error) {
	handler, err := amtuser1.NewHandler(
		ctx,
		handler.WithID(&in.ID),
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
