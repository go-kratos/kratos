package ratelimit

import (
	"context"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/middleware"
)

type Option func(*options)

// WithErrorCode set error code when limiter triggered,
// default error code 503
func WithErrorCode(code int) Option {
	return func(o *options) {
		o.errCode = code
	}
}

func WithErrorReason(reason string) Option {
	return func(o *options) {
		o.errReason = reason
	}
}

func WithErrorMessage(message string) Option {
	return func(o *options) {
		o.errMessage = message
	}
}

type options struct {
	errCode    int
	errReason  string
	errMessage string
}

// Limiter limit interface.
type Limiter interface {
	Allow(ctx context.Context) (done func(), err error)
}

// RateLimiter middleware
func RateLimiter(limiter Limiter, opts ...Option) middleware.Middleware {
	options := &options{
		errCode:    503,
		errReason:  "rate limit exceeded",
		errMessage: "service unavailable due to rate limit exceeded",
	}
	for _, o := range opts {
		o(options)
	}

	// errLimitExceed is returned when the rate limiter is
	// triggered and the request is rejected due to limit exceeded.
	errLimitExceed := errors.New(options.errCode, options.errReason, options.errMessage)
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
			done, e := limiter.Allow(ctx)
			if e != nil {
				// blocked
				return nil, errLimitExceed
			}
			// passed
			reply, err = handler(ctx, req)
			done()
			return
		}
	}
}
