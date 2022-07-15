package broker

import "github.com/go-kratos/kratos/v2/broker"

func WithConsumeHandler(handler broker.ConsumeHandler) ServerOption {
	return func(o *Server) {
		o.consumeHandler = handler
	}
}
