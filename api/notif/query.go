package notif

import (
	"context"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"
	npool "github.com/NpoolPlatform/message/npool/notif/gw/v1/notif"
	notif1 "github.com/NpoolPlatform/notif-gateway/pkg/notif"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) GetNotifs(ctx context.Context, in *npool.GetNotifsRequest) (*npool.GetNotifsResponse, error) {
	channel := basetypes.NotifChannel_ChannelFrontend
	hangler, err := notif1.NewHandler(
		ctx,
		notif1.WithAppID(&in.AppID, true),
		notif1.WithUserID(&in.UserID, true),
		notif1.WithLangID(&in.LangID, true),
		notif1.WithChannel(&channel, true),
		notif1.WithOffset(in.GetOffset()),
		notif1.WithLimit(in.GetLimit()),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"GetNotifs",
			"In", in,
			"Error", err,
		)
		return &npool.GetNotifsResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	infos, total, err := hangler.GetNotifs(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"GetNotifs",
			"In", in,
			"Error", err,
		)
		return &npool.GetNotifsResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.GetNotifsResponse{
		Infos: infos,
		Total: total,
	}, nil
}

func (s *Server) GetAppNotifs(ctx context.Context, in *npool.GetAppNotifsRequest) (*npool.GetAppNotifsResponse, error) {
	hangler, err := notif1.NewHandler(
		ctx,
		notif1.WithAppID(&in.AppID, true),
		notif1.WithOffset(in.GetOffset()),
		notif1.WithLimit(in.GetLimit()),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"GetAppNotifs",
			"In", in,
			"Error", err,
		)
		return &npool.GetAppNotifsResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	infos, total, err := hangler.GetNotifs(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"GetAppNotifs",
			"In", in,
			"Error", err,
		)
		return &npool.GetAppNotifsResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.GetAppNotifsResponse{
		Infos: infos,
		Total: total,
	}, nil
}
