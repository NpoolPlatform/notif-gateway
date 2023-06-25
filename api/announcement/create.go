package announcement

import (
	"context"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	npool "github.com/NpoolPlatform/message/npool/notif/gw/v1/announcement"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	amt1 "github.com/NpoolPlatform/notif-gateway/pkg/announcement"
)

func (s *Server) CreateAnnouncement(
	ctx context.Context,
	in *npool.CreateAnnouncementRequest,
) (
	*npool.CreateAnnouncementResponse,
	error,
) {
	handler, err := amt1.NewHandler(
		ctx,
		amt1.WithTitle(&in.Title),
		amt1.WithContent(&in.Content),
		amt1.WithAppID(&in.AppID),
		amt1.WithLangID(&in.AppID, &in.TargetLangID),
		amt1.WithChannel(&in.Channel),
		amt1.WithAnnouncementType(&in.AnnouncementType),
		amt1.WithStartAt(&in.StartAt, &in.StartAt),
		amt1.WithEndAt(&in.StartAt, &in.EndAt),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"CreateAnnouncement",
			"In", in,
			"Error", err,
		)
		return &npool.CreateAnnouncementResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	info, err := handler.CreateAnnouncement(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"CreateAnnouncement",
			"In", in,
			"Error", err,
		)
		return &npool.CreateAnnouncementResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.CreateAnnouncementResponse{
		Info: info,
	}, nil
}
