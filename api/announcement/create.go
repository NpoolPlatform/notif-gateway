package announcement

import (
	"context"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	npool "github.com/NpoolPlatform/message/npool/notif/gw/v1/announcement"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	announcement1 "github.com/NpoolPlatform/notif-gateway/pkg/announcement"
)

func (s *Server) CreateAnnouncement(
	ctx context.Context,
	in *npool.CreateAnnouncementRequest,
) (
	*npool.CreateAnnouncementResponse,
	error,
) {
	handler, err := announcement1.NewHandler(
		ctx,
		announcement1.WithTitle(&in.Title, true),
		announcement1.WithContent(&in.Content, true),
		announcement1.WithAppID(&in.AppID, true),
		announcement1.WithLangID(&in.TargetLangID, true),
		announcement1.WithChannel(&in.Channel, true),
		announcement1.WithAnnouncementType(&in.AnnouncementType, true),
		announcement1.WithStartAt(&in.StartAt, true),
		announcement1.WithEndAt(&in.EndAt, true),
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
