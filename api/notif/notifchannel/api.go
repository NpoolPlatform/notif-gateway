package notifchannel

import (
	"context"

	notifchannel "github.com/NpoolPlatform/message/npool/notif/gw/v1/notif/notifchannel"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

type Server struct {
	notifchannel.UnimplementedGatewayServer
}

func Register(server grpc.ServiceRegistrar) {
	notifchannel.RegisterGatewayServer(server, &Server{})
}

func RegisterGateway(mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) error {
	return notifchannel.RegisterGatewayHandlerFromEndpoint(context.Background(), mux, endpoint, opts)
}
