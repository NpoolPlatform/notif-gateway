package notif

import (
	"context"
	notif "github.com/NpoolPlatform/message/npool/notif/gw/v1/notif"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

type Server struct {
	notif.UnimplementedGatewayServer
}

func Register(server grpc.ServiceRegistrar) {
	notif.RegisterGatewayServer(server, &Server{})
}

func RegisterGateway(mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) error {
	return notif.RegisterGatewayHandlerFromEndpoint(context.Background(), mux, endpoint, opts)
}
