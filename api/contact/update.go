package contact

import (
	"context"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	npool "github.com/NpoolPlatform/message/npool/notif/gw/v1/contact"
	contact1 "github.com/NpoolPlatform/notif-gateway/pkg/contact"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) UpdateContact(ctx context.Context, in *npool.UpdateContactRequest) (*npool.UpdateContactResponse, error) {
	handler, err := contact1.NewHandler(
		ctx,
		contact1.WithID(&in.ID),
		contact1.WithAppID(&in.AppID),
		contact1.WithSender(in.Sender),
		contact1.WithAccount(in.Account),
		contact1.WithAccountType(in.AccountType),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"UpdateContact",
			"In", in,
			"Error", err,
		)
		return &npool.UpdateContactResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	info, err := handler.UpdateContact(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"UpdateContact",
			"In", in,
			"Error", err,
		)
		return &npool.UpdateContactResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.UpdateContactResponse{
		Info: info,
	}, nil
}

func (s *Server) UpdateAppContact(ctx context.Context, in *npool.UpdateAppContactRequest) (*npool.UpdateAppContactResponse, error) {
	resp, err := s.UpdateContact(ctx, &npool.UpdateContactRequest{
		ID:          in.ID,
		AppID:       in.TargetAppID,
		Account:     in.Account,
		AccountType: in.AccountType,
		Sender:      in.Sender,
	})
	if err != nil {
		logger.Sugar().Errorw(
			"UpdateAppContact",
			"In", in,
			"Error", err,
		)
		return &npool.UpdateAppContactResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.UpdateAppContactResponse{
		Info: resp.Info,
	}, nil
}
