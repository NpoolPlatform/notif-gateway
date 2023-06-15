package sendstate

import (
	"context"
	"fmt"

	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"
	"github.com/NpoolPlatform/notif-gateway/pkg/announcement/handler"
)

type Handler struct {
	*handler.Handler
	Channel *basetypes.NotifChannel
}

func NewHandler(ctx context.Context, options ...interface{}) (*Handler, error) {
	_handler, err := handler.NewHandler(ctx, options...)
	if err != nil {
		return nil, err
	}
	h := &Handler{
		Handler: _handler,
	}

	for _, opt := range options {
		_opt, ok := opt.(func(context.Context, *Handler) error)
		if !ok {
			continue
		}
		if err := _opt(ctx, h); err != nil {
			return nil, err
		}
	}
	return h, nil
}

func WithChannel(channel *basetypes.NotifChannel) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if channel == nil {
			return nil
		}
		switch *channel {
		case basetypes.NotifChannel_ChannelEmail:
		case basetypes.NotifChannel_ChannelSMS:
		case basetypes.NotifChannel_ChannelFrontend:
		default:
			return fmt.Errorf("channel %v invalid", *channel)
		}
		h.Channel = channel
		return nil
	}
}
