package logging

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"google.golang.org/grpc/codes"

	"github.com/go-kratos/kratos/v3/errors"
	"github.com/go-kratos/kratos/v3/middleware"
	"github.com/go-kratos/kratos/v3/transport"
	"github.com/go-kratos/kratos/v3/transport/http/status"
)

// Redacter defines how to log an object
type Redacter interface {
	Redact() string
}

// Server is an server logging middleware.
func Server(logger *slog.Logger) middleware.Middleware {
	if logger == nil {
		logger = slog.Default()
	}
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req any) (reply any, err error) {
			var (
				code      int32
				reason    string
				kind      string
				operation string
			)

			// default code
			code = int32(status.FromGRPCCode(codes.OK))

			startTime := time.Now()
			if info, ok := transport.FromServerContext(ctx); ok {
				kind = info.Kind().String()
				operation = info.Operation()
			}
			reply, err = handler(ctx, req)
			if se := errors.FromError(err); se != nil {
				code = se.Code
				reason = se.Reason
			}
			level, stack := extractError(err)
			attrs := []slog.Attr{
				slog.String("kind", "server"),
				slog.String("component", kind),
				slog.String("operation", operation),
				slog.String("args", extractArgs(req)),
				slog.Int64("code", int64(code)),
				slog.String("reason", reason),
				slog.Float64("latency", time.Since(startTime).Seconds()),
			}
			if err != nil {
				attrs = append(attrs, slog.Any("error", err))
				if stack != "" {
					attrs = append(attrs, slog.String("stack", stack))
				}
			}
			logger.LogAttrs(ctx, level, "server request", attrs...)
			return
		}
	}
}

// Client is a client logging middleware.
func Client(logger *slog.Logger) middleware.Middleware {
	if logger == nil {
		logger = slog.Default()
	}
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req any) (reply any, err error) {
			var (
				code      int32
				reason    string
				kind      string
				operation string
			)

			// default code
			code = int32(status.FromGRPCCode(codes.OK))

			startTime := time.Now()
			if info, ok := transport.FromClientContext(ctx); ok {
				kind = info.Kind().String()
				operation = info.Operation()
			}
			reply, err = handler(ctx, req)
			if se := errors.FromError(err); se != nil {
				code = se.Code
				reason = se.Reason
			}
			level, stack := extractError(err)
			attrs := []slog.Attr{
				slog.String("kind", "client"),
				slog.String("component", kind),
				slog.String("operation", operation),
				slog.String("args", extractArgs(req)),
				slog.Int64("code", int64(code)),
				slog.String("reason", reason),
				slog.Float64("latency", time.Since(startTime).Seconds()),
			}
			if err != nil {
				attrs = append(attrs, slog.Any("error", err))
				if stack != "" {
					attrs = append(attrs, slog.String("stack", stack))
				}
			}
			logger.LogAttrs(ctx, level, "client request", attrs...)
			return
		}
	}
}

// extractArgs returns the string of the req
func extractArgs(req any) string {
	if redacter, ok := req.(Redacter); ok {
		return redacter.Redact()
	}
	if stringer, ok := req.(fmt.Stringer); ok {
		return stringer.String()
	}
	return fmt.Sprintf("%+v", req)
}

// extractError returns the log level and an optional stack string for err.
// stack is empty when %+v adds nothing beyond Error(), so callers do not
// duplicate the same text under both "error" and "stack".
func extractError(err error) (slog.Level, string) {
	if err == nil {
		return slog.LevelInfo, ""
	}
	level := slog.LevelError
	if se := errors.FromError(err); se != nil {
		level = levelFromCode(se.Code)
	}
	stack := fmt.Sprintf("%+v", err)
	if stack == err.Error() {
		stack = ""
	}
	return level, stack
}

// levelFromCode maps an HTTP-style status code to a slog level.
// 4xx client/business errors and 503 Service Unavailable (overload or
// shedding, usually transient) are warnings so they do not page like
// other 5xx errors.
func levelFromCode(code int32) slog.Level {
	switch {
	case code >= 100 && code < 400:
		return slog.LevelInfo
	case code >= 400 && code < 500, code == 503:
		return slog.LevelWarn
	default:
		return slog.LevelError
	}
}
