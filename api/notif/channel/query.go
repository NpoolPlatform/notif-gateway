package channel

import (
	"context"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	npool "github.com/NpoolPlatform/message/npool/notif/gw/v1/notif/channel"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	channel1 "github.com/NpoolPlatform/notif-gateway/pkg/notif/channel"
)

//nolint:dupl
func (s *Server) GetAppChannels(ctx context.Context, in *npool.GetAppChannelsRequest) (*npool.GetAppChannelsResponse, error) {
	handler, err := channel1.NewHandler(
		ctx,
		channel1.WithAppID(&in.AppID, true),
		channel1.WithOffset(in.Offset),
		channel1.WithLimit(in.Limit),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"GetAppChannels",
			"In", in,
			"Error", err,
		)
		return &npool.GetAppChannelsResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	infos, total, err := handler.GetChannels(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"GetAppChannels",
			"In", in,
			"Error", err,
		)
		return &npool.GetAppChannelsResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.GetAppChannelsResponse{
		Infos: infos,
		Total: total,
	}, nil
}

//nolint:dupl
func (s *Server) GetNAppChannels(ctx context.Context, in *npool.GetNAppChannelsRequest) (*npool.GetNAppChannelsResponse, error) {
	handler, err := channel1.NewHandler(
		ctx,
		channel1.WithAppID(&in.TargetAppID, true),
		channel1.WithOffset(in.Offset),
		channel1.WithLimit(in.Limit),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"GetNAppChannels",
			"In", in,
			"Error", err,
		)
		return &npool.GetNAppChannelsResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	infos, total, err := handler.GetChannels(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"GetNAppChannels",
			"In", in,
			"Error", err,
		)
		return &npool.GetNAppChannelsResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.GetNAppChannelsResponse{
		Infos: infos,
		Total: total,
	}, nil
}
