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
func Server(logger log.Logger) middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
			reply, err = handler(ctx, req)
			if tr, ok := transport.FromContext(ctx); ok {
				switch tr.Kind {
				case transport.KindHTTP:
					httpServerLog(logger, ctx, extractArgs(req), err)
				case transport.KindGRPC:
					grpcServerLog(logger, ctx, extractArgs(req), err)
				}
			}
			return
		}
	}
}

// Client is an client logging middleware.
func Client(logger log.Logger) middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
			reply, err = handler(ctx, req)
			if tr, ok := transport.FromContext(ctx); ok {
				switch tr.Kind {
				case transport.KindHTTP:
					httpClientLog(logger, ctx, extractArgs(req), err)
				case transport.KindGRPC:
					grpcClientLog(logger, ctx, extractArgs(req), err)
				}
			}
			return
		}
	}
}

func extractArgs(req interface{}) string {
	if stringer, ok := req.(fmt.Stringer); ok {
		return stringer.String()
	}
	return fmt.Sprintf("%+v", req)
}

func extractTrace(ctx context.Context) (traceID, spanID string) {
	span := trace.SpanContextFromContext(ctx)
	if span.HasTraceID() {
		traceID = span.TraceID().String()
	}
	if span.HasSpanID() {
		spanID = span.SpanID().String()
	}
	return
}

func extractError(err error) (code int, errMsg string) {
	if err != nil {
		code = errors.Code(err)
		errMsg = fmt.Sprintf("%+v", err)
	}
	return
}
