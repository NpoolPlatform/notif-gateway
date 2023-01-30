package api

import (
	"context"
	"github.com/NpoolPlatform/notif-gateway/api/announcement"
	"github.com/NpoolPlatform/notif-gateway/api/announcement/readstate"

	"github.com/NpoolPlatform/notif-gateway/api/notif"

	v1 "github.com/NpoolPlatform/message/npool/notif/mw/v1"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

type Server struct {
	v1.UnimplementedMiddlewareServer
}

func Register(server grpc.ServiceRegistrar) {
	v1.RegisterMiddlewareServer(server, &Server{})
	announcement.Register(server)
	readstate.Register(server)
	notif.Register(server)
}

func RegisterGateway(mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) error {
	if err := v1.RegisterMiddlewareHandlerFromEndpoint(context.Background(), mux, endpoint, opts); err != nil {
		return err
	}
	if err := announcement.RegisterGateway(mux, endpoint, opts); err != nil {
		return err
	}
	if err := readstate.RegisterGateway(mux, endpoint, opts); err != nil {
		return err
	}
	if err := notif.RegisterGateway(mux, endpoint, opts); err != nil {
		return err
	}
	return nil
}
