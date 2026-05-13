package ratelimit

import (
	"context"

	"github.com/go-kratos/kratos/v3/errors"
	internalratelimit "github.com/go-kratos/kratos/v3/internal/ratelimit"
	"github.com/go-kratos/kratos/v3/middleware"
)

// ErrLimitExceed is service unavailable due to rate limit exceeded.
var ErrLimitExceed = errors.New(429, "RATELIMIT", "service unavailable due to rate limit exceeded")

// DoneFunc records request completion.
type DoneFunc = internalratelimit.DoneFunc

// DoneInfo contains request completion metadata.
type DoneInfo = internalratelimit.DoneInfo

// Limiter is a rate limiter.
type Limiter = internalratelimit.Limiter

// Option is ratelimit option.
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

// Server ratelimiter middleware
func Server(opts ...Option) middleware.Middleware {
	options := &options{
		limiter: internalratelimit.NewLimiter(),
	}
	for _, o := range opts {
		o(options)
	}
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req any) (reply any, err error) {
			done, e := options.limiter.Allow()
			if e != nil {
				// rejected
				return nil, ErrLimitExceed
			}
			// allowed
			reply, err = handler(ctx, req)
			done(DoneInfo{Err: err})
			return
		}
	}
}
