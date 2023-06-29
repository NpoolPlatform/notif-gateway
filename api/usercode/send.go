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
		usercode1.WithAppID(&in.AppID),
		usercode1.WithLangID(&in.AppID, &in.LangID),
		usercode1.WithUserID(&in.AppID, in.UserID),
		usercode1.WithAccount(in.Account),
		usercode1.WithAccountType(&in.AccountType),
		usercode1.WithUsedFor(&in.UsedFor),
		usercode1.WithToUsername(in.ToUsername),
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
