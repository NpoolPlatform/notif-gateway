package sendstate

import (
	"context"

	sendstate "github.com/NpoolPlatform/message/npool/notif/gw/v1/announcement/sendstate"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

type Server struct {
	sendstate.UnimplementedGatewayServer
}

func Register(server grpc.ServiceRegistrar) {
	sendstate.RegisterGatewayServer(server, &Server{})
}

func RegisterGateway(mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) error {
	return sendstate.RegisterGatewayHandlerFromEndpoint(context.Background(), mux, endpoint, opts)
}
