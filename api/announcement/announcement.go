//nolint:dupl
package announcement

import (
	"context"
	appcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/app"
	"github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	npoolpb "github.com/NpoolPlatform/message/npool"

	g11ncli "github.com/NpoolPlatform/g11n-middleware/pkg/client/applang"
	g11npb "github.com/NpoolPlatform/message/npool/g11n/mgr/v1/applang"

	mgrpb "github.com/NpoolPlatform/message/npool/notif/mgr/v1/announcement"

	"github.com/NpoolPlatform/message/npool/notif/mgr/v1/channel"

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

//nolint:funlen,gocyclo
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

	_, err = uuid.Parse(in.GetLangID())
	if err != nil {
		logger.Sugar().Errorw("CreateAnnouncement", "LangID", in.GetLangID(), "error", err)
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

	for _, val := range in.GetChannels() {
		switch val {
		case channel.NotifChannel_ChannelEmail:
		case channel.NotifChannel_ChannelSMS:
		default:
			logger.Sugar().Errorw("CreateAnnouncement", "Channel", in.GetChannels(), "error", err)
			return &npool.CreateAnnouncementResponse{}, status.Error(codes.InvalidArgument, "Channel is invalid")
		}
	}

	if in.GetEndAt() == 0 {
		logger.Sugar().Errorw("CreateAnnouncement", "EndAt", in.GetEndAt(), "error", err)
		return &npool.CreateAnnouncementResponse{}, status.Error(codes.InvalidArgument, "EndAt is empty")
	}

	switch in.AnnouncementType {
	case mgrpb.AnnouncementType_AppointUsers:
	case mgrpb.AnnouncementType_AllUsers:
	default:
		logger.Sugar().Errorw("CreateAnnouncement", "AnnouncementType", in.GetAnnouncementType(), "error", err)
		return &npool.CreateAnnouncementResponse{}, status.Error(codes.InvalidArgument, "AnnouncementType is invalid")
	}

	switch in.GetAnnouncementType() {
	case mgrpb.AnnouncementType_AllUsers:
	case mgrpb.AnnouncementType_AppointUsers:
	default:
		logger.Sugar().Errorw("CreateAnnouncement", "AnnouncementType", in.GetAnnouncementType(), "error", err)
		return &npool.CreateAnnouncementResponse{}, status.Error(codes.InvalidArgument, "AnnouncementType is invalid")
	}

	app, err := appcli.GetApp(ctx, in.GetAppID())
	if err != nil {
		logger.Sugar().Errorw("CreateReadState", "error", err)
		return &npool.CreateAnnouncementResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}
	if app == nil {
		logger.Sugar().Errorw(
			"CreateReadState",
			"AppID",
			in.GetAppID(),
			"error",
			"app not exist",
		)
		return &npool.CreateAnnouncementResponse{}, status.Error(codes.InvalidArgument, "app not exist")
	}

	appLang, err := g11ncli.GetLangOnly(ctx, &g11npb.Conds{
		AppID: &npoolpb.StringVal{
			Op:    cruder.EQ,
			Value: in.GetAppID(),
		},
		LangID: &npoolpb.StringVal{
			Op:    cruder.EQ,
			Value: in.GetLangID(),
		},
	})
	if err != nil {
		logger.Sugar().Errorw("CreateReadState", "error", err)
		return &npool.CreateAnnouncementResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}
	if appLang == nil {
		logger.Sugar().Errorw(
			"CreateReadState",
			"AppID",
			in.GetAppID(),
			"error",
			"app lang not exist",
		)
		return &npool.CreateAnnouncementResponse{}, status.Error(codes.InvalidArgument, "app lang not exist")
	}

	info, err := announcement1.CreateAnnouncement(
		ctx,
		in.GetAppID(),
		in.GetLangID(),
		in.GetTitle(),
		in.GetContent(),
		in.GetChannels(),
		in.GetEndAt(),
		in.GetAnnouncementType(),
	)
	if err != nil {
		logger.Sugar().Errorw("CreateAnnouncement", "error", err)
		return &npool.CreateAnnouncementResponse{}, status.Error(codes.Internal, err.Error())
	}
	return &npool.CreateAnnouncementResponse{
		Info: info,
	}, nil
}

//nolint:gocyclo
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

	if in.GetTitle() == "" && in.Title != nil {
		logger.Sugar().Errorw("UpdateAnnouncement", "Title", in.GetTitle(), "error", err)
		return &npool.UpdateAnnouncementResponse{}, status.Error(codes.InvalidArgument, "Title is empty")
	}

	if in.GetContent() == "" && in.Content != nil {
		logger.Sugar().Errorw("UpdateAnnouncement", "Content", in.GetContent(), "error", err)
		return &npool.UpdateAnnouncementResponse{}, status.Error(codes.InvalidArgument, "Content is empty")
	}

	if len(in.GetChannels()) == 0 {
		logger.Sugar().Errorw("UpdateAnnouncement", "Channels", in.GetChannels(), "error", err)
		return &npool.UpdateAnnouncementResponse{}, status.Error(codes.InvalidArgument, "Channels is empty")
	}

	for _, val := range in.GetChannels() {
		switch val {
		case channel.NotifChannel_ChannelEmail:
		case channel.NotifChannel_ChannelSMS:
		default:
			logger.Sugar().Errorw("CreateAnnouncement", "Channel", in.GetChannels(), "error", err)
			return &npool.UpdateAnnouncementResponse{}, status.Error(codes.InvalidArgument, "Channel is invalid")
		}
	}

	if in.GetEndAt() == 0 && in.EndAt != nil {
		logger.Sugar().Errorw("UpdateAnnouncement", "EndAt", in.GetEndAt(), "error", err)
		return &npool.UpdateAnnouncementResponse{}, status.Error(codes.InvalidArgument, "EndAt is empty")
	}

	if in.AnnouncementType != nil {
		switch in.GetAnnouncementType() {
		case mgrpb.AnnouncementType_AllUsers:
		case mgrpb.AnnouncementType_AppointUsers:
		default:
			logger.Sugar().Errorw("CreateAnnouncement", "AnnouncementType", in.GetAnnouncementType(), "error", err)
			return &npool.UpdateAnnouncementResponse{}, status.Error(codes.InvalidArgument, "AnnouncementType is invalid")
		}
	}

	info, err := announcement1.UpdateAnnouncement(
		ctx,
		in.GetID(),
		in.Title,
		in.Content,
		in.GetChannels(),
		in.EndAt,
		in.AnnouncementType,
	)
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

	_, err = uuid.Parse(in.GetUserID())
	if err != nil {
		logger.Sugar().Errorw("GetAnnouncements", "UserID", in.GetUserID(), "error", err)
		return &npool.GetAnnouncementsResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}
	_, err = uuid.Parse(in.GetUserID())
	if err != nil {
		logger.Sugar().Errorw("GetAnnouncements", "UserID", in.GetUserID(), "error", err)
		return &npool.GetAnnouncementsResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}
	_, err = uuid.Parse(in.GetLangID())
	if err != nil {
		logger.Sugar().Errorw("GetAnnouncements", "LangID", in.GetLangID(), "error", err)
		return &npool.GetAnnouncementsResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	infos, total, err := announcement1.GetAnnouncements(
		ctx,
		in.GetAppID(),
		in.GetUserID(),
		in.GetLangID(),
		in.GetOffset(),
		in.GetLimit(),
	)
	if err != nil {
		logger.Sugar().Errorw("GetAnnouncements", "error", err)
		return &npool.GetAnnouncementsResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.GetAnnouncementsResponse{
		Infos: infos,
		Total: total,
	}, nil
}

func (s *Server) GetAppAnnouncements(
	ctx context.Context,
	in *npool.GetAppAnnouncementsRequest,
) (
	*npool.GetAppAnnouncementsResponse,
	error,
) {
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
		return &npool.GetAppAnnouncementsResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	infos, total, err := announcement1.GetAppAnnouncements(ctx, in.GetAppID(), in.GetOffset(), in.GetLimit())
	if err != nil {
		logger.Sugar().Errorw("GetAnnouncements", "error", err)
		return &npool.GetAppAnnouncementsResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.GetAppAnnouncementsResponse{
		Infos: infos,
		Total: total,
	}, nil
}

func (s *Server) GetNAppAnnouncements(
	ctx context.Context,
	in *npool.GetNAppAnnouncementsRequest,
) (
	*npool.GetNAppAnnouncementsResponse,
	error,
) {
	var err error

	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "GetAnnouncements")
	defer span.End()

	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	_, err = uuid.Parse(in.GetTargetAppID())
	if err != nil {
		logger.Sugar().Errorw("GetAnnouncements", "TargetAppID", in.GetTargetAppID(), "error", err)
		return &npool.GetNAppAnnouncementsResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	infos, total, err := announcement1.GetAppAnnouncements(ctx, in.GetTargetAppID(), in.GetOffset(), in.GetLimit())
	if err != nil {
		logger.Sugar().Errorw("GetAnnouncements", "error", err)
		return &npool.GetNAppAnnouncementsResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.GetNAppAnnouncementsResponse{
		Infos: infos,
		Total: total,
	}, nil
}
