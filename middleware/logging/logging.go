package logging

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel/trace"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
)

// Server is an server logging middleware.
func Server(l log.Logger) middleware.Middleware {
	logger := log.NewHelper("middleware/logging", l)
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
			traceID = trace.SpanContextFromContext(ctx).TraceID().String()
			if stringer, ok := req.(fmt.Stringer); ok {
				args = stringer.String()
			} else {
				args = fmt.Sprintf("%+v", req)
			}
			if tr, ok := transport.FromContext(ctx); ok {
				component = string(tr.Kind)
				path = tr.Request.FullPath
				method = tr.Request.Method
				query = tr.Request.Query
			}

			reply, err := handler(ctx, req)

			kvPairs := []interface{}{
				"kind", "server",
				"component", component,
				"traceID", traceID,
				"path", path,
				"method", method,
				"args", args,
				"query", query,
				"code", uint32(errors.Code(err)),
			}
			if err != nil {
				kvPairs = append(kvPairs, "error", err.Error())
				logger.Errorw(kvPairs...)
			} else {
				logger.Infow(kvPairs...)
			}
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
				args      string
				component string
				query     string
				traceID   string
			)
			traceID = trace.SpanContextFromContext(ctx).TraceID().String()
			if stringer, ok := req.(fmt.Stringer); ok {
				args = stringer.String()
			} else {
				args = fmt.Sprintf("%+v", req)
			}
			if tr, ok := transport.FromContext(ctx); ok {
				component = string(tr.Kind)
				path = tr.Request.FullPath
				method = tr.Request.Method
				query = tr.Request.Query
			}

			reply, err := handler(ctx, req)

			kvPairs := []interface{}{
				"kind", "client",
				"component", component,
				"traceID", traceID,
				"path", path,
				"method", method,
				"args", args,
				"query", query,
				"code", uint32(errors.Code(err)),
			}
			if err != nil {
				kvPairs = append(kvPairs, "error", err.Error())
				logger.Errorw(kvPairs...)
			} else {
				logger.Infow(kvPairs...)
			}
			return reply, nil
		}
	}
}
