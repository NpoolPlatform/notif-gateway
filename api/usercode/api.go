package usercode

import (
	"context"

	"github.com/NpoolPlatform/message/npool/notif/gw/v1/usercode"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

type Server struct {
	usercode.UnimplementedGatewayServer
}

func Register(server grpc.ServiceRegistrar) {
	usercode.RegisterGatewayServer(server, &Server{})
}

func RegisterGateway(mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) error {
	if err := usercode.RegisterGatewayHandlerFromEndpoint(context.Background(), mux, endpoint, opts); err != nil {
		return err
	}
	return nil
}
