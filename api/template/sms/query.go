//nolint:dupl
package sms

import (
	"context"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	npool "github.com/NpoolPlatform/message/npool/notif/gw/v1/template/sms"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	smstemplate1 "github.com/NpoolPlatform/notif-gateway/pkg/template/sms"
)

func (s *Server) GetSMSTemplate(ctx context.Context, in *npool.GetSMSTemplateRequest) (*npool.GetSMSTemplateResponse, error) {
	handler, err := smstemplate1.NewHandler(
		ctx,
		smstemplate1.WithEntID(&in.EntID, true),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"GetSMSTemplate",
			"In", in,
			"Error", err,
		)
		return &npool.GetSMSTemplateResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	info, err := handler.GetSMSTemplate(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"GetSMSTemplate",
			"In", in,
			"Error", err,
		)
		return &npool.GetSMSTemplateResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.GetSMSTemplateResponse{
		Info: info,
	}, nil
}

func (s *Server) GetSMSTemplates(ctx context.Context, in *npool.GetSMSTemplatesRequest) (*npool.GetSMSTemplatesResponse, error) {
	handler, err := smstemplate1.NewHandler(
		ctx,
		smstemplate1.WithAppID(&in.AppID, true),
		smstemplate1.WithOffset(in.GetOffset()),
		smstemplate1.WithLimit(in.GetLimit()),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"GetSMSTemplates",
			"In", in,
			"Error", err,
		)
		return &npool.GetSMSTemplatesResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	infos, total, err := handler.GetSMSTemplates(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"GetSMSTemplates",
			"In", in,
			"Error", err,
		)
		return &npool.GetSMSTemplatesResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.GetSMSTemplatesResponse{
		Infos: infos,
		Total: total,
	}, nil
}

func (s *Server) GetAppSMSTemplates(ctx context.Context, in *npool.GetAppSMSTemplatesRequest) (*npool.GetAppSMSTemplatesResponse, error) {
	handler, err := smstemplate1.NewHandler(
		ctx,
		smstemplate1.WithAppID(&in.TargetAppID, true),
		smstemplate1.WithOffset(in.GetOffset()),
		smstemplate1.WithLimit(in.GetLimit()),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"GetAppSMSTemplates",
			"In", in,
			"Error", err,
		)
		return &npool.GetAppSMSTemplatesResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	infos, total, err := handler.GetSMSTemplates(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"GetAppSMSTemplates",
			"In", in,
			"Error", err,
		)
		return &npool.GetAppSMSTemplatesResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.GetAppSMSTemplatesResponse{
		Infos: infos,
		Total: total,
	}, nil
}
