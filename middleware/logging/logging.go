package logging

import (
	"context"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
)

// Option is HTTP logging option.
type Option func(*options)

type options struct {
	logger log.Logger
}

// WithLogger with middleware logger.
func WithLogger(logger log.Logger) Option {
	return func(o *options) {
		o.logger = logger
	}
}

// Server is an server logging middleware.
func Server(opts ...Option) middleware.Middleware {
	options := options{
		logger: log.DefaultLogger,
	}
	for _, o := range opts {
		o(&options)
	}
	log := log.NewHelper("middleware/logging", options.logger)
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			var (
				path   string
				method string
			)
			if info, ok := http.FromServerContext(ctx); ok {
				path = info.Request.RequestURI
				method = info.Request.Method
			} else if info, ok := grpc.FromServerContext(ctx); ok {
				path = info.FullMethod
				method = "POST"
			}
			reply, err := handler(ctx, req)
			if err != nil {
				log.Errorw(
					"kind", "server",
					"path", path,
					"method", method,
					"code", errors.Code(err),
					"error", err.Error(),
				)
				return nil, err
			}
			log.Infow(
				"kind", "server",
				"path", path,
				"method", method,
				"code", 0,
			)
			return reply, nil
		}
	}
}

// Client is an client logging middleware.
func Client(opts ...Option) middleware.Middleware {
	options := options{
		logger: log.DefaultLogger,
	}
	for _, o := range opts {
		o(&options)
	}
	log := log.NewHelper("middleware/logging", options.logger)
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			var (
				component string
				path      string
				method    string
			)
			if info, ok := http.FromClientContext(ctx); ok {
				component = "HTTP"
				path = info.Request.RequestURI
				method = info.Request.Method
			} else if info, ok := grpc.FromClientContext(ctx); ok {
				component = "gRPC"
				path = info.FullMethod
				method = "POST"
			}
			reply, err := handler(ctx, req)
			if err != nil {
				log.Errorw(
					"kind", "client",
					"component", component,
					"path", path,
					"method", method,
					"code", errors.Code(err),
					"error", err.Error(),
				)
				return nil, err
			}
			log.Infow(
				"kind", "client",
				"component", component,
				"path", path,
				"method", method,
				"code", 0,
			)
			return reply, nil
		}
	}
}
