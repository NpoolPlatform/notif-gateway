package email

import (
	"context"
	"fmt"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/NpoolPlatform/message/npool/notif/gw/v1/template/email"

	mgrpb "github.com/NpoolPlatform/message/npool/notif/mgr/v1/template/email"
	mgrcli "github.com/NpoolPlatform/notif-manager/pkg/client/template/email"

	appmwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/app"

	applangmwcli "github.com/NpoolPlatform/g11n-middleware/pkg/client/applang"
	applangmgrpb "github.com/NpoolPlatform/message/npool/g11n/mgr/v1/applang"

	cruder "github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	commonpb "github.com/NpoolPlatform/message/npool"
)

//nolint
func validate(ctx context.Context, in *email.CreateEmailTemplateRequest) error {
	if _, err := uuid.Parse(in.GetAppID()); err != nil {
		logger.Sugar().Errorw("validate", "AppID", in.GetAppID())
		return status.Error(codes.InvalidArgument, "AppID is invalid")
	}

	exist, err := appmwcli.ExistApp(ctx, in.GetAppID())
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

	usedFor := false
	for key := range basetypes.UsedFor_value {
		if key == in.UsedFor.String() && in.UsedFor != basetypes.UsedFor_DefaultUsedFor {
			usedFor = true
		}
	}

	if !usedFor {
		logger.Sugar().Errorw("validate", "UsedFor", in.GetUsedFor())
		return status.Error(codes.InvalidArgument, "UsedFor is invalid")
	}

	if in.GetSender() == "" {
		logger.Sugar().Errorw("validate", "Sender", in.GetSender())
		return status.Error(codes.InvalidArgument, "Sender is empty")
	}
	if in.GetSubject() == "" {
		logger.Sugar().Errorw("validate", "Subject", in.GetSubject())
		return status.Error(codes.InvalidArgument, "Subject is empty")
	}
	if in.GetDefaultToUsername() == "" {
		logger.Sugar().Errorw("validate", "DefaultToUsername", in.GetDefaultToUsername())
		return status.Error(codes.InvalidArgument, "DefaultToUsername is empty")
	}

	exist, err = mgrcli.ExistEmailTemplateConds(ctx, &mgrpb.Conds{
		AppID: &commonpb.StringVal{
			Op:    cruder.EQ,
			Value: in.GetAppID(),
		},
		LangID: &commonpb.StringVal{
			Op:    cruder.EQ,
			Value: in.GetTargetLangID(),
		},
		UsedFor: &commonpb.Int32Val{
			Op:    cruder.EQ,
			Value: int32(in.GetUsedFor().Number()),
		},
	})
	if err != nil {
		logger.Sugar().Errorw("validate", "err", err)
		return status.Error(codes.Internal, err.Error())
	}
	if exist {
		logger.Sugar().Errorw("validate", "Email template already exists")
		return status.Error(codes.AlreadyExists, "Email template already exists")
	}

	return nil
}
