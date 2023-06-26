package announcement

import (
	"context"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	npool "github.com/NpoolPlatform/message/npool/notif/gw/v1/announcement"
	announcement1 "github.com/NpoolPlatform/notif-gateway/pkg/announcement"
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
	handler, err := announcement1.NewHandler(
		ctx,
		announcement1.WithID(&in.ID),
		announcement1.WithAppID(&in.AppID),
		announcement1.WithTitle(in.Title),
		announcement1.WithContent(in.Content),
		announcement1.WithAnnouncementType(in.AnnouncementType),
		announcement1.WithStartAt(in.StartAt, in.EndAt),
		announcement1.WithEndAt(in.StartAt, in.EndAt),
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
