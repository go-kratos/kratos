package logging

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel/trace"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
)

// Server is an server logging middleware.
func Server(l log.Logger) middleware.Middleware {
	logger := log.NewHelper(l)
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			var (
				path      string
				method    string
				args      string
				component string
				query     string
				traceID   string
			)
			if tid := trace.SpanContextFromContext(ctx).TraceID(); tid.IsValid() {
				traceID = tid.String()
			}
			if stringer, ok := req.(fmt.Stringer); ok {
				args = stringer.String()
			} else {
				args = fmt.Sprintf("%+v", req)
			}
			if info, ok := http.FromServerContext(ctx); ok {
				component = "HTTP"
				path = info.Request.URL.Path
				method = info.Request.Method
				query = info.Request.URL.RawQuery
			} else if info, ok := grpc.FromServerContext(ctx); ok {
				component = "gRPC"
				path = info.FullMethod
				method = "POST"
			}
			reply, err := handler(ctx, req)
			if component == "HTTP" {
				if err != nil {
					logger.Errorw(
						"kind", "server",
						"component", component,
						"traceID", traceID,
						"path", path,
						"method", method,
						"args", args,
						"query", query,
						"code", uint32(errors.Code(err)),
						"error", err.Error(),
					)
					return nil, err
				}
				logger.Infow(
					"kind", "server",
					"component", component,
					"traceID", traceID,
					"path", path,
					"method", method,
					"args", args,
					"query", query,
					"code", 0,
				)
			} else {
				if err != nil {
					logger.Errorw(
						"kind", "server",
						"component", component,
						"traceID", traceID,
						"path", path,
						"method", method,
						"args", args,
						"code", uint32(errors.Code(err)),
						"error", err.Error(),
					)
					return nil, err
				}
				logger.Infow(
					"kind", "server",
					"component", component,
					"traceID", traceID,
					"path", path,
					"method", method,
					"args", args,
					"code", 0,
				)
			}
			return reply, nil
		}
	}
}

// Client is an client logging middleware.
func Client(l log.Logger) middleware.Middleware {
	logger := log.NewHelper(l)
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			var (
				path      string
				method    string
				args      string
				component string
				query     string
				traceID   string
			)
			if tid := trace.SpanContextFromContext(ctx).TraceID(); tid.IsValid() {
				traceID = tid.String()
			}
			if info, ok := http.FromClientContext(ctx); ok {
				component = "HTTP"
				path = info.Request.URL.Path
				method = info.Request.Method
				args = req.(fmt.Stringer).String()
				query = info.Request.URL.RawQuery
			} else if info, ok := grpc.FromClientContext(ctx); ok {
				path = info.FullMethod
				method = "POST"
				component = "gRPC"
				args = req.(fmt.Stringer).String()
			}
			reply, err := handler(ctx, req)
			if component == "HTTP" {
				if err != nil {
					logger.Errorw(
						"kind", "client",
						"component", component,
						"traceID", traceID,
						"path", path,
						"method", method,
						"args", args,
						"query", query,
						"code", uint32(errors.Code(err)),
						"error", err.Error(),
					)
					return nil, err
				}
				logger.Infow(
					"kind", "client",
					"component", component,
					"traceID", traceID,
					"path", path,
					"method", method,
					"args", args,
					"query", query,
					"code", 0,
				)
			} else {
				if err != nil {
					logger.Errorw(
						"kind", "client",
						"component", component,
						"traceID", traceID,
						"path", path,
						"method", method,
						"args", args,
						"code", uint32(errors.Code(err)),
						"error", err.Error(),
					)
					return nil, err
				}
				logger.Infow(
					"kind", "client",
					"component", component,
					"traceID", traceID,
					"path", path,
					"method", method,
					"args", args,
					"code", 0,
				)
			}
			return reply, nil
		}
	}
}
