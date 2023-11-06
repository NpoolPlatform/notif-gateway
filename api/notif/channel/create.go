//nolint:dupl
package channel

import (
	"context"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	npool "github.com/NpoolPlatform/message/npool/notif/gw/v1/notif/channel"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	channel1 "github.com/NpoolPlatform/notif-gateway/pkg/notif/channel"
)

func (s *Server) CreateChannel(ctx context.Context, in *npool.CreateChannelRequest) (*npool.CreateChannelResponse, error) {
	handler, err := channel1.NewHandler(
		ctx,
		channel1.WithAppID(&in.AppID, true),
		channel1.WithChannel(&in.Channel, true),
		channel1.WithEventType(&in.EventType, true),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"CreateChannel",
			"In", in,
			"Error", err,
		)
		return &npool.CreateChannelResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	info, err := handler.CreateChannel(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"CreateChannel",
			"In", in,
			"Error", err,
		)
		return &npool.CreateChannelResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.CreateChannelResponse{
		Info: info,
	}, nil
}
