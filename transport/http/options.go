package http

import (
	"context"
	"net/http"

	"github.com/go-kratos/kratos/v2/encoding"
)

// ServerOption is HTTP server option.
type ServerOption func(o *ServerOptions)

// ServerOptions is HTTP server options.
type ServerOptions struct {
	Address      string
	ErrorHandler ErrorHandler
}

// ErrorHandler is encoding an error to the ResponseWriter.
type ErrorHandler func(ctx context.Context, err error, codec encoding.Codec, w http.ResponseWriter)

// WithAddress is bind address option.
func WithAddress(a string) ServerOption {
	return func(o *ServerOptions) {
		o.Address = a
	}
}
