package announcement

import (
	"context"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	npool "github.com/NpoolPlatform/message/npool/notif/gw/v1/announcement"
	amt1 "github.com/NpoolPlatform/notif-gateway/pkg/announcement"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) UpdateAnnouncement(
	ctx context.Context,
	in *npool.UpdateAnnouncementRequest,
) (
	*npool.UpdateAnnouncementResponse,
	error,
) {
	handler, err := amt1.NewHandler(
		ctx,
		amt1.WithID(&in.ID),
		amt1.WithAppID(&in.AppID),
		amt1.WithTitle(in.Title),
		amt1.WithContent(in.Content),
		amt1.WithAnnouncementType(in.AnnouncementType),
		amt1.WithEndAt(in.EndAt),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"UpdateAnnouncement",
			"In", in,
			"Error", err,
		)
		return &npool.UpdateAnnouncementResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	info, err := handler.UpdateAnnouncement(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"UpdateAnnouncement",
			"In", in,
			"Error", err,
		)
		return &npool.UpdateAnnouncementResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.UpdateAnnouncementResponse{
		Info: info,
	}, nil
}
