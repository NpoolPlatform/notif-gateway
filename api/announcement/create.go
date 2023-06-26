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
		announcement1.WithTitle(&in.Title),
		announcement1.WithContent(&in.Content),
		announcement1.WithAppID(&in.AppID),
		announcement1.WithLangID(&in.AppID, &in.TargetLangID),
		announcement1.WithChannel(&in.Channel),
		announcement1.WithAnnouncementType(&in.AnnouncementType),
		announcement1.WithStartAt(&in.StartAt),
		announcement1.WithEndAt(&in.EndAt),
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
