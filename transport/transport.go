package transport

import (
	"context"
	"net/url"

	// init encoding
	_ "github.com/go-kratos/kratos/v2/encoding/json"
	_ "github.com/go-kratos/kratos/v2/encoding/proto"
	_ "github.com/go-kratos/kratos/v2/encoding/xml"
	_ "github.com/go-kratos/kratos/v2/encoding/yaml"
)

// Server is transport server.
type Server interface {
	Start(context.Context) error
	Stop(context.Context) error
}

// Endpointer is registry endpoint.
type Endpointer interface {
	Endpoint() (*url.URL, error)
}

// Transport is transport context value.
type Transport struct {
	Kind     Kind
	Endpoint string
}

// Kind defines the type of Transport
type Kind string

// Defines a set of transport kind
const (
	KindGRPC Kind = "gRPC"
	KindHTTP Kind = "HTTP"
)

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
