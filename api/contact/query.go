package contact

import (
	"context"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	npool "github.com/NpoolPlatform/message/npool/notif/gw/v1/contact"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	contact1 "github.com/NpoolPlatform/notif-gateway/pkg/contact"
)

func (s *Server) GetContact(ctx context.Context, in *npool.GetContactRequest) (*npool.GetContactResponse, error) {
	handler, err := contact1.NewHandler(
		ctx,
		contact1.WithEntID(&in.EntID, true),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"GetContact",
			"In", in,
			"Error", err,
		)
		return &npool.GetContactResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	info, err := handler.GetContact(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"GetContact",
			"In", in,
			"Error", err,
		)
		return &npool.GetContactResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.GetContactResponse{
		Info: info,
	}, nil
}

func (s *Server) GetContacts(ctx context.Context, in *npool.GetContactsRequest) (*npool.GetContactsResponse, error) {
	handler, err := contact1.NewHandler(
		ctx,
		contact1.WithAppID(&in.AppID, true),
		contact1.WithOffset(in.Offset),
		contact1.WithLimit(in.Limit),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"GetContacts",
			"In", in,
			"Error", err,
		)
		return &npool.GetContactsResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	infos, total, err := handler.GetContacts(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"GetContacts",
			"In", in,
			"Error", err,
		)
		return &npool.GetContactsResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.GetContactsResponse{
		Infos: infos,
		Total: total,
	}, nil
}

func (s *Server) GetAppContacts(ctx context.Context, in *npool.GetAppContactsRequest) (*npool.GetAppContactsResponse, error) {
	resp, err := s.GetContacts(ctx, &npool.GetContactsRequest{
		AppID:  in.TargetAppID,
		Offset: in.Offset,
		Limit:  in.Limit,
	})

	if err != nil {
		logger.Sugar().Errorw(
			"GetAppContacts",
			"In", in,
			"Error", err,
		)
		return &npool.GetAppContactsResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	return &npool.GetAppContactsResponse{
		Infos: resp.Infos,
		Total: resp.Total,
	}, nil
}
