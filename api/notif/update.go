//nolint:nolintlint,dupl
package notif

import (
	"context"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"

	npool "github.com/NpoolPlatform/message/npool/notif/gw/v1/notif"
	mwpb "github.com/NpoolPlatform/message/npool/notif/mw/v1/notif"
	notif1 "github.com/NpoolPlatform/notif-gateway/pkg/notif"
	mwcli "github.com/NpoolPlatform/notif-middleware/pkg/client/notif"

	cruder "github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) UpdateNotifs(ctx context.Context, in *npool.UpdateNotifsRequest) (*npool.UpdateNotifsResponse, error) {
	handler, err := notif1.NewHandler(
		ctx,
		notif1.WithIDs(in.IDs),
		notif1.WithAppID(&in.AppID),
		notif1.WithUserID(&in.AppID, &in.UserID),
		notif1.WithNotified(&in.Notified),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"UpdateNotifs",
			"In", in,
			"Error", err,
		)
		return &npool.UpdateNotifsResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	rows, _, err := mwcli.GetNotifs(ctx, &mwpb.Conds{
		IDs: &basetypes.StringSliceVal{Op: cruder.IN, Value: in.GetIDs()},
	}, 0, int32(len(in.GetIDs())))
	if err != nil {
		logger.Sugar().Errorw("GetNotif", "error", err)
		return &npool.UpdateNotifsResponse{}, status.Error(codes.Internal, err.Error())
	}

	for _, val := range rows {
		if val.AppID != in.GetAppID() || val.UserID != in.GetUserID() {
			logger.Sugar().Errorw("GetNotif", "error", err)
			return &npool.UpdateNotifsResponse{}, status.Error(codes.PermissionDenied, "permission denied")
		}
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
