package frontend

import (
	"context"

	"github.com/NpoolPlatform/message/npool/notif/gw/v1/template/frontend"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

type Server struct {
	frontend.UnimplementedGatewayServer
}

func Register(server grpc.ServiceRegistrar) {
	frontend.RegisterGatewayServer(server, &Server{})
}

func RegisterGateway(mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) error {
	if err := frontend.RegisterGatewayHandlerFromEndpoint(context.Background(), mux, endpoint, opts); err != nil {
		return err
	}
	return nil
}
