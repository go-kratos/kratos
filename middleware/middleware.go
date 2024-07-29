package middleware

import (
	"context"
	"google.golang.org/grpc"
)

// Handler defines the handler invoked by Middleware.
type Handler func(ctx context.Context, req interface{}) (interface{}, error)

// StreamHandler defines the handler invoked by Middleware for stream calls.
type StreamHandler func(srv interface{}, stream grpc.ServerStream) error

// Middleware is HTTP/gRPC transport middleware.
type Middleware func(Handler) Handler

// StreamMiddleware is gRPC stream transport middleware.
type StreamMiddleware func(StreamHandler) StreamHandler

// Chain returns a Middleware that specifies the chained handler for endpoint.
func Chain(m ...Middleware) Middleware {
	return func(next Handler) Handler {
		for i := len(m) - 1; i >= 0; i-- {
			next = m[i](next)
		}
		return next
	}
}

// ChainStream returns a StreamMiddleware that specifies the chained handler for endpoint.
func ChainStream(m ...StreamMiddleware) StreamMiddleware {
	return func(next StreamHandler) StreamHandler {
		for i := len(m) - 1; i >= 0; i-- {
			next = m[i](next)
		}
		return next
	}
}
