package sms

import (
	"context"

	"github.com/NpoolPlatform/message/npool/notif/gw/v1/template/sms"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

type Server struct {
	sms.UnimplementedGatewayServer
}

func Register(server grpc.ServiceRegistrar) {
	sms.RegisterGatewayServer(server, &Server{})
}

func RegisterGateway(mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) error {
	if err := sms.RegisterGatewayHandlerFromEndpoint(context.Background(), mux, endpoint, opts); err != nil {
		return err
	}
	return nil
}
