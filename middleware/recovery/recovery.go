package recovery

import (
	"context"
	"fmt"
	"runtime"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware"
)

// HandlerFunc is recovery handler func.
type HandlerFunc func(ctx context.Context, req, err interface{}) error

// Option is recovery option.
type Option func(*options)

type options struct {
	handler HandlerFunc
	logger  log.Logger
}

// WithHandler with recovery handler.
func WithHandler(h HandlerFunc) Option {
	return func(o *options) {
		o.handler = h
	}
}

// WithLogger with recovery logger.
func WithLogger(logger log.Logger) Option {
	return func(o *options) {
		o.logger = logger
	}
}

// Recovery is a server middleware that recovers from any panics.
func Recovery(opts ...Option) middleware.Middleware {
	options := options{
		logger: log.DefaultLogger,
		handler: func(ctx context.Context, req, err interface{}) error {
			return errors.InternalServer("RECOVERY", fmt.Sprintf("panic triggered: %v", err))
		},
	}
	for _, o := range opts {
		o(&options)
	}
	logger := log.NewHelper(options.logger)
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
			defer func() {
				if rerr := recover(); rerr != nil {
					buf := make([]byte, 64<<10)
					n := runtime.Stack(buf, false)
					buf = buf[:n]
					logger.Errorf("%v: %+v\n%s\n", rerr, req, buf)

					err = options.handler(ctx, req, rerr)
				}
			}()
			return handler(ctx, req)
		}
	}
}
