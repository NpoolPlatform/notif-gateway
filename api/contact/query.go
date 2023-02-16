//nolint:dupl
package contact

import (
	"context"

	"go.opentelemetry.io/otel/attribute"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	"github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	npool "github.com/NpoolPlatform/message/npool/notif/gw/v1/contact"
	constant "github.com/NpoolPlatform/notif-gateway/pkg/message/const"
	commontracer "github.com/NpoolPlatform/notif-gateway/pkg/tracer"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	scodes "go.opentelemetry.io/otel/codes"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	contactpb "github.com/NpoolPlatform/message/npool/notif/mgr/v1/contact"
	"github.com/NpoolPlatform/notif-manager/pkg/client/contact"

	npoolpb "github.com/NpoolPlatform/message/npool"
)

func (s *Server) GetContact(ctx context.Context, in *npool.GetContactRequest) (*npool.GetContactResponse, error) {
	var err error

	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "GetContact")
	defer span.End()

	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	commontracer.TraceID(span, in.GetID())

	if _, err := uuid.Parse(in.GetID()); err != nil {
		logger.Sugar().Errorw("validate", "ID", in.GetID())
		return &npool.GetContactResponse{}, status.Error(codes.InvalidArgument, "ID is invalid")
	}

	span = commontracer.TraceInvoker(span, "contact", "manager", "GetContact")

	info, err := contact.GetContact(ctx, in.GetID())
	if err != nil {
		logger.Sugar().Errorw("validate", "err", err)
		return &npool.GetContactResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.GetContactResponse{
		Info: info,
	}, nil
}

func (s *Server) GetContacts(ctx context.Context, in *npool.GetContactsRequest) (*npool.GetContactsResponse, error) {
	var err error

	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "GetContact")
	defer span.End()

	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	span.SetAttributes(
		attribute.String("AppID", in.GetAppID()),
	)
	span = commontracer.TraceOffsetLimit(span, int(in.GetOffset()), int(in.GetLimit()))

	if _, err := uuid.Parse(in.GetAppID()); err != nil {
		logger.Sugar().Errorw("validate", "AppID", in.GetAppID())
		return &npool.GetContactsResponse{}, status.Error(codes.InvalidArgument, "AppID is invalid")
	}

	span = commontracer.TraceInvoker(span, "contact", "manager", "GetContacts")

	infos, total, err := contact.GetContacts(ctx, &contactpb.Conds{
		AppID: &npoolpb.StringVal{
			Op:    cruder.EQ,
			Value: in.GetAppID(),
		},
	}, in.GetOffset(), in.GetLimit())
	if err != nil {
		logger.Sugar().Errorw("validate", "err", err)
		return &npool.GetContactsResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.GetContactsResponse{
		Infos: infos,
		Total: total,
	}, nil
}

func (s *Server) GetAppContacts(ctx context.Context, in *npool.GetAppContactsRequest) (*npool.GetAppContactsResponse, error) {
	var err error

	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "GetContact")
	defer span.End()

	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	span.SetAttributes(
		attribute.String("AppID", in.GetTargetAppID()),
	)
	span = commontracer.TraceOffsetLimit(span, int(in.GetOffset()), int(in.GetLimit()))

	if _, err := uuid.Parse(in.GetTargetAppID()); err != nil {
		logger.Sugar().Errorw("validate", "TargetAppID", in.GetTargetAppID())
		return &npool.GetAppContactsResponse{}, status.Error(codes.InvalidArgument, "TargetAppID is invalid")
	}

	span = commontracer.TraceInvoker(span, "contact", "manager", "GetAppContacts")

	infos, total, err := contact.GetContacts(ctx, &contactpb.Conds{
		AppID: &npoolpb.StringVal{
			Op:    cruder.EQ,
			Value: in.GetTargetAppID(),
		},
	}, in.GetOffset(), in.GetLimit())
	if err != nil {
		logger.Sugar().Errorw("validate", "err", err)
		return &npool.GetAppContactsResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.GetAppContactsResponse{
		Infos: infos,
		Total: total,
	}, nil
}
