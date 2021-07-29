package ratelimit

import (
	"context"
	"errors"

	"github.com/go-kratos/kratos/v2/middleware"
)

// ErrLimitExceed is returned when the rate limiter is
// triggered and the request is rejected due to limit exceeded.
var ErrLimitExceed = errors.New("rate limit exceeded")

// Limiter limit interface.
type Limiter interface {
	Allow(ctx context.Context) (done func(), err error)
}

// RateLimiter middleware
func RateLimiter(limiter Limiter) middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
			done, e := limiter.Allow(ctx)
			if e != nil {
				// blocked
				return nil, ErrLimitExceed
			}
			// passed
			reply, err = handler(ctx, req)
			done()
			return
		}
	}
}
