package sendstate

import (
	"context"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	npool "github.com/NpoolPlatform/message/npool/notif/gw/v1/announcement/sendstate"

	handler1 "github.com/NpoolPlatform/notif-gateway/pkg/announcement/handler"
	amtsend1 "github.com/NpoolPlatform/notif-gateway/pkg/announcement/sendstate"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

//nolint:dupl
func (s *Server) GetSendStates(
	ctx context.Context,
	in *npool.GetSendStatesRequest,
) (
	*npool.GetSendStatesResponse,
	error,
) {
	handler, err := amtsend1.NewHandler(
		ctx,
		handler1.WithAppID(&in.AppID, true),
		handler1.WithUserID(&in.UserID, true),
		amtsend1.WithChannel(in.Channel),
		handler1.WithOffset(in.Offset),
		handler1.WithLimit(in.Limit),
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

//nolint:dupl
func (s *Server) GetAppUserSendStates(ctx context.Context, in *npool.GetAppUserSendStatesRequest) (*npool.GetAppUserSendStatesResponse, error) {
	handler, err := amtsend1.NewHandler(
		ctx,
		handler1.WithAppID(&in.TargetAppID, true),
		handler1.WithUserID(&in.TargetUserID, true),
		amtsend1.WithChannel(in.Channel),
		handler1.WithOffset(in.Offset),
		handler1.WithLimit(in.Limit),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"GetAppUserSendStates",
			"In", in,
			"Error", err,
		)
		return &npool.GetAppUserSendStatesResponse{}, status.Error(codes.Internal, err.Error())
	}

	infos, total, err := handler.GetSendStates(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"GetAppUserSendStates",
			"In", in,
			"Error", err,
		)
		return &npool.GetAppUserSendStatesResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.GetAppUserSendStatesResponse{
		Infos: infos,
		Total: total,
	}, nil
}

//nolint:dupl
func (s *Server) GetAppSendStates(ctx context.Context, in *npool.GetAppSendStatesRequest) (*npool.GetAppSendStatesResponse, error) {
	handler, err := amtsend1.NewHandler(
		ctx,
		handler1.WithAppID(&in.AppID, true),
		amtsend1.WithChannel(in.Channel),
		handler1.WithOffset(in.Offset),
		handler1.WithLimit(in.Limit),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"GetAppSendStates",
			"In", in,
			"Error", err,
		)
		return &npool.GetAppSendStatesResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	infos, total, err := handler.GetSendStates(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"GetAppSendStates",
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

//nolint:dupl
func (s *Server) GetNAppSendStates(ctx context.Context, in *npool.GetNAppSendStatesRequest) (*npool.GetNAppSendStatesResponse, error) {
	handler, err := amtsend1.NewHandler(
		ctx,
		handler1.WithAppID(&in.TargetAppID, true),
		amtsend1.WithChannel(in.Channel),
		handler1.WithOffset(in.Offset),
		handler1.WithLimit(in.Limit),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"GetNAppSendStates",
			"In", in,
			"Error", err,
		)
		return &npool.GetNAppSendStatesResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	infos, total, err := handler.GetSendStates(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"GetNAppSendStates",
			"In", in,
			"Error", err,
		)
		return &npool.GetNAppSendStatesResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.GetNAppSendStatesResponse{
		Infos: infos,
		Total: total,
	}, nil
}
