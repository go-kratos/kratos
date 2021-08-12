package breaker

import (
	"context"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/middleware"
)

// Breaker interface defines a circuit breaker for Handler
type Breaker interface {
	// Allow check whether current request is allowed to execute.
	// it will return a not nil error when breaker is on.
	Allow(ctx context.Context) error
	// Check breaker use Check(error) to check whether
	// the request succeeded of failed,
	Check(err error) bool
	// Mark whether the current request is success
	Mark(isSuccess bool)
}

type Option func(*options)

// WithBreaker set circuit breaker implentation
func WithBreaker(breaker Breaker) Option {
	return func(o *options) {
		o.breaker = breaker
	}
}

type options struct {
	breaker Breaker
}

// CircuitBreaker middleware will return errBreakerTriggered when the circuit
// breaker is triggered and the request is rejected directly.
func CircuitBreaker(opts ...Option) middleware.Middleware {
	options := &options{}
	for _, o := range opts {
		o(options)
	}

	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			if err := options.breaker.Allow(ctx); err != nil {
				// rejected
				return nil, errors.New(503, "BREAKER", "request failed due to circuit breaker triggered")
			}
			// allowed
			reply, err := handler(ctx, req)
			options.breaker.Mark(options.breaker.Check(err))
			return reply, err
		}
	}
}
