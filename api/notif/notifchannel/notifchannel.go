//nolint:dupl
package notifchannel

import (
	"context"

	appcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/app"
	mgrcli "github.com/NpoolPlatform/notif-manager/pkg/client/notif/notifchannel"

	"github.com/NpoolPlatform/message/npool/notif/mgr/v1/channel"
	"github.com/NpoolPlatform/message/npool/third/mgr/v1/usedfor"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"

	npool "github.com/NpoolPlatform/message/npool/notif/gw/v1/notif/notifchannel"
	constant "github.com/NpoolPlatform/notif-gateway/pkg/message/const"

	commontracer "github.com/NpoolPlatform/notif-gateway/pkg/tracer"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	scodes "go.opentelemetry.io/otel/codes"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	notifchannel1 "github.com/NpoolPlatform/notif-gateway/pkg/notif/notifchannel"
)

//nolint:gocyclo
func (s *Server) CreateNotifChannels(
	ctx context.Context,
	in *npool.CreateNotifChannelsRequest,
) (
	*npool.CreateNotifChannelsResponse,
	error,
) {
	var err error

	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "CreateNotifChannels")
	defer span.End()

	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	if _, err := uuid.Parse(in.GetAppID()); err != nil {
		logger.Sugar().Errorw("validate", "AppID", in.GetAppID(), "error", err)
		return &npool.CreateNotifChannelsResponse{}, status.Error(codes.Internal, "appID is invalid")
	}

	for _, val := range in.GetEventTypes() {
		switch val {
		case usedfor.UsedFor_WithdrawalRequest:
		case usedfor.UsedFor_WithdrawalCompleted:
		case usedfor.UsedFor_DepositReceived:
		case usedfor.UsedFor_KYCApproved:
		case usedfor.UsedFor_KYCRejected:
		case usedfor.UsedFor_Announcement:
		default:
			logger.Sugar().Errorw("validate", "EventType", val, "error", "EventTypes is invalid")
			return &npool.CreateNotifChannelsResponse{}, status.Error(codes.Internal, "EventTypes is invalid")
		}
	}

	switch in.GetChannel() {
	case channel.NotifChannel_ChannelEmail:
	case channel.NotifChannel_ChannelSMS:
	case channel.NotifChannel_ChannelFrontend:
	default:
		logger.Sugar().Errorw("validate", "Channel", in.GetChannel(), "error", "Channel is invalid")
		return &npool.CreateNotifChannelsResponse{}, status.Error(codes.Internal, "Channel is invalid")
	}

	appInfo, err := appcli.GetApp(ctx, in.GetAppID())
	if err != nil {
		logger.Sugar().Errorw("validate", "error", err.Error())
		return &npool.CreateNotifChannelsResponse{}, status.Error(codes.Internal, err.Error())
	}
	if appInfo == nil {
		logger.Sugar().Errorw("validate", "Channel", in.GetChannel(), "error", "AppID not exists")
		return &npool.CreateNotifChannelsResponse{}, status.Error(codes.Internal, "AppID not exists")
	}
	rows, err := notifchannel1.CreateNotifChannels(
		ctx,
		in.GetAppID(),
		in.GetEventTypes(),
		in.GetChannel(),
	)
	if err != nil {
		logger.Sugar().Errorw("CreateNotifChannels", "error", err)
		return &npool.CreateNotifChannelsResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.CreateNotifChannelsResponse{
		Infos: rows,
	}, nil
}

func (s *Server) GetAppNotifChannels(
	ctx context.Context,
	in *npool.GetAppNotifChannelsRequest,
) (
	*npool.GetAppNotifChannelsResponse,
	error,
) {
	var err error

	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "GetAppNotifChannels")
	defer span.End()

	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	span = commontracer.TraceOffsetLimit(span, int(in.GetOffset()), int(in.GetLimit()))

	if _, err := uuid.Parse(in.GetAppID()); err != nil {
		logger.Sugar().Errorw("validate", "AppID", in.GetAppID(), "error", err)
		return &npool.GetAppNotifChannelsResponse{}, status.Error(codes.Internal, "appID is invalid")
	}

	rows, total, err := notifchannel1.GetNotifChannels(
		ctx,
		in.GetAppID(),
		in.GetOffset(),
		in.GetLimit(),
	)
	if err != nil {
		logger.Sugar().Errorw("GetAppNotifChannels", "error", err)
		return &npool.GetAppNotifChannelsResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.GetAppNotifChannelsResponse{
		Infos: rows,
		Total: total,
	}, nil
}

func (s *Server) DeleteNotifChannel(ctx context.Context, in *npool.DeleteNotifChannelRequest) (*npool.DeleteNotifChannelResponse, error) {
	var err error

	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "DeleteNotifChannel")
	defer span.End()

	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	info, err := mgrcli.GetNotifChannel(ctx, in.GetID())
	if err != nil {
		logger.Sugar().Errorw("validate", "ID", in.GetID(), "error", err)
		return &npool.DeleteNotifChannelResponse{}, status.Error(codes.Internal, err.Error())
	}
	if _, err := uuid.Parse(in.GetAppID()); err != nil {
		logger.Sugar().Errorw("validate", "AppID", in.GetAppID(), "error", err)
		return &npool.DeleteNotifChannelResponse{}, status.Error(codes.Internal, "appID is invalid")
	}

	if info.AppID != in.GetAppID() {
		logger.Sugar().Errorw("validate", "AppID", in.GetAppID(), "error", err)
		return &npool.DeleteNotifChannelResponse{}, status.Error(codes.PermissionDenied, "permission denied")
	}

	row, err := notifchannel1.DeleteNotifChannel(
		ctx,
		in.GetID(),
	)
	if err != nil {
		logger.Sugar().Errorw("DeleteNotifChannel", "error", err)
		return &npool.DeleteNotifChannelResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.DeleteNotifChannelResponse{
		Info: row,
	}, nil
}

func (s *Server) GetNAppNotifChannels(
	ctx context.Context,
	in *npool.GetNAppNotifChannelsRequest,
) (
	*npool.GetNAppNotifChannelsResponse,
	error,
) {
	var err error

	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "GetAppNotifChannels")
	defer span.End()

	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	span = commontracer.TraceOffsetLimit(span, int(in.GetOffset()), int(in.GetLimit()))

	if _, err := uuid.Parse(in.GetTargetAppID()); err != nil {
		logger.Sugar().Errorw("validate", "TargetAppID", in.GetTargetAppID(), "error", err)
		return &npool.GetNAppNotifChannelsResponse{}, status.Error(codes.Internal, "GetTargetAppID is invalid")
	}

	rows, total, err := notifchannel1.GetNotifChannels(
		ctx,
		in.GetTargetAppID(),
		in.GetOffset(),
		in.GetLimit(),
	)
	if err != nil {
		logger.Sugar().Errorw("GetAppNotifChannels", "error", err)
		return &npool.GetNAppNotifChannelsResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.GetNAppNotifChannelsResponse{
		Infos: rows,
		Total: total,
	}, nil
}
