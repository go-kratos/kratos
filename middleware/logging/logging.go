package logging

import (
	"context"
	"fmt"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
)

type Option func(*options)

type options struct {
	printArgs bool
}

// WithPrintArgs print args
func WithPrintArgs() Option {
	return func(opts *options) {
		opts.printArgs = true
	}
}

// Server is an server logging middleware.
func Server(logger log.Logger, opts ...Option) middleware.Middleware {
	options := options{}
	for _, o := range opts {
		o(&options)
	}

	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
			reply, err = handler(ctx, req)
			if tr, ok := transport.FromContext(ctx); ok {
				switch tr.Kind {
				case transport.KindHTTP:
					httpServerLog(logger, ctx, extractArgs(req, options.printArgs), err)
				case transport.KindGRPC:
					grpcServerLog(logger, ctx, extractArgs(req, options.printArgs), err)
				}
			}
			return
		}
	}
}

// Client is an client logging middleware.
func Client(logger log.Logger, opts ...Option) middleware.Middleware {
	options := options{}
	for _, o := range opts {
		o(&options)
	}
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
			reply, err = handler(ctx, req)
			if tr, ok := transport.FromContext(ctx); ok {
				switch tr.Kind {
				case transport.KindHTTP:
					httpClientLog(logger, ctx, extractArgs(req, options.printArgs), err)
				case transport.KindGRPC:
					grpcClientLog(logger, ctx, extractArgs(req, options.printArgs), err)
				}
			}
			return
		}
	}
}

func extractArgs(req interface{}, printArgs bool) string {
	if printArgs {
		if stringer, ok := req.(fmt.Stringer); ok {
			return stringer.String()
		}
		return fmt.Sprintf("%+v", req)
	}
	return ""
}

func extractError(err error) (code int, errMsg string) {
	if err != nil {
		code = errors.Code(err)
		errMsg = fmt.Sprintf("%+v", err)
	}
	return
}
