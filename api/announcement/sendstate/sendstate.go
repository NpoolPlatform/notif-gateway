//nolint:dupl
package sendstate

import (
	"context"

	"github.com/NpoolPlatform/message/npool/notif/mgr/v1/channel"

	constant "github.com/NpoolPlatform/notif-gateway/pkg/message/const"

	"go.opentelemetry.io/otel"
	scodes "go.opentelemetry.io/otel/codes"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	npool "github.com/NpoolPlatform/message/npool/notif/gw/v1/announcement/sendstate"

	"github.com/google/uuid"

	sendstate1 "github.com/NpoolPlatform/notif-gateway/pkg/announcement/sendstate"
)

func (s *Server) GetSendStates(
	ctx context.Context,
	in *npool.GetSendStatesRequest,
) (
	*npool.GetSendStatesResponse,
	error,
) {
	var err error

	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "GetSendStates")
	defer span.End()

	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	_, err = uuid.Parse(in.GetAppID())
	if err != nil {
		logger.Sugar().Errorw("GetSendStates", "AppID", in.GetAppID(), "error", err)
		return &npool.GetSendStatesResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}
	_, err = uuid.Parse(in.GetUserID())
	if err != nil {
		logger.Sugar().Errorw("GetSendStates", "UserID", in.GetUserID(), "error", err)
		return &npool.GetSendStatesResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	if in.Channel != nil {
		switch in.GetChannel() {
		case channel.NotifChannel_ChannelEmail:
		case channel.NotifChannel_ChannelSMS:
		default:
			logger.Sugar().Errorw("GetSendStates", "Channel", in.GetChannel(), "error", err)
			return &npool.GetSendStatesResponse{}, status.Error(codes.InvalidArgument, "Channel is invalid")
		}
	}

	infos, total, err := sendstate1.GetSendStates(
		ctx,
		in.GetAppID(),
		in.GetUserID(),
		in.GetOffset(),
		in.GetLimit(),
		in.Channel,
	)
	if err != nil {
		logger.Sugar().Errorw("GetSendStates", "error", err)
		return &npool.GetSendStatesResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.GetSendStatesResponse{
		Infos: infos,
		Total: total,
	}, nil
}

func (s *Server) GetAppUserSendStates(
	ctx context.Context,
	in *npool.GetAppUserSendStatesRequest,
) (
	*npool.GetAppUserSendStatesResponse,
	error,
) {
	var err error

	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "GetAppUserSendStates")
	defer span.End()

	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	_, err = uuid.Parse(in.GetTargetAppID())
	if err != nil {
		logger.Sugar().Errorw("GetAppUserSendStates", "TargetAppID", in.GetTargetAppID(), "error", err)
		return &npool.GetAppUserSendStatesResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}
	_, err = uuid.Parse(in.GetTargetUserID())
	if err != nil {
		logger.Sugar().Errorw("GetAppUserSendStates", "TargetUserID", in.GetTargetUserID(), "error", err)
		return &npool.GetAppUserSendStatesResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	if in.Channel != nil {
		switch in.GetChannel() {
		case channel.NotifChannel_ChannelEmail:
		case channel.NotifChannel_ChannelSMS:
		default:
			logger.Sugar().Errorw("GetSendStates", "Channel", in.GetChannel(), "error", err)
			return &npool.GetAppUserSendStatesResponse{}, status.Error(codes.InvalidArgument, "Channel is invalid")
		}
	}

	infos, total, err := sendstate1.GetSendStates(
		ctx,
		in.GetTargetAppID(),
		in.GetTargetUserID(),
		in.GetOffset(),
		in.GetLimit(),
		in.Channel,
	)
	if err != nil {
		logger.Sugar().Errorw("GetAppUserSendStates", "error", err)
		return &npool.GetAppUserSendStatesResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.GetAppUserSendStatesResponse{
		Infos: infos,
		Total: total,
	}, nil
}

func (s *Server) GetAppSendStates(
	ctx context.Context,
	in *npool.GetAppSendStatesRequest,
) (
	*npool.GetAppSendStatesResponse,
	error,
) {
	var err error

	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "GetAppSendStates")
	defer span.End()

	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	_, err = uuid.Parse(in.GetAppID())
	if err != nil {
		logger.Sugar().Errorw("GetAppSendStates", "AppID", in.GetAppID(), "error", err)
		return &npool.GetAppSendStatesResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	if in.Channel != nil {
		switch in.GetChannel() {
		case channel.NotifChannel_ChannelEmail:
		case channel.NotifChannel_ChannelSMS:
		default:
			logger.Sugar().Errorw("GetSendStates", "Channel", in.GetChannel(), "error", err)
			return &npool.GetAppSendStatesResponse{}, status.Error(codes.InvalidArgument, "Channel is invalid")
		}
	}

	infos, total, err := sendstate1.GetAppSendStates(
		ctx,
		in.GetAppID(),
		in.GetOffset(),
		in.GetLimit(),
		in.Channel,
	)
	if err != nil {
		logger.Sugar().Errorw("GetAppSendStates", "error", err)
		return &npool.GetAppSendStatesResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.GetAppSendStatesResponse{
		Infos: infos,
		Total: total,
	}, nil
}

func (s *Server) GetNAppSendStates(
	ctx context.Context,
	in *npool.GetNAppSendStatesRequest,
) (
	*npool.GetNAppSendStatesResponse,
	error,
) {
	var err error

	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "GetNAppSendStates")
	defer span.End()

	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	_, err = uuid.Parse(in.GetTargetAppID())
	if err != nil {
		logger.Sugar().Errorw("GetNAppSendStates", "TargetAppID", in.GetTargetAppID(), "error", err)
		return &npool.GetNAppSendStatesResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	if in.Channel != nil {
		switch in.GetChannel() {
		case channel.NotifChannel_ChannelEmail:
		case channel.NotifChannel_ChannelSMS:
		default:
			logger.Sugar().Errorw("GetSendStates", "Channel", in.GetChannel(), "error", err)
			return &npool.GetNAppSendStatesResponse{}, status.Error(codes.InvalidArgument, "Channel is invalid")
		}
	}

	infos, total, err := sendstate1.GetAppSendStates(
		ctx,
		in.GetTargetAppID(),
		in.GetOffset(),
		in.GetLimit(),
		in.Channel,
	)
	if err != nil {
		logger.Sugar().Errorw("GetNAppSendStates", "error", err)
		return &npool.GetNAppSendStatesResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.GetNAppSendStatesResponse{
		Infos: infos,
		Total: total,
	}, nil
}
