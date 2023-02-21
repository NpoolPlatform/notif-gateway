package channel

import (
	"context"

	channel "github.com/NpoolPlatform/message/npool/notif/gw/v1/notif/channel"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

type Server struct {
	channel.UnimplementedGatewayServer
}

func Register(server grpc.ServiceRegistrar) {
	channel.RegisterGatewayServer(server, &Server{})
}

func RegisterGateway(mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) error {
	return channel.RegisterGatewayHandlerFromEndpoint(context.Background(), mux, endpoint, opts)
}
