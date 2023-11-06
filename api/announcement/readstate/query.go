package readstate

import (
	"context"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	npool "github.com/NpoolPlatform/message/npool/notif/gw/v1/announcement/readstate"

	handler1 "github.com/NpoolPlatform/notif-gateway/pkg/announcement/handler"
	amtread1 "github.com/NpoolPlatform/notif-gateway/pkg/announcement/readstate"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) GetReadState(
	ctx context.Context,
	in *npool.GetReadStateRequest,
) (
	*npool.GetReadStateResponse,
	error,
) {
	handler, err := amtread1.NewHandler(
		ctx,
		handler1.WithAppID(&in.AppID, true),
		handler1.WithUserID(&in.UserID, true),
		handler1.WithAnnouncementID(&in.AppID, &in.AnnouncementID, true),
		handler1.WithOffset(0),
		handler1.WithLimit(1),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"GetReadState",
			"In", in,
			"Error", err,
		)
		return &npool.GetReadStateResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	info, err := handler.GetReadState(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"GetReadState",
			"In", in,
			"Error", err,
		)
		return &npool.GetReadStateResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.GetReadStateResponse{
		Info: info,
	}, nil
}

//nolint:dupl
func (s *Server) GetReadStates(ctx context.Context, in *npool.GetReadStatesRequest) (*npool.GetReadStatesResponse, error) {
	handler, err := amtread1.NewHandler(
		ctx,
		handler1.WithAppID(&in.AppID, true),
		handler1.WithUserID(&in.UserID, true),
		handler1.WithOffset(in.Offset),
		handler1.WithLimit(in.Limit),
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

//nolint:dupl
func (s *Server) GetAppUserReadStates(ctx context.Context, in *npool.GetAppUserReadStatesRequest) (*npool.GetAppUserReadStatesResponse, error) {
	handler, err := amtread1.NewHandler(
		ctx,
		handler1.WithAppID(&in.TargetAppID, true),
		handler1.WithUserID(&in.TargetUserID, true),
		handler1.WithOffset(in.Offset),
		handler1.WithLimit(in.Limit),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"GetAppUserReadStates",
			"In", in,
			"Error", err,
		)
		return &npool.GetAppUserReadStatesResponse{}, status.Error(codes.Internal, err.Error())
	}

	infos, total, err := handler.GetReadStates(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"GetAppUserReadStates",
			"In", in,
			"Error", err,
		)
		return &npool.GetAppUserReadStatesResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.GetAppUserReadStatesResponse{
		Infos: infos,
		Total: total,
	}, nil
}

//nolint:dupl
func (s *Server) GetAppReadStates(ctx context.Context, in *npool.GetAppReadStatesRequest) (*npool.GetAppReadStatesResponse, error) {
	handler, err := amtread1.NewHandler(
		ctx,
		handler1.WithAppID(&in.AppID, true),
		handler1.WithOffset(in.Offset),
		handler1.WithLimit(in.Limit),
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

//nolint:dupl
func (s *Server) GetNAppReadStates(ctx context.Context, in *npool.GetNAppReadStatesRequest) (*npool.GetNAppReadStatesResponse, error) {
	handler, err := amtread1.NewHandler(
		ctx,
		handler1.WithAppID(&in.TargetAppID, true),
		handler1.WithOffset(in.Offset),
		handler1.WithLimit(in.Limit),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"GetNAppReadStates",
			"In", in,
			"Error", err,
		)
		return &npool.GetNAppReadStatesResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	infos, total, err := handler.GetReadStates(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"GetNAppReadStates",
			"In", in,
			"Error", err,
		)
		return &npool.GetNAppReadStatesResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.GetNAppReadStatesResponse{
		Infos: infos,
		Total: total,
	}, nil
}
