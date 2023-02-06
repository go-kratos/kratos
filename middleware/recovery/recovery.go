package recovery

import (
	"context"
	"fmt"
	"runtime"
	"time"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
)

// Redacter defines how to log an object
type Redacter interface {
	Redact() string
}

// ErrUnknownRequest is unknown request error.
var ErrUnknownRequest = errors.InternalServer("UNKNOWN", "unknown request error")

// HandlerFunc is recovery handler func.
type HandlerFunc func(ctx context.Context, req, err interface{}) error

// Option is recovery option.
type Option func(*options)

type options struct {
	handler HandlerFunc
}

// WithHandler with recovery handler.
func WithHandler(h HandlerFunc) Option {
	return func(o *options) {
		o.handler = h
	}
}

// Recovery is a server middleware that recovers from any panics.
func Recovery(opts ...Option) middleware.Middleware {
	op := options{
		handler: func(ctx context.Context, req, err interface{}) error {
			return ErrUnknownRequest
		},
	}
	for _, o := range opts {
		o(&op)
	}
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
			startTime := time.Now()
			defer func() {
				if rerr := recover(); rerr != nil {
					var (
						component string
						operation string
					)
					if info, ok := transport.FromServerContext(ctx); ok {
						component = info.Kind().String()
						operation = info.Operation()
					}
					if info, ok := transport.FromClientContext(ctx); ok {
						component = info.Kind().String()
						operation = info.Operation()
					}
					buf := make([]byte, 64<<10) //nolint:gomnd
					n := runtime.Stack(buf, false)
					buf = buf[:n]
					stack := fmt.Sprintf("%v: %+v\n%s\n", rerr, req, buf)
					log.Context(ctx).Log(log.LevelError,
						"component", component,
						"operation", operation,
						"args", extractArgs(req),
						"code", 500,
						"reason", "UNKNOWN",
						"stack", stack,
						"latency", time.Since(startTime).Seconds(),
					)
				}
			}()
			return handler(ctx, req)
		}
	}
}

// extractArgs returns the string of the req
func extractArgs(req interface{}) string {
	if redacter, ok := req.(Redacter); ok {
		return redacter.Redact()
	}
	if stringer, ok := req.(fmt.Stringer); ok {
		return stringer.String()
	}
	return fmt.Sprintf("%+v", req)
}
