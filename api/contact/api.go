package contact

import (
	"context"

	"github.com/NpoolPlatform/message/npool/notif/gw/v1/contact"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

type Server struct {
	contact.UnimplementedGatewayServer
}

func Register(server grpc.ServiceRegistrar) {
	contact.RegisterGatewayServer(server, &Server{})
}

func RegisterGateway(mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) error {
	if err := contact.RegisterGatewayHandlerFromEndpoint(context.Background(), mux, endpoint, opts); err != nil {
		return err
	}
	return nil
}
