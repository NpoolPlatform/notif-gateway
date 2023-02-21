package frontend

import (
	"context"
	"fmt"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"

	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/NpoolPlatform/message/npool/notif/gw/v1/template/frontend"

	mgrpb "github.com/NpoolPlatform/message/npool/notif/mgr/v1/template/frontend"
	mgrcli "github.com/NpoolPlatform/notif-manager/pkg/client/template/frontend"

	appusermgrcli "github.com/NpoolPlatform/appuser-manager/pkg/client/app"

	applangmwcli "github.com/NpoolPlatform/g11n-middleware/pkg/client/applang"
	applangmgrpb "github.com/NpoolPlatform/message/npool/g11n/mgr/v1/applang"

	cruder "github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	commonpb "github.com/NpoolPlatform/message/npool"
)

//nolint
func validate(ctx context.Context, in *frontend.CreateFrontendTemplateRequest) error {
	if _, err := uuid.Parse(in.GetAppID()); err != nil {
		logger.Sugar().Errorw("validate", "AppID", in.GetAppID())
		return status.Error(codes.InvalidArgument, "AppID is invalid")
	}

	exist, err := appusermgrcli.ExistApp(ctx, in.GetAppID())
	if err != nil {
		logger.Sugar().Errorw("validate", "err", err)
		return status.Error(codes.Internal, err.Error())
	}

	if !exist {
		logger.Sugar().Errorw("validate", "AppID", in.GetAppID())
		return status.Error(codes.InvalidArgument, "AppID is not exist")
	}

	if _, err := uuid.Parse(in.GetTargetLangID()); err != nil {
		logger.Sugar().Errorw("validate", "TargetLangID", in.GetTargetLangID())
		return status.Error(codes.InvalidArgument, "TargetLangID is invalid")
	}

	appLang, err := applangmwcli.GetLangOnly(ctx, &applangmgrpb.Conds{
		AppID: &commonpb.StringVal{
			Op:    cruder.EQ,
			Value: in.GetAppID(),
		},
		LangID: &commonpb.StringVal{
			Op:    cruder.EQ,
			Value: in.GetTargetLangID(),
		},
	})
	if err != nil {
		logger.Sugar().Errorw("validate", "err", err)
		return err
	}
	if appLang == nil {
		return fmt.Errorf("applang not exist")
	}

	switch in.GetUsedFor() {
	case basetypes.UsedFor_WithdrawalRequest:
	case basetypes.UsedFor_WithdrawalCompleted:
	case basetypes.UsedFor_DepositReceived:
	case basetypes.UsedFor_KYCApproved:
	case basetypes.UsedFor_KYCRejected:
	case basetypes.UsedFor_Announcement:
	default:
		return fmt.Errorf("UsedFor is invalid")
	}

	if in.GetTitle() == "" {
		logger.Sugar().Errorw("validate", "Title", in.GetTitle())
		return status.Error(codes.InvalidArgument, "Title is empty")
	}
	if in.GetContent() == "" {
		logger.Sugar().Errorw("validate", "Content", in.GetContent())
		return status.Error(codes.InvalidArgument, "Content is empty")
	}
	exist, err = mgrcli.ExistFrontendTemplateConds(ctx, &mgrpb.Conds{
		AppID: &commonpb.StringVal{
			Op:    cruder.EQ,
			Value: in.GetAppID(),
		},
		LangID: &commonpb.StringVal{
			Op:    cruder.EQ,
			Value: in.GetTargetLangID(),
		},
		UsedFor: &commonpb.Uint32Val{
			Op:    cruder.EQ,
			Value: uint32(in.GetUsedFor()),
		},
	})
	if err != nil {
		logger.Sugar().Errorw("validate", "err", err)
		return status.Error(codes.Internal, err.Error())
	}
	if exist {
		logger.Sugar().Errorw("validate", "Frontend template already exists")
		return status.Error(codes.AlreadyExists, "Frontend template already exists")
	}

	return nil
}
