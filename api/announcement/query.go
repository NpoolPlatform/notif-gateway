package announcement

import (
	"context"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	npool "github.com/NpoolPlatform/message/npool/notif/gw/v1/announcement"

	amt1 "github.com/NpoolPlatform/notif-gateway/pkg/announcement"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

//nolint
func (s *Server) GetAnnouncement(ctx context.Context, in *npool.GetAnnouncementRequest) (*npool.GetAnnouncementResponse, error) {
	handler, err := amt1.NewHandler(
		ctx,
		amt1.WithID(&in.ID),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"GetAnnouncement",
			"In", in,
			"Error", err,
		)
		return &npool.GetAnnouncementResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	info, err := handler.GetAnnouncement(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"GetAnnouncement",
			"In", in,
			"Error", err,
		)
		return &npool.GetAnnouncementResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.GetAnnouncementResponse{
		Info: info,
	}, nil
}

func (s *Server) GetAnnouncements(ctx context.Context, in *npool.GetAnnouncementsRequest) (*npool.GetAnnouncementsResponse, error) {
	handler, err := amt1.NewHandler(
		ctx,
		amt1.WithAppID(&in.AppID),
		amt1.WithUserID(&in.UserID),
		amt1.WithLangID(&in.AppID, &in.LangID),
		amt1.WithOffset(in.Offset),
		amt1.WithLimit(in.Limit),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"GetAnnouncements",
			"In", in,
			"Error", err,
		)
		return &npool.GetAnnouncementsResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	infos, total, err := handler.GetAnnouncements(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"GetAnnouncements",
			"In", in,
			"Error", err,
		)
		return &npool.GetAnnouncementsResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.GetAnnouncementsResponse{
		Infos: infos,
		Total: total,
	}, nil
}

//nolint
func (s *Server) GetAppAnnouncements(ctx context.Context, in *npool.GetAppAnnouncementsRequest) (*npool.GetAppAnnouncementsResponse, error) {
	resp, err := s.GetAnnouncements(ctx, &npool.GetAnnouncementsRequest{
		AppID:  in.AppID,
		Offset: in.Offset,
		Limit:  in.Limit,
	})

	if err != nil {
		logger.Sugar().Errorw(
			"GetAppAnnouncements",
			"In", in,
			"Error", err,
		)
		return &npool.GetAppAnnouncementsResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	return &npool.GetAppAnnouncementsResponse{
		Infos: resp.Infos,
		Total: resp.Total,
	}, nil
}

//nolint
func (s *Server) GetNAppAnnouncements(ctx context.Context, in *npool.GetNAppAnnouncementsRequest) (*npool.GetNAppAnnouncementsResponse, error) {
	resp, err := s.GetAnnouncements(ctx, &npool.GetAnnouncementsRequest{
		AppID:  in.TargetAppID,
		Offset: in.Offset,
		Limit:  in.Limit,
	})

	if err != nil {
		logger.Sugar().Errorw(
			"GetNAppAnnouncements",
			"In", in,
			"Error", err,
		)
		return &npool.GetNAppAnnouncementsResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}
	return &npool.GetNAppAnnouncementsResponse{
		Infos: resp.Infos,
		Total: resp.Total,
	}, nil
}
