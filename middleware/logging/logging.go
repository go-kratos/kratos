package logging

import (
	"context"
	"fmt"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
)

// Server is an server logging middleware.
func Server(l log.Logger) middleware.Middleware {
	logger := log.NewHelper("middleware/logging", l)
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			var (
				path      string
				method    string
				params    string
				component string
			)
			if info, ok := http.FromServerContext(ctx); ok {
				component = "HTTP"
				path = info.Request.RequestURI
				method = info.Request.Method
				params = info.Request.Form.Encode()
			} else if info, ok := grpc.FromServerContext(ctx); ok {
				component = "gRPC"
				path = info.FullMethod
				method = "POST"
				params = req.(fmt.Stringer).String()
			}
			reply, err := handler(ctx, req)
			if err != nil {
				logger.Errorw(
					"kind", "server",
					"component", component,
					"path", path,
					"method", method,
					"params", params,
					"code", errors.Code(err),
					"error", err.Error(),
				)
				return nil, err
			}
			logger.Infow(
				"kind", "server",
				"component", component,
				"path", path,
				"method", method,
				"params", params,
				"code", 0,
			)
			return reply, nil
		}
	}
}

// Client is an client logging middleware.
func Client(l log.Logger) middleware.Middleware {
	logger := log.NewHelper("middleware/logging", l)
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			var (
				path      string
				method    string
				params    string
				component string
			)
			if info, ok := http.FromClientContext(ctx); ok {
				component = "HTTP"
				path = info.Request.RequestURI
				method = info.Request.Method
				params = info.Request.Form.Encode()
			} else if info, ok := grpc.FromClientContext(ctx); ok {
				path = info.FullMethod
				method = "POST"
				component = "gRPC"
				params = req.(fmt.Stringer).String()
			}
			reply, err := handler(ctx, req)
			if err != nil {
				logger.Errorw(
					"kind", "client",
					"component", component,
					"path", path,
					"method", method,
					"params", params,
					"code", errors.Code(err),
					"error", err.Error(),
				)
				return nil, err
			}
			logger.Infow(
				"kind", "client",
				"component", component,
				"path", path,
				"method", method,
				"params", params,
				"code", 0,
			)
			return reply, nil
		}
	}
}
