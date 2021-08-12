package ratelimit

import (
	"context"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/sra/ratelimit/bbr"
)

// Limiter limit interface.
type Limiter interface {
	Allow(ctx context.Context) (done func(), err error)
}

type Option func(*options)

// WithLimiter set Limiter implementation,
// default is bbr limiter
func WithLimiter(limiter Limiter) Option {
	return func(o *options) {
		o.limiter = limiter
	}
}

type options struct {
	limiter Limiter
}

// RateLimiter middleware
func RateLimiter(limiter Limiter, opts ...Option) middleware.Middleware {
	options := &options{
		limiter: bbr.NewLimiter(),
	}
	for _, o := range opts {
		o(options)
	}
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
			done, e := limiter.Allow(ctx)
			if e != nil {
				// rejected
				return nil, errors.New(429, "RATELIMIT", "service unavailable due to rate limit exceeded")
			}
			// allowed
			reply, err = handler(ctx, req)
			done()
			return
		}
	}
}
