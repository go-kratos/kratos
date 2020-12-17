package http

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-kratos/kratos/v2/middleware"
)

// ServerOption is HTTP server option.
type ServerOption func(o *serverOptions)

// serverOptions is HTTP server options.
type serverOptions struct {
	errorHandler    ErrorHandler
	responseHandler ResponseHandler
	middleware      middleware.Middleware
}

// ErrorHandler is encoding an error to the ResponseWriter.
type ErrorHandler func(ctx context.Context, err error, m Marshaler, w http.ResponseWriter)

// ResponseHandler is encoding an data to the ResponseWriter.
type ResponseHandler func(ctx context.Context, out interface{}, m Marshaler, w http.ResponseWriter)

// ServerErrorHandler with error handler option.
func ServerErrorHandler(h ErrorHandler) ServerOption {
	return func(o *serverOptions) {
		o.errorHandler = h
	}
}

// ServerResponseHandler with error handler option.
func ServerResponseHandler(h ResponseHandler) ServerOption {
	return func(o *serverOptions) {
		o.responseHandler = h
	}
}

// ServerMiddleware with server middleware option.
func ServerMiddleware(m ...middleware.Middleware) ServerOption {
	return func(o *serverOptions) {
		o.middleware = middleware.Chain(m[0], m[1:]...)
	}
}

// DefaultErrorHandler is default errors handler.
func DefaultErrorHandler(ctx context.Context, err error, m Marshaler, w http.ResponseWriter) {
	se := StatusError(err)
	w.WriteHeader(se.Code)
	if m != nil {
		b, _ := m.Marshal(se)
		w.Write(b)
	} else {
		b, _ := json.Marshal(se)
		w.Write(b)
	}
}

// DefaultResponseHandler is default response handler.
func DefaultResponseHandler(ctx context.Context, out interface{}, m Marshaler, w http.ResponseWriter) {
	data, err := m.Marshal(out)
	if err != nil {
		DefaultErrorHandler(ctx, err, m, w)
		return
	}
	_, err = w.Write(data)
	return
}
