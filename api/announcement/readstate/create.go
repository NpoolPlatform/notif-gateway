package readstate

import (
	"context"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	npool "github.com/NpoolPlatform/message/npool/notif/gw/v1/announcement/readstate"
	handler1 "github.com/NpoolPlatform/notif-gateway/pkg/announcement/handler"
	amtreadstate1 "github.com/NpoolPlatform/notif-gateway/pkg/announcement/readstate"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) CreateReadState(ctx context.Context, in *npool.CreateReadStateRequest) (*npool.CreateReadStateResponse, error) { //nolint
	handler, err := amtreadstate1.NewHandler(
		ctx,
		handler1.WithAppID(&in.AppID),
		handler1.WithUserID(&in.AppID, &in.UserID),
		handler1.WithAnnouncementID(&in.AppID, &in.AnnouncementID),
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
