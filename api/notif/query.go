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

func (s *Server) GetNotif(ctx context.Context, in *npool.GetNotifRequest) (*npool.GetNotifResponse, error) {
	handler, err := notif1.NewHandler(
		ctx,
		notif1.WithID(&in.ID),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"GetNotif",
			"In", in,
			"Error", err,
		)
		return &npool.GetNotifResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	info, err := handler.GetNotif(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"GetNotif",
			"In", in,
			"Error", err,
		)
		return &npool.GetNotifResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.GetNotifResponse{
		Info: info,
	}, nil
}

func (s *Server) GetNotifs(ctx context.Context, in *npool.GetNotifsRequest) (*npool.GetNotifsResponse, error) {
	channel := basetypes.NotifChannel_ChannelFrontend
	hangler, err := notif1.NewHandler(
		ctx,
		notif1.WithAppID(&in.AppID),
		notif1.WithUserID(&in.UserID),
		notif1.WithLangID(&in.LangID),
		notif1.WithChannel(&channel),
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
		notif1.WithAppID(&in.AppID),
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
