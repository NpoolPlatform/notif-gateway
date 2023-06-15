package channel

import (
	"context"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	npool "github.com/NpoolPlatform/message/npool/notif/gw/v1/notif/channel"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	channel1 "github.com/NpoolPlatform/notif-gateway/pkg/notif/channel"
)

func (s *Server) GetAppChannels(ctx context.Context, in *npool.GetAppChannelsRequest) (*npool.GetAppChannelsResponse, error) {
	handler, err := channel1.NewHandler(
		ctx,
		channel1.WithAppID(&in.AppID),
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

func (s *Server) GetNAppChannels(ctx context.Context, in *npool.GetNAppChannelsRequest) (*npool.GetNAppChannelsResponse, error) {
	resp, err := s.GetAppChannels(ctx, &npool.GetAppChannelsRequest{
		AppID:  in.TargetAppID,
		Offset: in.Offset,
		Limit:  in.Limit,
	})

	if err != nil {
		logger.Sugar().Errorw(
			"GetAppChannels",
			"In", in,
			"Error", err,
		)
		return &npool.GetNAppChannelsResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	return &npool.GetNAppChannelsResponse{
		Infos: resp.Infos,
		Total: resp.Total,
	}, nil
}
