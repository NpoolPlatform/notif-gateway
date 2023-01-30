package readstate

import (
	"context"

	readstate "github.com/NpoolPlatform/message/npool/notif/gw/v1/announcement/readstate"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

type Server struct {
	readstate.UnimplementedGatewayServer
}

func Register(server grpc.ServiceRegistrar) {
	readstate.RegisterGatewayServer(server, &Server{})
}

func RegisterGateway(mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) error {
	return readstate.RegisterGatewayHandlerFromEndpoint(context.Background(), mux, endpoint, opts)
}
