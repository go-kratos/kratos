package http

import (
	"context"
	"net/http"

	"github.com/go-kratos/kratos/v2/middleware"
)

// ServerOption is HTTP server option.
type ServerOption func(o *serverOptions)

// serverOptions is HTTP server options.
type serverOptions struct {
	errorHandler ErrorHandler
	middleware   middleware.Middleware
	certFile     string
	keyFile      string
}

// ErrorHandler is encoding an error to the ResponseWriter.
type ErrorHandler func(ctx context.Context, err error, m Marshaler, w http.ResponseWriter)

// WithErrorHandler with error handler option.
func WithErrorHandler(h ErrorHandler) ServerOption {
	return func(o *serverOptions) {
		o.errorHandler = h
	}
}

func ServerMiddleware(m middleware.Middleware) ServerOption {
	return func(o *serverOptions) { o.middleware = m }
}

func ServerTLSCert(certFile, keyFile string) ServerOption {
	return func(o *serverOptions) {
		o.certFile = certFile
		o.keyFile = keyFile
	}
}
