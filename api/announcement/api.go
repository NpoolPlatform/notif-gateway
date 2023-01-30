package announcement

import (
	"context"

	announcement "github.com/NpoolPlatform/message/npool/notif/gw/v1/announcement"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

type Server struct {
	announcement.UnimplementedGatewayServer
}

func Register(server grpc.ServiceRegistrar) {
	announcement.RegisterGatewayServer(server, &Server{})
}

func RegisterGateway(mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) error {
	return announcement.RegisterGatewayHandlerFromEndpoint(context.Background(), mux, endpoint, opts)
}
