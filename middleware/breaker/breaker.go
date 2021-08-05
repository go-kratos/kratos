package breaker

import (
	"context"
	"errors"

	"github.com/go-kratos/kratos/v2/middleware"
)

// ErrBreakerTriggered is returned when the
// circuit breaker is triggered and the request
// is rejected directly.
var ErrBreakerTriggered = errors.New("circuit breaker triggered")

// Breaker interface defines a circuit breaker for Handler
type Breaker interface {
	// Allow check whether current request is allowed to execute.
	// it will return a not nil error when breaker is on.
	Allow(ctx context.Context) error
	// Check breaker use Check(error) to check whether
	// the request succeeded of failed,
	Check(err error) bool
	// Mark wheather the current request is success
	Mark(isSuccess bool)
}

func CircuitBreaker(breaker Breaker) middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			if err := breaker.Allow(ctx); err != nil {
				// reject request
				return nil, ErrBreakerTriggered
			}
			// allow request
			reply, err := handler(ctx, req)
			breaker.Mark(breaker.Check(err))
			return reply, err
		}
	}
}
