package logging

import (
	"context"
	"fmt"
	"time"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
)

// Server is an server logging middleware.
func Server(logger log.Logger) middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
			startTime := time.Now()
			reply, err = handler(ctx, req)
			tr, _ := transport.FromContext(ctx)
			info, _ := middleware.FromContext(ctx)
			level, errMsg := extractError(err)
			log.WithContext(ctx, logger).Log(level,
				"kind", "server",
				"component", tr.Kind,
				"method", info.FullMethod,
				"args", extractArgs(req),
				"code", errors.Code(err),
				"error", errMsg,
				"latency", time.Since(startTime).Seconds(),
			)
			return
		}
	}
}

// Client is an client logging middleware.
func Client(logger log.Logger) middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
			startTime := time.Now()
			reply, err = handler(ctx, req)
			tr, _ := transport.FromContext(ctx)
			info, _ := middleware.FromContext(ctx)
			level, errMsg := extractError(err)
			log.WithContext(ctx, logger).Log(level,
				"kind", "client",
				"component", tr.Kind,
				"method", info.FullMethod,
				"args", extractArgs(req),
				"code", errors.Code(err),
				"error", errMsg,
				"latency", time.Since(startTime).Seconds(),
			)
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
func extractError(err error) (log.Level, string) {
	if err != nil {
		return log.LevelError, fmt.Sprintf("%+v", err)
	}
	return log.LevelInfo, ""
}
