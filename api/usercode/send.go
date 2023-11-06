package usercode

import (
	"context"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"

	npool "github.com/NpoolPlatform/message/npool/notif/gw/v1/usercode"

	usercode1 "github.com/NpoolPlatform/notif-gateway/pkg/usercode"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) SendCode(ctx context.Context, in *npool.SendCodeRequest) (*npool.SendCodeResponse, error) {
	handler, err := usercode1.NewHandler(
		ctx,
		usercode1.WithAppID(&in.AppID, true),
		usercode1.WithLangID(&in.LangID, true),
		usercode1.WithUserID(in.UserID, false),
		usercode1.WithAccount(in.Account, false),
		usercode1.WithAccountType(&in.AccountType, true),
		usercode1.WithUsedFor(&in.UsedFor, true),
		usercode1.WithToUsername(in.ToUsername, true),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"SendCode",
			"In", in,
			"Error", err,
		)
		return &npool.SendCodeResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	err = handler.SendCode(ctx)
	if err != nil {
		logger.Sugar().Errorw("SendCode", "err", err)
		return &npool.SendCodeResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	return &npool.SendCodeResponse{}, nil
}
