//nolint:nolintlint,dupl
package notif

import (
	"context"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"

	npool "github.com/NpoolPlatform/message/npool/notif/gw/v1/notif"
	notif1 "github.com/NpoolPlatform/notif-gateway/pkg/notif"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) UpdateNotifs(ctx context.Context, in *npool.UpdateNotifsRequest) (*npool.UpdateNotifsResponse, error) {
	reqs := in.GetInfos()
	handler, err := notif1.NewHandler(
		ctx,
		notif1.WithAppID(&in.AppID),
		notif1.WithUserID(&in.AppID, &in.UserID),
		notif1.WithReqs(reqs),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"UpdateNotifs",
			"In", in,
			"Error", err,
		)
		return &npool.UpdateNotifsResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	infos, err := handler.UpdateNotifs(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"UpdateNotifs",
			"In", in,
			"Error", err,
		)
		return &npool.UpdateNotifsResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.UpdateNotifsResponse{
		Infos: infos,
	}, nil
}
