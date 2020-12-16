package grpc

import (
	"github.com/go-kratos/kratos/v2/middleware"
)

// ServerOption is gRPC server option.
type ServerOption func(o *serverOptions)

type serverOptions struct {
	middleware middleware.Middleware
}

// ServerMiddleware .
func ServerMiddleware(m middleware.Middleware) ServerOption {
	return func(o *serverOptions) { o.middleware = m }
}
