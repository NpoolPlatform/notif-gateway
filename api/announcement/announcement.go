//nolint:dupl
package announcement

import (
	"context"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	npool "github.com/NpoolPlatform/message/npool/notif/gw/v1/announcement"
	constant "github.com/NpoolPlatform/notif-gateway/pkg/message/const"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	scodes "go.opentelemetry.io/otel/codes"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	announcement1 "github.com/NpoolPlatform/notif-gateway/pkg/announcement"
)

func (s *Server) CreateAnnouncement(ctx context.Context, in *npool.CreateAnnouncementRequest) (*npool.CreateAnnouncementResponse, error) {
	var err error

	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "CreateAnnouncement")
	defer span.End()

	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	_, err = uuid.Parse(in.GetAppID())
	if err != nil {
		logger.Sugar().Errorw("CreateAnnouncement", "AppID", in.GetAppID(), "error", err)
		return &npool.CreateAnnouncementResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	if in.GetTitle() == "" {
		logger.Sugar().Errorw("CreateAnnouncement", "Title", in.GetTitle(), "error", err)
		return &npool.CreateAnnouncementResponse{}, status.Error(codes.InvalidArgument, "Title is empty")
	}

	if in.GetContent() == "" {
		logger.Sugar().Errorw("CreateAnnouncement", "Content", in.GetContent(), "error", err)
		return &npool.CreateAnnouncementResponse{}, status.Error(codes.InvalidArgument, "Content is empty")
	}

	if len(in.GetChannels()) == 0 {
		logger.Sugar().Errorw("CreateAnnouncement", "Channels", in.GetChannels(), "error", err)
		return &npool.CreateAnnouncementResponse{}, status.Error(codes.InvalidArgument, "Channels is empty")
	}

	info, err := announcement1.CreateAnnouncement(ctx, in.GetAppID(), in.GetTitle(), in.GetContent(), in.GetChannels())
	if err != nil {
		logger.Sugar().Errorw("CreateAnnouncement", "error", err)
		return &npool.CreateAnnouncementResponse{}, status.Error(codes.Internal, err.Error())
	}
	return &npool.CreateAnnouncementResponse{
		Info: info,
	}, nil
}

func (s *Server) UpdateAnnouncement(ctx context.Context, in *npool.UpdateAnnouncementRequest) (*npool.UpdateAnnouncementResponse, error) {
	var err error

	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "UpdateAnnouncement")
	defer span.End()

	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	_, err = uuid.Parse(in.GetID())
	if err != nil {
		logger.Sugar().Errorw("UpdateAnnouncement", "ID", in.GetID(), "error", err)
		return &npool.UpdateAnnouncementResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	if in.GetTitle() == "" && in.Title == nil {
		logger.Sugar().Errorw("UpdateAnnouncement", "Title", in.GetTitle(), "error", err)
		return &npool.UpdateAnnouncementResponse{}, status.Error(codes.InvalidArgument, "Title is empty")
	}

	if in.GetContent() == "" && in.Content == nil {
		logger.Sugar().Errorw("UpdateAnnouncement", "Content", in.GetContent(), "error", err)
		return &npool.UpdateAnnouncementResponse{}, status.Error(codes.InvalidArgument, "Content is empty")
	}

	if len(in.GetChannels()) == 0 {
		logger.Sugar().Errorw("UpdateAnnouncement", "Channels", in.GetChannels(), "error", err)
		return &npool.UpdateAnnouncementResponse{}, status.Error(codes.InvalidArgument, "Channels is empty")
	}

	info, err := announcement1.UpdateAnnouncement(ctx, in.GetID(), in.Title, in.Content, in.GetChannels())
	if err != nil {
		logger.Sugar().Errorw("UpdateAnnouncement", "error", err)
		return &npool.UpdateAnnouncementResponse{}, status.Error(codes.Internal, err.Error())
	}
	return &npool.UpdateAnnouncementResponse{
		Info: info,
	}, nil
}

func (s *Server) DeleteAnnouncement(ctx context.Context, in *npool.DeleteAnnouncementRequest) (*npool.DeleteAnnouncementResponse, error) {
	var err error

	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "DeleteAnnouncement")
	defer span.End()

	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	_, err = uuid.Parse(in.GetID())
	if err != nil {
		logger.Sugar().Errorw("DeleteAnnouncement", "ID", in.GetID(), "error", err)
		return &npool.DeleteAnnouncementResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	info, err := announcement1.DeleteAnnouncement(ctx, in.GetID())
	if err != nil {
		logger.Sugar().Errorw("DeleteAnnouncement", "error", err)
		return &npool.DeleteAnnouncementResponse{}, status.Error(codes.Internal, err.Error())
	}
	return &npool.DeleteAnnouncementResponse{
		Info: info,
	}, nil
}

func (s *Server) GetAnnouncement(ctx context.Context, in *npool.GetAnnouncementRequest) (*npool.GetAnnouncementResponse, error) {
	var err error

	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "GetAnnouncement")
	defer span.End()

	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	_, err = uuid.Parse(in.GetID())
	if err != nil {
		logger.Sugar().Errorw("GetAnnouncements", "D", in.GetID(), "error", err)
		return &npool.GetAnnouncementResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	info, err := announcement1.GetAnnouncement(ctx, in.GetID())
	if err != nil {
		logger.Sugar().Errorw("GetAnnouncements", "error", err)
		return &npool.GetAnnouncementResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.GetAnnouncementResponse{
		Info: info,
	}, nil
}

func (s *Server) GetAnnouncements(ctx context.Context, in *npool.GetAnnouncementsRequest) (*npool.GetAnnouncementsResponse, error) {
	var err error

	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "GetAnnouncements")
	defer span.End()

	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	_, err = uuid.Parse(in.GetAppID())
	if err != nil {
		logger.Sugar().Errorw("GetAnnouncements", "AppID", in.GetAppID(), "error", err)
		return &npool.GetAnnouncementsResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	infos, total, err := announcement1.GetAnnouncements(ctx, in.GetAppID(), in.GetOffset(), in.GetLimit())
	if err != nil {
		logger.Sugar().Errorw("GetAnnouncements", "error", err)
		return &npool.GetAnnouncementsResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.GetAnnouncementsResponse{
		Infos: infos,
		Total: total,
	}, nil
}
