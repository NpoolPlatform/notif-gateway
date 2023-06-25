package announcement

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	npool "github.com/NpoolPlatform/message/npool/notif/gw/v1/announcement"

	amt1 "github.com/NpoolPlatform/notif-gateway/pkg/announcement"
)

func (s *Server) DeleteAnnouncement(
	ctx context.Context,
	in *npool.DeleteAnnouncementRequest,
) (
	*npool.DeleteAnnouncementResponse,
	error,
) {
	handler, err := amt1.NewHandler(
		ctx,
		amt1.WithID(&in.ID),
		amt1.WithAppID(&in.AppID),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"DeleteAnnouncement",
			"In", in,
			"Error", err,
		)
		return &npool.DeleteAnnouncementResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	info, err := handler.DeleteAnnouncement(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"DeleteAnnouncement",
			"In", in,
			"Error", err,
		)
		return &npool.DeleteAnnouncementResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.DeleteAnnouncementResponse{
		Info: info,
	}, nil
}
