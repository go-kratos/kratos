package circuitbreaker

import (
	"context"

	"github.com/go-kratos/aegis/circuitbreaker"
	"github.com/go-kratos/aegis/circuitbreaker/sre"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/internal/group"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
)

// ErrNotAllowed is request failed due to circuit breaker triggered.
var ErrNotAllowed = errors.New(503, "CIRCUITBREAKER", "request failed due to circuit breaker triggered")

// Option is circuit breaker option.
type Option func(*options)

// WithGroup with circuit breaker group.
// NOTE: implements generics circuitbreaker.CircuitBreaker
func WithGroup(g *group.Group) Option {
	return func(o *options) {
		o.group = g
	}
}

// WithCircuitBreaker with circuit breaker genFunc.
func WithCircuitBreaker(genBreakerFunc func() circuitbreaker.CircuitBreaker) Option {
	return func(o *options) {
		o.group = group.NewGroup(func() any {
			return genBreakerFunc()
		})
	}
}

type options struct {
	group *group.Group
}

// Client circuitbreaker middleware will return errBreakerTriggered when the circuit
// breaker is triggered and the request is rejected directly.
func Client(opts ...Option) middleware.Middleware {
	opt := &options{
		group: group.NewGroup(func() any {
			return sre.NewBreaker()
		}),
	}
	for _, o := range opts {
		o(opt)
	}
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req any) (any, error) {
			info, _ := transport.FromClientContext(ctx)
			breaker := opt.group.Get(info.Operation()).(circuitbreaker.CircuitBreaker)
			if err := breaker.Allow(); err != nil {
				// rejected
				// NOTE: when client reject requests locally,
				// continue to add counter let the drop ratio higher.
				breaker.MarkFailed()
				return nil, ErrNotAllowed
			}
			// allowed
			reply, err := handler(ctx, req)
			if err != nil && (errors.IsInternalServer(err) || errors.IsServiceUnavailable(err) || errors.IsGatewayTimeout(err)) {
				breaker.MarkFailed()
			} else {
				breaker.MarkSuccess()
			}
			return reply, err
		}
	}
}
