package sentry

import (
	"context"
	"fmt"
	"net"
	"os"
	"strings"
	"time"

	"github.com/getsentry/sentry-go"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
)

type ctxKey struct{}

type Option func(*options)

type options struct {
	repanic         bool
	waitForDelivery bool
	timeout         time.Duration
	tags            map[string]interface{}
}

// WithRepanic repanic configures whether Sentry should repanic after recovery, in most cases it should be set to true.
func WithRepanic(repanic bool) Option {
	return func(opts *options) {
		opts.repanic = repanic
	}
}

// WithWaitForDelivery waitForDelivery configures whether you want to block the request before moving forward with the response.
func WithWaitForDelivery(waitForDelivery bool) Option {
	return func(opts *options) {
		opts.waitForDelivery = waitForDelivery
	}
}

// WithTimeout timeout for the event delivery requests.
func WithTimeout(timeout time.Duration) Option {
	return func(opts *options) {
		opts.timeout = timeout
	}
}

// WithTags global tags injection, the value type must be string or log.Valuer
func WithTags(kvs map[string]interface{}) Option {
	return func(opts *options) {
		opts.tags = kvs
	}
}

// Server returns a new server middleware for Sentry.
func Server(opts ...Option) middleware.Middleware {
	conf := options{repanic: true}
	for _, o := range opts {
		o(&conf)
	}
	if conf.timeout == 0 {
		conf.timeout = 2 * time.Second
	}
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
			hub := GetHubFromContext(ctx)
			scope := hub.Scope()

			for k, v := range conf.tags {
				switch val := v.(type) {
				case string:
					scope.SetTag(k, val)
				case log.Valuer:
					scope.SetTag(k, fmt.Sprintf("%v", val(ctx)))
				}
			}

			if tr, ok := transport.FromServerContext(ctx); ok {
				switch tr.Kind() {
				case transport.KindGRPC:
					gtr := tr.(*grpc.Transport)
					scope.SetContext("gRPC", map[string]interface{}{
						"endpoint":  gtr.Endpoint(),
						"operation": gtr.Operation(),
					})
					headers := make(map[string]interface{})
					for _, k := range gtr.RequestHeader().Keys() {
						headers[k] = gtr.RequestHeader().Get(k)
					}
					scope.SetContext("Headers", headers)
				case transport.KindHTTP:
					htr := tr.(*http.Transport)
					r := htr.Request()
					scope.SetRequest(r)
				}
			}

			ctx = context.WithValue(ctx, ctxKey{}, hub)
			defer recoverWithSentry(ctx, conf, hub, req)
			return handler(ctx, req)
		}
	}
}

func recoverWithSentry(ctx context.Context, opts options, hub *sentry.Hub, req interface{}) {
	if err := recover(); err != nil {
		if !isBrokenPipeError(err) {
			eventID := hub.RecoverWithContext(
				context.WithValue(ctx, sentry.RequestContextKey, req),
				err,
			)
			if eventID != nil && opts.waitForDelivery {
				hub.Flush(opts.timeout)
			}
		}
		if opts.repanic {
			panic(err)
		}
	}
}

func isBrokenPipeError(err interface{}) bool {
	if netErr, ok := err.(*net.OpError); ok {
		if sysErr, ok := netErr.Err.(*os.SyscallError); ok {
			if strings.Contains(strings.ToLower(sysErr.Error()), "broken pipe") ||
				strings.Contains(strings.ToLower(sysErr.Error()), "connection reset by peer") {
				return true
			}
		}
	}
	return false
}

// GetHubFromContext retrieves attached *sentry.Hub instance from context or sentry.
// You can use this hub for extra information reporting
func GetHubFromContext(ctx context.Context) *sentry.Hub {
	if hub, ok := ctx.Value(ctxKey{}).(*sentry.Hub); ok {
		return hub
	}
	return sentry.CurrentHub().Clone()
}
