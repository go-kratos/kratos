package breaker

import (
	"context"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/middleware"
)

type Option func(*options)

// WithErrorCode set error code when breaker is triggered,
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

// CircuitBreaker middleware will return errBreakerTriggered when the circuit
// breaker is triggered and the request is rejected directly.
func CircuitBreaker(breaker Breaker, opts ...Option) middleware.Middleware {
	options := &options{
		errCode:    503,
		errReason:  "circuit breaker triggered",
		errMessage: "request failed due to circuit breaker triggered",
	}
	for _, o := range opts {
		o(options)
	}

	errBreakerTriggered := errors.New(options.errCode, options.errReason, options.errMessage)
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			if err := breaker.Allow(ctx); err != nil {
				// reject request
				return nil, errBreakerTriggered
			}
			// allow request
			reply, err := handler(ctx, req)
			breaker.Mark(breaker.Check(err))
			return reply, err
		}
	}
}
