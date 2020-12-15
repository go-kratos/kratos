package grpc

import "github.com/go-kratos/kratos/v2/middleware"

// ServerOption is gRPC server option.
type ServerOption func(o *serverOptions)

type serverOptions struct {
	middlewares []middleware.Middleware
}

func ServerMiddleware(m ...middleware.Middleware) ServerOption {
	return func(o *serverOptions) { o.middlewares = append(o.middlewares, m...) }
}
