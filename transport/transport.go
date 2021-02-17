package transport

import (
	"context"

	// init encoding
	_ "github.com/go-kratos/kratos/v2/encoding/json"
	_ "github.com/go-kratos/kratos/v2/encoding/proto"
)

// Server is transport server.
type Server interface {
	Endpoint() (string, error)
	Start() error
	Stop() error
}

// Transport is transport context value.
type Transport struct {
	Kind string
}

type transportKey struct{}

// NewContext returns a new Context that carries value.
func NewContext(ctx context.Context, tr Transport) context.Context {
	return context.WithValue(ctx, transportKey{}, tr)
}

// FromContext returns the Transport value stored in ctx, if any.
func FromContext(ctx context.Context) (tr Transport, ok bool) {
	tr, ok = ctx.Value(transportKey{}).(Transport)
	return
}
