//nolint:dupl
package channel

import (
	"context"

	appcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/app"
	mgrcli "github.com/NpoolPlatform/notif-manager/pkg/client/notif/channel"

	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"
	"github.com/NpoolPlatform/message/npool/notif/mgr/v1/channel"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"

	npool "github.com/NpoolPlatform/message/npool/notif/gw/v1/notif/channel"
	constant "github.com/NpoolPlatform/notif-gateway/pkg/message/const"

	commontracer "github.com/NpoolPlatform/notif-gateway/pkg/tracer"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	scodes "go.opentelemetry.io/otel/codes"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	channel1 "github.com/NpoolPlatform/notif-gateway/pkg/notif/channel"
)

//nolint:gocyclo
func (s *Server) CreateChannels(
	ctx context.Context,
	in *npool.CreateChannelsRequest,
) (
	*npool.CreateChannelsResponse,
	error,
) {
	var err error

	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "CreateChannels")
	defer span.End()

	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	if _, err := uuid.Parse(in.GetAppID()); err != nil {
		logger.Sugar().Errorw("validate", "AppID", in.GetAppID(), "error", err)
		return &npool.CreateChannelsResponse{}, status.Error(codes.Internal, "appID is invalid")
	}

	for _, val := range in.GetEventTypes() {
		switch val {
		case basetypes.UsedFor_WithdrawalRequest:
		case basetypes.UsedFor_WithdrawalCompleted:
		case basetypes.UsedFor_DepositReceived:
		case basetypes.UsedFor_KYCApproved:
		case basetypes.UsedFor_KYCRejected:
		case basetypes.UsedFor_Announcement:
		default:
			logger.Sugar().Errorw("validate", "EventType", val, "error", "EventTypes is invalid")
			return &npool.CreateChannelsResponse{}, status.Error(codes.Internal, "EventTypes is invalid")
		}
	}

	switch in.GetChannel() {
	case channel.NotifChannel_ChannelEmail:
	case channel.NotifChannel_ChannelSMS:
	case channel.NotifChannel_ChannelFrontend:
	default:
		logger.Sugar().Errorw("validate", "Channel", in.GetChannel(), "error", "Channel is invalid")
		return &npool.CreateChannelsResponse{}, status.Error(codes.Internal, "Channel is invalid")
	}

	appInfo, err := appcli.GetApp(ctx, in.GetAppID())
	if err != nil {
		logger.Sugar().Errorw("validate", "error", err.Error())
		return &npool.CreateChannelsResponse{}, status.Error(codes.Internal, err.Error())
	}
	if appInfo == nil {
		logger.Sugar().Errorw("validate", "Channel", in.GetChannel(), "error", "AppID not exists")
		return &npool.CreateChannelsResponse{}, status.Error(codes.Internal, "AppID not exists")
	}
	rows, err := channel1.CreateChannels(
		ctx,
		in.GetAppID(),
		in.GetEventTypes(),
		in.GetChannel(),
	)
	if err != nil {
		logger.Sugar().Errorw("CreateChannels", "error", err)
		return &npool.CreateChannelsResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.CreateChannelsResponse{
		Infos: rows,
	}, nil
}

func (s *Server) GetAppChannels(
	ctx context.Context,
	in *npool.GetAppChannelsRequest,
) (
	*npool.GetAppChannelsResponse,
	error,
) {
	var err error

	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "GetAppChannels")
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
		return &npool.GetAppChannelsResponse{}, status.Error(codes.Internal, "appID is invalid")
	}

	rows, total, err := channel1.GetChannels(
		ctx,
		in.GetAppID(),
		in.GetOffset(),
		in.GetLimit(),
	)
	if err != nil {
		logger.Sugar().Errorw("GetAppChannels", "error", err)
		return &npool.GetAppChannelsResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.GetAppChannelsResponse{
		Infos: rows,
		Total: total,
	}, nil
}

func (s *Server) DeleteChannel(ctx context.Context, in *npool.DeleteChannelRequest) (*npool.DeleteChannelResponse, error) {
	var err error

	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "DeleteChannel")
	defer span.End()

	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	info, err := mgrcli.GetChannel(ctx, in.GetID())
	if err != nil {
		logger.Sugar().Errorw("validate", "ID", in.GetID(), "error", err)
		return &npool.DeleteChannelResponse{}, status.Error(codes.Internal, err.Error())
	}
	if _, err := uuid.Parse(in.GetAppID()); err != nil {
		logger.Sugar().Errorw("validate", "AppID", in.GetAppID(), "error", err)
		return &npool.DeleteChannelResponse{}, status.Error(codes.Internal, "appID is invalid")
	}

	if info.AppID != in.GetAppID() {
		logger.Sugar().Errorw("validate", "AppID", in.GetAppID(), "error", err)
		return &npool.DeleteChannelResponse{}, status.Error(codes.PermissionDenied, "permission denied")
	}

	row, err := channel1.DeleteChannel(
		ctx,
		in.GetID(),
	)
	if err != nil {
		logger.Sugar().Errorw("DeleteChannel", "error", err)
		return &npool.DeleteChannelResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.DeleteChannelResponse{
		Info: row,
	}, nil
}

func (s *Server) GetNAppChannels(
	ctx context.Context,
	in *npool.GetNAppChannelsRequest,
) (
	*npool.GetNAppChannelsResponse,
	error,
) {
	var err error

	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "GetAppChannels")
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
		return &npool.GetNAppChannelsResponse{}, status.Error(codes.Internal, "GetTargetAppID is invalid")
	}

	rows, total, err := channel1.GetChannels(
		ctx,
		in.GetTargetAppID(),
		in.GetOffset(),
		in.GetLimit(),
	)
	if err != nil {
		logger.Sugar().Errorw("GetAppChannels", "error", err)
		return &npool.GetNAppChannelsResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.GetNAppChannelsResponse{
		Infos: rows,
		Total: total,
	}, nil
}
