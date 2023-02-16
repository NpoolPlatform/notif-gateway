package usercode

import (
	"context"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"

	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"
	npool "github.com/NpoolPlatform/message/npool/notif/gw/v1/usercode"

	usercode1 "github.com/NpoolPlatform/notif-gateway/pkg/usercode"

	constant "github.com/NpoolPlatform/notif-gateway/pkg/message/const"
	commontracer "github.com/NpoolPlatform/notif-gateway/pkg/tracer"

	"go.opentelemetry.io/otel"
	scodes "go.opentelemetry.io/otel/codes"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/google/uuid"
)

func (s *Server) SendCode(ctx context.Context, in *npool.SendCodeRequest) (*npool.SendCodeResponse, error) { //nolint
	var err error

	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "CreateContact")
	defer span.End()

	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	if _, err := uuid.Parse(in.GetAppID()); err != nil {
		logger.Sugar().Errorw("SendCode", "AppID", in.GetAppID())
		return &npool.SendCodeResponse{}, status.Error(codes.InvalidArgument, "AppID is invalid")
	}

	switch in.GetUsedFor() {
	case basetypes.UsedFor_Signup:
		fallthrough //nolint
	case basetypes.UsedFor_Signin:
		fallthrough //nolint
	case basetypes.UsedFor_Update:
		fallthrough //nolint
	case basetypes.UsedFor_SetWithdrawAddress:
		fallthrough //nolint
	case basetypes.UsedFor_Withdraw:
		fallthrough //nolint
	case basetypes.UsedFor_CreateInvitationCode:
		fallthrough //nolint
	case basetypes.UsedFor_SetCommission:
		fallthrough //nolint
	case basetypes.UsedFor_SetTransferTargetUser:
		fallthrough //nolint
	case basetypes.UsedFor_Transfer:
		if _, err := uuid.Parse(in.GetUserID()); err != nil {
			logger.Sugar().Errorw("SendCode", "UserID", in.GetUserID())
			return &npool.SendCodeResponse{}, status.Error(codes.InvalidArgument, "UserID is invalid")
		}
	default:
		logger.Sugar().Errorw("SendCode", "UsedFor", in.GetUsedFor())
		return &npool.SendCodeResponse{}, status.Error(codes.InvalidArgument, "UsedFor is invalid")
	}

	if in.Account == nil && in.UserID == nil {
		logger.Sugar().Errorw("SendCode", "Account", in.GetAccount(), "UserID", in.GetUserID())
		return &npool.SendCodeResponse{}, status.Error(codes.InvalidArgument, "Account and UserID cannot all be empty")
	}

	switch in.GetAccountType() {
	case basetypes.SignMethod_Email:
	case basetypes.SignMethod_Mobile:
	default:
		logger.Sugar().Errorw("SendCode", "AccountType", in.GetAccountType())
		return &npool.SendCodeResponse{}, status.Error(codes.InvalidArgument, "AccountType is invalid")
	}

	span = commontracer.TraceInvoker(span, "contact", "manager", "CreateEmailTemplate")

	err = usercode1.SendCode(
		ctx,
		in.GetAppID(),
		in.GetLangID(),
		in.UserID,
		in.Account,
		in.GetAccountType(),
		in.GetUsedFor(),
		in.ToUsername,
	)
	if err != nil {
		logger.Sugar().Errorw("SendCode", "err", err)
		return &npool.SendCodeResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	return &npool.SendCodeResponse{}, nil
}
