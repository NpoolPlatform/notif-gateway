package channel

import (
	"context"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	npool "github.com/NpoolPlatform/message/npool/notif/gw/v1/notif/channel"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	channel1 "github.com/NpoolPlatform/notif-gateway/pkg/notif/channel"
)

func (s *Server) CreateChannels(ctx context.Context, in *npool.CreateChannelsRequest) (*npool.CreateChannelsResponse, error) {
	handler, err := channel1.NewHandler(
		ctx,
		channel1.WithAppID(&in.AppID),
		channel1.WithChannel(&in.Channel),
		channel1.WithEventTypes(in.EventTypes),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"CreateChannels",
			"In", in,
			"Error", err,
		)
		return &npool.CreateChannelsResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	info, err := handler.CreateChannels(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"CreateChannels",
			"In", in,
			"Error", err,
		)
		return &npool.CreateChannelsResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.CreateChannelsResponse{
		Infos: info,
	}, nil
}
