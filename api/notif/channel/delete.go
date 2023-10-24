package channel

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	npool "github.com/NpoolPlatform/message/npool/notif/gw/v1/notif/channel"

	channel1 "github.com/NpoolPlatform/notif-gateway/pkg/notif/channel"
)

func (s *Server) DeleteChannel(ctx context.Context, in *npool.DeleteChannelRequest) (*npool.DeleteChannelResponse, error) {
	handler, err := channel1.NewHandler(
		ctx,
		channel1.WithID(&in.ID, true),
		channel1.WithAppID(&in.AppID, true),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"DeleteChannel",
			"In", in,
			"Error", err,
		)
		return &npool.DeleteChannelResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	info, err := handler.DeleteChannel(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"DeleteChannel",
			"In", in,
			"Error", err,
		)
		return &npool.DeleteChannelResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.DeleteChannelResponse{
		Info: info,
	}, nil
}
