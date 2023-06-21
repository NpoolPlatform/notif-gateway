package readstate

import (
	"context"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	npool "github.com/NpoolPlatform/message/npool/notif/gw/v1/notif/readstate"
	notifreadstate1 "github.com/NpoolPlatform/notif-gateway/pkg/notif/readstate"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) CreateReadState(
	ctx context.Context,
	in *npool.CreateReadStateRequest,
) (
	*npool.CreateReadStateResponse,
	error,
) {
	handler, err := notifreadstate1.NewHandler(
		ctx,
		notifreadstate1.WithAppID(&in.AppID),
		notifreadstate1.WithUserID(&in.AppID, &in.UserID),
		notifreadstate1.WithNotifID(&in.AppID, &in.NotifID),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"CreateReadState",
			"In", in,
			"Error", err,
		)
		return &npool.CreateReadStateResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	info, err := handler.CreateReadState(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"CreateReadState",
			"In", in,
			"Error", err,
		)
		return &npool.CreateReadStateResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.CreateReadStateResponse{
		Info: info,
	}, nil
}
