//nolint:nolintlint,dupl
package sms

import (
	"context"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	npool "github.com/NpoolPlatform/message/npool/notif/gw/v1/template/sms"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	smstemplate1 "github.com/NpoolPlatform/notif-gateway/pkg/template/sms"
)

func (s *Server) UpdateSMSTemplate(
	ctx context.Context,
	in *npool.UpdateSMSTemplateRequest,
) (
	*npool.UpdateSMSTemplateResponse,
	error,
) {
	handler, err := smstemplate1.NewHandler(
		ctx,
		smstemplate1.WithID(&in.ID, true),
		smstemplate1.WithEntID(&in.EntID, true),
		smstemplate1.WithAppID(&in.AppID, true),
		smstemplate1.WithSubject(in.Subject, false),
		smstemplate1.WithMessage(in.Message, false),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"UpdateSMSTemplate",
			"In", in,
			"Error", err,
		)
		return &npool.UpdateSMSTemplateResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	info, err := handler.UpdateSMSTemplate(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"UpdateSMSTemplate",
			"In", in,
			"Error", err,
		)
		return &npool.UpdateSMSTemplateResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.UpdateSMSTemplateResponse{
		Info: info,
	}, nil
}

func (s *Server) UpdateAppSMSTemplate(
	ctx context.Context,
	in *npool.UpdateAppSMSTemplateRequest,
) (
	*npool.UpdateAppSMSTemplateResponse,
	error,
) {
	handler, err := smstemplate1.NewHandler(
		ctx,
		smstemplate1.WithID(&in.ID, true),
		smstemplate1.WithEntID(&in.EntID, true),
		smstemplate1.WithAppID(&in.TargetAppID, true),
		smstemplate1.WithSubject(in.Subject, false),
		smstemplate1.WithMessage(in.Message, false),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"UpdateAppSMSTemplate",
			"In", in,
			"Error", err,
		)
		return &npool.UpdateAppSMSTemplateResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	info, err := handler.UpdateSMSTemplate(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"UpdateAppSMSTemplate",
			"In", in,
			"Error", err,
		)
		return &npool.UpdateAppSMSTemplateResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.UpdateAppSMSTemplateResponse{
		Info: info,
	}, nil
}
