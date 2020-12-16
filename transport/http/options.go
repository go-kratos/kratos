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
	errorHandler    errorHandler
	responseHandler responseHandler
	middleware      middleware.Middleware
}

// ErrorHandler is encoding an error to the ResponseWriter.
type errorHandler func(ctx context.Context, err error, m Marshaler, w http.ResponseWriter)

// ResponseHandler is encoding an data to the ResponseWriter.
type responseHandler func(ctx context.Context, out interface{}, m Marshaler, w http.ResponseWriter) error

// ErrorHandler with error handler option.
func ErrorHandler(h errorHandler) ServerOption {
	return func(o *serverOptions) {
		o.errorHandler = h
	}
}

// ResponseHandler with error handler option.
func ResponseHandler(h responseHandler) ServerOption {
	return func(o *serverOptions) {
		o.responseHandler = h
	}
}

// ServerMiddleware .
func ServerMiddleware(m middleware.Middleware) ServerOption {
	return func(o *serverOptions) { o.middleware = m }
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
func DefaultResponseHandler(ctx context.Context, out interface{}, m Marshaler, w http.ResponseWriter) error {
	data, err := m.Marshal(out)
	if err != nil {
		return err
	}
	_, err = w.Write(data)
	return err
}
