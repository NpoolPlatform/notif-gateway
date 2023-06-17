package user

import (
	"context"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	npool "github.com/NpoolPlatform/message/npool/notif/gw/v1/notif/user"

	notifuser1 "github.com/NpoolPlatform/notif-gateway/pkg/notif/user"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) GetNotifUsers(ctx context.Context, in *npool.GetNotifUsersRequest) (*npool.GetNotifUsersResponse, error) {
	handler, err := notifuser1.NewHandler(
		ctx,
		notifuser1.WithAppID(&in.AppID),
		notifuser1.WithNotifID(&in.NotifID),
		notifuser1.WithOffset(in.Offset),
		notifuser1.WithLimit(in.Limit),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"GetNotifUsers",
			"In", in,
			"Error", err,
		)
		return &npool.GetNotifUsersResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	infos, total, err := handler.GetNotifUsers(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"GetNotifUsers",
			"In", in,
			"Error", err,
		)
		return &npool.GetNotifUsersResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.GetNotifUsersResponse{
		Infos: infos,
		Total: total,
	}, nil
}

func (s *Server) GetAppNotifUsers(ctx context.Context, in *npool.GetAppNotifUsersRequest) (*npool.GetAppNotifUsersResponse, error) {
	handler, err := notifuser1.NewHandler(
		ctx,
		notifuser1.WithAppID(&in.AppID),
		notifuser1.WithOffset(in.Offset),
		notifuser1.WithLimit(in.Limit),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"GetNotifUsers",
			"In", in,
			"Error", err,
		)
		return &npool.GetAppNotifUsersResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	infos, total, err := handler.GetNotifUsers(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"GetNotifUsers",
			"In", in,
			"Error", err,
		)
		return &npool.GetAppNotifUsersResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.GetAppNotifUsersResponse{
		Infos: infos,
		Total: total,
	}, nil
}

func (s *Server) GetNAppNotifUsers(ctx context.Context, in *npool.GetNAppNotifUsersRequest) (*npool.GetNAppNotifUsersResponse, error) {
	resp, err := s.GetAppNotifUsers(ctx, &npool.GetAppNotifUsersRequest{
		AppID:  in.TargetAppID,
		Offset: in.Offset,
		Limit:  in.Limit,
	})

	if err != nil {
		logger.Sugar().Errorw(
			"GetNAppNotifUsers",
			"In", in,
			"Error", err,
		)
		return &npool.GetNAppNotifUsersResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}
	return &npool.GetNAppNotifUsersResponse{
		Infos: resp.Infos,
		Total: resp.Total,
	}, nil
}
