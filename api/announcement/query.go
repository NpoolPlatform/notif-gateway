package announcement

import (
	"context"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	npool "github.com/NpoolPlatform/message/npool/notif/gw/v1/announcement"

	announcement1 "github.com/NpoolPlatform/notif-gateway/pkg/announcement"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) GetAnnouncements(ctx context.Context, in *npool.GetAnnouncementsRequest) (*npool.GetAnnouncementsResponse, error) {
	handler, err := announcement1.NewHandler(
		ctx,
		announcement1.WithAppID(&in.AppID, true),
		announcement1.WithUserID(&in.UserID, true),
		announcement1.WithLangID(&in.LangID, true),
		announcement1.WithOffset(in.Offset),
		announcement1.WithLimit(in.Limit),
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

func (s *Server) GetAppAnnouncements(
	ctx context.Context,
	in *npool.GetAppAnnouncementsRequest,
) (
	*npool.GetAppAnnouncementsResponse,
	error,
) {
	handler, err := announcement1.NewHandler(
		ctx,
		announcement1.WithAppID(&in.AppID, true),
		announcement1.WithOffset(in.Offset),
		announcement1.WithLimit(in.Limit),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"GetAppAnnouncements",
			"In", in,
			"Error", err,
		)
		return &npool.GetAppAnnouncementsResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	infos, total, err := handler.GetAppAnnouncements(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"GetAppAnnouncements",
			"In", in,
			"Error", err,
		)
		return &npool.GetAppAnnouncementsResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.GetAppAnnouncementsResponse{
		Infos: infos,
		Total: total,
	}, nil
}

func (s *Server) GetNAppAnnouncements(
	ctx context.Context,
	in *npool.GetNAppAnnouncementsRequest,
) (
	*npool.GetNAppAnnouncementsResponse,
	error,
) {
	resp, err := s.GetAppAnnouncements(ctx, &npool.GetAppAnnouncementsRequest{
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
