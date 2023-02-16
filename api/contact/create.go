package contact

import (
	"context"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	npool "github.com/NpoolPlatform/message/npool/notif/gw/v1/contact"
	constant "github.com/NpoolPlatform/notif-gateway/pkg/message/const"
	"go.opentelemetry.io/otel"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	scodes "go.opentelemetry.io/otel/codes"

	commontracer "github.com/NpoolPlatform/notif-gateway/pkg/tracer"

	contactpb "github.com/NpoolPlatform/message/npool/notif/mgr/v1/contact"
	"github.com/NpoolPlatform/notif-manager/pkg/client/contact"

	tracer "github.com/NpoolPlatform/notif-manager/pkg/tracer/contact"
)

func (s *Server) CreateContact(ctx context.Context, in *npool.CreateContactRequest) (*npool.CreateContactResponse, error) {
	var err error

	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "CreateContact")
	defer span.End()

	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	contactInfo := &contactpb.ContactReq{
		AppID:       &in.AppID,
		Account:     &in.Account,
		AccountType: &in.AccountType,
		UsedFor:     &in.UsedFor,
		Sender:      &in.Sender,
	}

	tracer.Trace(span, contactInfo)

	err = validate(ctx, in)
	if err != nil {
		return nil, err
	}

	span = commontracer.TraceInvoker(span, "contact", "manager", "CreateContact")

	info, err := contact.CreateContact(ctx, contactInfo)
	if err != nil {
		logger.Sugar().Errorw("validate", "err", err)
		return &npool.CreateContactResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.CreateContactResponse{
		Info: info,
	}, nil
}

func (s *Server) CreateAppContact(ctx context.Context, in *npool.CreateAppContactRequest) (*npool.CreateAppContactResponse, error) {
	var err error

	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "CreateAppContact")
	defer span.End()

	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	contactInfo := &contactpb.ContactReq{
		AppID:       &in.TargetAppID,
		Account:     &in.Account,
		AccountType: &in.AccountType,
		UsedFor:     &in.UsedFor,
		Sender:      &in.Sender,
	}

	tracer.Trace(span, contactInfo)

	err = validate(ctx, &npool.CreateContactRequest{
		AppID:       in.TargetAppID,
		Account:     in.Account,
		AccountType: in.AccountType,
		UsedFor:     in.UsedFor,
		Sender:      in.Sender,
	})
	if err != nil {
		return nil, err
	}

	span = commontracer.TraceInvoker(span, "contact", "manager", "CreateContact")

	info, err := contact.CreateContact(ctx, contactInfo)
	if err != nil {
		logger.Sugar().Errorw("validate", "err", err)
		return &npool.CreateAppContactResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.CreateAppContactResponse{
		Info: info,
	}, nil
}
