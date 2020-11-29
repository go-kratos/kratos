package middleware

import "context"

// Endpoint is the server endpoint.
type Endpoint func(ctx context.Context, request interface{}) (response interface{}, err error)

// Middleware is transport middleware.
type Middleware func(Endpoint) Endpoint

// Chain is the middleare function for the given pattern.
func Chain(outer Middleware, others ...Middleware) Middleware {
	return func(next Endpoint) Endpoint {
		for i := len(others) - 1; i >= 0; i-- {
			next = others[i](next)
		}
		return outer(next)
	}
}
