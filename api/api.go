package api

import (
	"context"

	"github.com/NpoolPlatform/notif-gateway/api/announcement"
	"github.com/NpoolPlatform/notif-gateway/api/announcement/readstate"
	"github.com/NpoolPlatform/notif-gateway/api/announcement/sendstate"
	"github.com/NpoolPlatform/notif-gateway/api/announcement/user"
	"github.com/NpoolPlatform/notif-gateway/api/contact"
	"github.com/NpoolPlatform/notif-gateway/api/notif"
	"github.com/NpoolPlatform/notif-gateway/api/notif/channel"
	"github.com/NpoolPlatform/notif-gateway/api/template/email"
	"github.com/NpoolPlatform/notif-gateway/api/template/frontend"
	"github.com/NpoolPlatform/notif-gateway/api/template/sms"
	"github.com/NpoolPlatform/notif-gateway/api/usercode"

	v1 "github.com/NpoolPlatform/message/npool/notif/gw/v1"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

type Server struct {
	v1.UnimplementedGatewayServer
}

func Register(server grpc.ServiceRegistrar) {
	v1.RegisterGatewayServer(server, &Server{})
	announcement.Register(server)
	readstate.Register(server)
	sendstate.Register(server)
	notif.Register(server)
	contact.Register(server)
	email.Register(server)
	frontend.Register(server)
	sms.Register(server)
	user.Register(server)
	channel.Register(server)
	usercode.Register(server)
}

func RegisterGateway(mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) error {
	if err := v1.RegisterGatewayHandlerFromEndpoint(context.Background(), mux, endpoint, opts); err != nil {
		return err
	}
	if err := announcement.RegisterGateway(mux, endpoint, opts); err != nil {
		return err
	}
	if err := readstate.RegisterGateway(mux, endpoint, opts); err != nil {
		return err
	}
	if err := sendstate.RegisterGateway(mux, endpoint, opts); err != nil {
		return err
	}
	if err := notif.RegisterGateway(mux, endpoint, opts); err != nil {
		return err
	}
	if err := contact.RegisterGateway(mux, endpoint, opts); err != nil {
		return err
	}
	if err := email.RegisterGateway(mux, endpoint, opts); err != nil {
		return err
	}
	if err := frontend.RegisterGateway(mux, endpoint, opts); err != nil {
		return err
	}
	if err := sms.RegisterGateway(mux, endpoint, opts); err != nil {
		return err
	}
	if err := user.RegisterGateway(mux, endpoint, opts); err != nil {
		return err
	}
	if err := channel.RegisterGateway(mux, endpoint, opts); err != nil {
		return err
	}
	if err := usercode.RegisterGateway(mux, endpoint, opts); err != nil {
		return err
	}
	return nil
}
