package email

import (
	"context"

	"github.com/NpoolPlatform/message/npool/notif/gw/v1/template/email"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

type Server struct {
	email.UnimplementedGatewayServer
}

func Register(server grpc.ServiceRegistrar) {
	email.RegisterGatewayServer(server, &Server{})
}

func RegisterGateway(mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) error {
	if err := email.RegisterGatewayHandlerFromEndpoint(context.Background(), mux, endpoint, opts); err != nil {
		return err
	}
	return nil
}
