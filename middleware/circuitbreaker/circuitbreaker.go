package circuitbreaker

import (
	"context"

	"github.com/go-kratos/aegis/circuitbreaker"
	"github.com/go-kratos/aegis/circuitbreaker/sre"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/middleware"
)

// Option is circuit breaker option.
type Option func(*options)

// WithBreaker set circuit breaker implentation
func WithBreaker(breaker circuitbreaker.CircuitBreaker) Option {
	return func(o *options) {
		o.breaker = breaker
	}
}

type options struct {
	breaker circuitbreaker.CircuitBreaker
}

// Client circuitbreaker middleware will return errBreakerTriggered when the circuit
// breaker is triggered and the request is rejected directly.
func Client(opts ...Option) middleware.Middleware {
	options := &options{
		breaker: sre.NewBreaker(),
	}
	for _, o := range opts {
		o(options)
	}
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			if err := options.breaker.Allow(); err != nil {
				// rejected
				// NOTE: when client reject requets locally,
				// continue add counter let the drop ratio higher.
				options.breaker.MarkFailed()
				return nil, errors.New(503, "CIRCUITBREAKER", "request failed due to circuit breaker triggered")
			}
			// allowed
			reply, err := handler(ctx, req)
			if err != nil {
				options.breaker.MarkFailed()
			} else {
				options.breaker.MarkSuccess()
			}
			return reply, err
		}
	}
}
