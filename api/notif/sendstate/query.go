package sendstate

import (
	"context"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	npool "github.com/NpoolPlatform/message/npool/notif/gw/v1/notif/sendstate"

	notifsendstate1 "github.com/NpoolPlatform/notif-gateway/pkg/notif/sendstate"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) GetSendStates(ctx context.Context, in *npool.GetSendStatesRequest) (*npool.GetSendStatesResponse, error) {
	handler, err := notifsendstate1.NewHandler(
		ctx,
		notifsendstate1.WithAppID(&in.AppID),
		notifsendstate1.WithUserID(&in.AppID, &in.UserID),
		notifsendstate1.WithChannel(in.Channel),
		notifsendstate1.WithOffset(in.Offset),
		notifsendstate1.WithLimit(in.Limit),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"GetSendStates",
			"In", in,
			"Error", err,
		)
		return &npool.GetSendStatesResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	infos, total, err := handler.GetSendStates(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"GetSendStates",
			"In", in,
			"Error", err,
		)
		return &npool.GetSendStatesResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.GetSendStatesResponse{
		Infos: infos,
		Total: total,
	}, nil
}

func (s *Server) GetAppUserSendStates(ctx context.Context, in *npool.GetAppUserSendStatesRequest) (*npool.GetAppUserSendStatesResponse, error) {
	resp, err := s.GetSendStates(ctx, &npool.GetSendStatesRequest{
		AppID:   in.TargetAppID,
		UserID:  in.TargetUserID,
		Channel: in.Channel,
		Offset:  in.Offset,
		Limit:   in.Limit,
	})

	if err != nil {
		logger.Sugar().Errorw(
			"GetAppUserSendStates",
			"In", in,
			"Error", err,
		)
		return &npool.GetAppUserSendStatesResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.GetAppUserSendStatesResponse{
		Infos: resp.Infos,
		Total: resp.Total,
	}, nil
}

func (s *Server) GetAppSendStates(ctx context.Context, in *npool.GetAppSendStatesRequest) (*npool.GetAppSendStatesResponse, error) {
	handler, err := notifsendstate1.NewHandler(
		ctx,
		notifsendstate1.WithAppID(&in.AppID),
		notifsendstate1.WithOffset(in.Offset),
		notifsendstate1.WithLimit(in.Limit),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"GetSendStates",
			"In", in,
			"Error", err,
		)
		return &npool.GetAppSendStatesResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	infos, total, err := handler.GetSendStates(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"GetSendStates",
			"In", in,
			"Error", err,
		)
		return &npool.GetAppSendStatesResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.GetAppSendStatesResponse{
		Infos: infos,
		Total: total,
	}, nil
}

func (s *Server) GetNAppSendStates(ctx context.Context, in *npool.GetNAppSendStatesRequest) (*npool.GetNAppSendStatesResponse, error) {
	resp, err := s.GetAppSendStates(ctx, &npool.GetAppSendStatesRequest{
		AppID:  in.TargetAppID,
		Offset: in.Offset,
		Limit:  in.Limit,
	})

	if err != nil {
		logger.Sugar().Errorw(
			"GetNAppSendStates",
			"In", in,
			"Error", err,
		)
		return &npool.GetNAppSendStatesResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}
	return &npool.GetNAppSendStatesResponse{
		Infos: resp.Infos,
		Total: resp.Total,
	}, nil
}
