package contact

import (
	"context"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	npool "github.com/NpoolPlatform/message/npool/notif/gw/v1/contact"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	contact1 "github.com/NpoolPlatform/notif-gateway/pkg/contact"
)

func (s *Server) CreateContact(ctx context.Context, in *npool.CreateContactRequest) (*npool.CreateContactResponse, error) {
	handler, err := contact1.NewHandler(
		ctx,
		contact1.WithAppID(&in.AppID, true),
		contact1.WithAccount(&in.Account, true),
		contact1.WithAccountType(&in.AccountType, true),
		contact1.WithUsedFor(&in.UsedFor, true),
		contact1.WithSender(&in.Sender, true),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"CreateContact",
			"In", in,
			"Error", err,
		)
		return &npool.CreateContactResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	info, err := handler.CreateContact(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"CreateContact",
			"In", in,
			"Error", err,
		)
		return &npool.CreateContactResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.CreateContactResponse{
		Info: info,
	}, nil
}

func (s *Server) CreateAppContact(ctx context.Context, in *npool.CreateAppContactRequest) (*npool.CreateAppContactResponse, error) {
	resp, err := s.CreateContact(ctx, &npool.CreateContactRequest{
		AppID:       in.TargetAppID,
		Account:     in.Account,
		AccountType: in.AccountType,
		UsedFor:     in.UsedFor,
		Sender:      in.Sender,
	})
	if err != nil {
		logger.Sugar().Errorw(
			"CreateAppContact",
			"In", in,
			"Error", err,
		)
		return &npool.CreateAppContactResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	return &npool.CreateAppContactResponse{
		Info: resp.Info,
	}, nil
}
