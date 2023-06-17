package readstate

import (
	"context"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	npool "github.com/NpoolPlatform/message/npool/notif/gw/v1/notif/readstate"

	notifreadstate1 "github.com/NpoolPlatform/notif-gateway/pkg/notif/readstate"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) GetReadState(ctx context.Context, in *npool.GetReadStateRequest) (*npool.GetReadStateResponse, error) {
	handler, err := notifreadstate1.NewHandler(
		ctx,
		notifreadstate1.WithAppID(&in.AppID),
		notifreadstate1.WithUserID(&in.AppID, &in.UserID),
		notifreadstate1.WithNotifID(&in.NotifID),
		notifreadstate1.WithOffset(0),
		notifreadstate1.WithLimit(1),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"GetReadStates",
			"In", in,
			"Error", err,
		)
		return &npool.GetReadStateResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	infos, _, err := handler.GetReadStates(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"GetReadStates",
			"In", in,
			"Error", err,
		)
		return &npool.GetReadStateResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.GetReadStateResponse{
		Info: infos[0],
	}, nil
}

func (s *Server) GetReadStates(ctx context.Context, in *npool.GetReadStatesRequest) (*npool.GetReadStatesResponse, error) {
	handler, err := notifreadstate1.NewHandler(
		ctx,
		notifreadstate1.WithAppID(&in.AppID),
		notifreadstate1.WithUserID(&in.AppID, &in.UserID),
		notifreadstate1.WithOffset(in.Offset),
		notifreadstate1.WithLimit(in.Limit),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"GetReadStates",
			"In", in,
			"Error", err,
		)
		return &npool.GetReadStatesResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	infos, total, err := handler.GetReadStates(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"GetReadStates",
			"In", in,
			"Error", err,
		)
		return &npool.GetReadStatesResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.GetReadStatesResponse{
		Infos: infos,
		Total: total,
	}, nil
}

func (s *Server) GetAppUserReadStates(
	ctx context.Context,
	in *npool.GetAppUserReadStatesRequest,
) (
	*npool.GetAppUserReadStatesResponse,
	error,
) {
	resp, err := s.GetReadStates(ctx, &npool.GetReadStatesRequest{
		AppID:  in.TargetAppID,
		UserID: in.TargetUserID,
		Offset: in.Offset,
		Limit:  in.Limit,
	})

	if err != nil {
		logger.Sugar().Errorw(
			"GetAppUserReadStates",
			"In", in,
			"Error", err,
		)
		return &npool.GetAppUserReadStatesResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.GetAppUserReadStatesResponse{
		Infos: resp.Infos,
		Total: resp.Total,
	}, nil
}

func (s *Server) GetAppReadStates(ctx context.Context, in *npool.GetAppReadStatesRequest) (*npool.GetAppReadStatesResponse, error) {
	handler, err := notifreadstate1.NewHandler(
		ctx,
		notifreadstate1.WithAppID(&in.AppID),
		notifreadstate1.WithOffset(in.Offset),
		notifreadstate1.WithLimit(in.Limit),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"GetReadStates",
			"In", in,
			"Error", err,
		)
		return &npool.GetAppReadStatesResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	infos, total, err := handler.GetReadStates(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"GetReadStates",
			"In", in,
			"Error", err,
		)
		return &npool.GetAppReadStatesResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.GetAppReadStatesResponse{
		Infos: infos,
		Total: total,
	}, nil
}

func (s *Server) GetNAppReadStates(ctx context.Context, in *npool.GetNAppReadStatesRequest) (*npool.GetNAppReadStatesResponse, error) {
	resp, err := s.GetAppReadStates(ctx, &npool.GetAppReadStatesRequest{
		AppID:  in.TargetAppID,
		Offset: in.Offset,
		Limit:  in.Limit,
	})

	if err != nil {
		logger.Sugar().Errorw(
			"GetNAppReadStates",
			"In", in,
			"Error", err,
		)
		return &npool.GetNAppReadStatesResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}
	return &npool.GetNAppReadStatesResponse{
		Infos: resp.Infos,
		Total: resp.Total,
	}, nil
}
