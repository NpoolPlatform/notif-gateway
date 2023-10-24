package user

import (
	"context"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	npool "github.com/NpoolPlatform/message/npool/notif/gw/v1/announcement/user"

	handler1 "github.com/NpoolPlatform/notif-gateway/pkg/announcement/handler"
	amtuser1 "github.com/NpoolPlatform/notif-gateway/pkg/announcement/user"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) GetAnnouncementUsers(
	ctx context.Context,
	in *npool.GetAnnouncementUsersRequest,
) (
	*npool.GetAnnouncementUsersResponse,
	error,
) {
	handler, err := amtuser1.NewHandler(
		ctx,
		handler1.WithAppID(&in.AppID, true),
		handler1.WithOffset(in.Offset),
		handler1.WithLimit(in.Limit),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"GetAnnouncementUsers",
			"In", in,
			"Error", err,
		)
		return &npool.GetAnnouncementUsersResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	infos, total, err := handler.GetAnnouncementUsers(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"GetAnnouncementUsers",
			"In", in,
			"Error", err,
		)
		return &npool.GetAnnouncementUsersResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.GetAnnouncementUsersResponse{
		Infos: infos,
		Total: total,
	}, nil
}

func (s *Server) GetAppAnnouncementUsers(
	ctx context.Context,
	in *npool.GetAppAnnouncementUsersRequest,
) (
	*npool.GetAppAnnouncementUsersResponse,
	error,
) {
	resp, err := s.GetAnnouncementUsers(ctx, &npool.GetAnnouncementUsersRequest{
		AppID:  in.TargetAppID,
		Offset: in.Offset,
		Limit:  in.Limit,
	})
	if err != nil {
		logger.Sugar().Errorw(
			"GetAppAnnouncementUsers",
			"In", in,
			"Error", err,
		)
		return &npool.GetAppAnnouncementUsersResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	return &npool.GetAppAnnouncementUsersResponse{
		Infos: resp.Infos,
		Total: resp.Total,
	}, nil
}
