package http

import (
	"context"
	"net/http"
)

// ServerOption is HTTP server option.
type ServerOption func(o *ServerOptions)

// ServerOptions is HTTP server options.
type ServerOptions struct {
	ErrorHandler ErrorHandler
}

// ErrorHandler is encoding an error to the ResponseWriter.
type ErrorHandler func(ctx context.Context, err error, m Marshaler, w http.ResponseWriter)

// WithErrorHandler with error handler option.
func WithErrorHandler(h ErrorHandler) ServerOption {
	return func(o *ServerOptions) {
		o.ErrorHandler = h
	}
}
