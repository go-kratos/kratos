package grpc

import (
	"context"

	"github.com/go-kratos/kratos/v2/middleware"
)

// ServerOption is gRPC server option.
type ServerOption func(o *serverOptions)

type serverOptions struct {
	middleware   middleware.Middleware
	errorEncoder EncodeErrorFunc
}

// EncodeErrorFunc .
type EncodeErrorFunc func(ctx context.Context, err error) error

// ServerMiddleware .
func ServerMiddleware(m ...middleware.Middleware) ServerOption {
	return func(o *serverOptions) {
		o.middleware = middleware.Chain(m[0], m[1:]...)
	}
}
