package contact

import (
	"context"

	"github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	"github.com/NpoolPlatform/message/npool"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"

	"github.com/NpoolPlatform/message/npool/notif/gw/v1/contact"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/google/uuid"

	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"

	mgrpb "github.com/NpoolPlatform/message/npool/notif/mgr/v1/contact"
	mgrcli "github.com/NpoolPlatform/notif-manager/pkg/client/contact"
)

func validate(ctx context.Context, info *contact.CreateContactRequest) error {
	if _, err := uuid.Parse(info.GetAppID()); err != nil {
		logger.Sugar().Errorw("validate", "AppID", info.GetAppID())
		return status.Error(codes.InvalidArgument, "AppID is invalid")
	}

	switch info.GetUsedFor() {
	case basetypes.UsedFor_Contact:
	default:
		logger.Sugar().Errorw("validate", "UsedFor", info.GetUsedFor())
		return status.Error(codes.InvalidArgument, "UsedFor is invalid")
	}

	if info.GetAccount() == "" {
		logger.Sugar().Errorw("validate", "Account", info.GetAccount())
		return status.Error(codes.InvalidArgument, "Account is empty")
	}

	switch info.GetAccountType() {
	case basetypes.SignMethod_Email:
	default:
		logger.Sugar().Errorw("validate", "AccountType", info.GetAccountType())
		return status.Error(codes.InvalidArgument, "AccountType is invalid")
	}

	if info.GetSender() == "" {
		logger.Sugar().Errorw("validate", "Sender", info.GetSender())
		return status.Error(codes.InvalidArgument, "Sender is empty")
	}

	exist, err := mgrcli.ExistContactConds(ctx, &mgrpb.Conds{
		AppID: &npool.StringVal{
			Op:    cruder.EQ,
			Value: info.GetAppID(),
		},
		AccountType: &npool.Int32Val{
			Op:    cruder.EQ,
			Value: int32(info.GetAccountType().Number()),
		},
		UsedFor: &npool.Int32Val{
			Op:    cruder.EQ,
			Value: int32(info.GetUsedFor().Number()),
		},
	})
	if err != nil {
		logger.Sugar().Errorw("validate", "err", err)
		return status.Error(codes.Internal, err.Error())
	}

	if exist {
		logger.Sugar().Errorw("validate", "Contact already exists")
		return status.Error(codes.AlreadyExists, "Contact already exists")
	}
	return nil
}
