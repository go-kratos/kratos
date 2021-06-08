package logging

import (
	"context"
	"fmt"

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

// extractArgs returns the string of the req
func extractArgs(req interface{}) string {
	if stringer, ok := req.(fmt.Stringer); ok {
		return stringer.String()
	}
	return fmt.Sprintf("%+v", req)
}

// extractError returns the string of the error
func extractError(err error) (errMsg string) {
	if err != nil {
		errMsg = fmt.Sprintf("%+v", err)
	}
	return
}
