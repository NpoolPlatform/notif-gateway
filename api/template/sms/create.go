package sms

import (
	"context"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	npool "github.com/NpoolPlatform/message/npool/notif/gw/v1/template/sms"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	smstemplate1 "github.com/NpoolPlatform/notif-gateway/pkg/template/sms"
)

func (s *Server) CreateSMSTemplate(
	ctx context.Context,
	in *npool.CreateSMSTemplateRequest,
) (
	*npool.CreateSMSTemplateResponse,
	error,
) {
	handler, err := smstemplate1.NewHandler(
		ctx,
		smstemplate1.WithAppID(&in.AppID),
		smstemplate1.WithLangID(&in.TargetLangID),
		smstemplate1.WithUsedFor(&in.UsedFor),
		smstemplate1.WithSubject(&in.Subject),
		smstemplate1.WithMessage(&in.Message),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"CreateSMSTemplate",
			"In", in,
			"Error", err,
		)
		return &npool.CreateSMSTemplateResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	info, err := handler.CreateSMSTemplate(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"CreateSMSTemplate",
			"In", in,
			"Error", err,
		)
		return &npool.CreateSMSTemplateResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.CreateSMSTemplateResponse{
		Info: info,
	}, nil
}

func (s *Server) CreateAppSMSTemplate(
	ctx context.Context,
	in *npool.CreateAppSMSTemplateRequest,
) (
	*npool.CreateAppSMSTemplateResponse,
	error,
) {
	handler, err := smstemplate1.NewHandler(
		ctx,
		smstemplate1.WithAppID(&in.TargetAppID),
		smstemplate1.WithLangID(&in.TargetLangID),
		smstemplate1.WithUsedFor(&in.UsedFor),
		smstemplate1.WithSubject(&in.Subject),
		smstemplate1.WithMessage(&in.Message),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"CreateAppSMSTemplate",
			"In", in,
			"Error", err,
		)
		return &npool.CreateAppSMSTemplateResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	info, err := handler.CreateSMSTemplate(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"CreateAppSMSTemplate",
			"In", in,
			"Error", err,
		)
		return &npool.CreateAppSMSTemplateResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.CreateAppSMSTemplateResponse{
		Info: info,
	}, nil
}
