package http

import (
	"context"
	"net/http"

	"github.com/go-kratos/kratos/v2/middleware"
)

// ServerOption is HTTP server option.
type ServerOption func(*serverOptions)

// serverOptions is HTTP server options.
type serverOptions struct {
	requestDecoder  DecodeRequestFunc
	responseEncoder EncodeResponseFunc
	errorEncoder    EncodeErrorFunc
	middleware      middleware.Middleware
}

// DecodeRequestFunc is decode request func.
type DecodeRequestFunc func(ctx context.Context, in interface{}, req *http.Request) error

// EncodeResponseFunc is encode response func.
type EncodeResponseFunc func(ctx context.Context, out interface{}, res http.ResponseWriter, req *http.Request) error

// EncodeErrorFunc is encode error func.
type EncodeErrorFunc func(ctx context.Context, err error, res http.ResponseWriter, req *http.Request)

// ServerDecodeRequestFunc with decode request option.
func ServerDecodeRequestFunc(fn EncodeErrorFunc) ServerOption {
	return func(o *serverOptions) {
		o.errorEncoder = fn
	}
}

// ServerEncodeResponseFunc with response handler option.
func ServerEncodeResponseFunc(fn EncodeResponseFunc) ServerOption {
	return func(o *serverOptions) {
		o.responseEncoder = fn
	}
}

// ServerEncodeErrorFunc with error handler option.
func ServerEncodeErrorFunc(fn EncodeErrorFunc) ServerOption {
	return func(o *serverOptions) {
		o.errorEncoder = fn
	}
}

// ServerMiddleware with server middleware option.
func ServerMiddleware(m ...middleware.Middleware) ServerOption {
	return func(o *serverOptions) {
		o.middleware = middleware.Chain(m[0], m[1:]...)
	}
}
