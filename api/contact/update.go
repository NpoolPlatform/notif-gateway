package contact

import (
	"context"

	"github.com/google/uuid"

	tracer "github.com/NpoolPlatform/notif-manager/pkg/tracer/contact"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	npool "github.com/NpoolPlatform/message/npool/notif/gw/v1/contact"
	contactpb "github.com/NpoolPlatform/message/npool/notif/mgr/v1/contact"
	constant "github.com/NpoolPlatform/notif-gateway/pkg/message/const"
	commontracer "github.com/NpoolPlatform/notif-gateway/pkg/tracer"
	"github.com/NpoolPlatform/notif-manager/pkg/client/contact"
	"go.opentelemetry.io/otel"
	scodes "go.opentelemetry.io/otel/codes"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) UpdateContact(ctx context.Context, in *npool.UpdateContactRequest) (*npool.UpdateContactResponse, error) {
	var err error

	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "UpdateContact")
	defer span.End()

	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	if _, err := uuid.Parse(in.GetID()); err != nil {
		logger.Sugar().Errorw("validate", "ID", in.GetID())
		return &npool.UpdateContactResponse{}, status.Error(codes.InvalidArgument, "ID is invalid")
	}

	info, err := contact.GetContact(ctx, in.GetID())
	if err != nil {
		logger.Sugar().Errorw("validate", "err", err)
		return &npool.UpdateContactResponse{}, status.Error(codes.Internal, err.Error())
	}

	if info.GetAppID() != in.GetAppID() {
		logger.Sugar().Errorw("validate", "err", err)
		return &npool.UpdateContactResponse{}, status.Error(codes.PermissionDenied, "permission denied")
	}

	contactInfo := &contactpb.ContactReq{
		ID:          &in.ID,
		Account:     in.Account,
		AccountType: in.AccountType,
		Sender:      in.Sender,
	}
	tracer.Trace(span, contactInfo)

	span = commontracer.TraceInvoker(span, "contact", "manager", "UpdateContact")

	info, err = contact.UpdateContact(ctx, contactInfo)
	if err != nil {
		logger.Sugar().Errorw("validate", "err", err)
		return &npool.UpdateContactResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.UpdateContactResponse{
		Info: info,
	}, nil
}

func (s *Server) UpdateAppContact(ctx context.Context, in *npool.UpdateAppContactRequest) (*npool.UpdateAppContactResponse, error) {
	var err error

	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "UpdateAppContact")
	defer span.End()

	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	if _, err := uuid.Parse(in.GetID()); err != nil {
		logger.Sugar().Errorw("validate", "ID", in.GetID())
		return &npool.UpdateAppContactResponse{}, status.Error(codes.InvalidArgument, "ID is invalid")
	}

	contactInfo := &contactpb.ContactReq{
		ID:          &in.ID,
		Account:     in.Account,
		AccountType: in.AccountType,
		Sender:      in.Sender,
	}
	tracer.Trace(span, contactInfo)

	span = commontracer.TraceInvoker(span, "contact", "manager", "UpdateContact")

	info, err := contact.UpdateContact(ctx, contactInfo)
	if err != nil {
		logger.Sugar().Errorw("validate", "err", err)
		return &npool.UpdateAppContactResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.UpdateAppContactResponse{
		Info: info,
	}, nil
}
