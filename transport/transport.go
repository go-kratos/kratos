package transport

import (
	"context"
	"net/url"

	// init encoding
	_ "github.com/go-kratos/kratos/v2/encoding/json"
	_ "github.com/go-kratos/kratos/v2/encoding/proto"
	_ "github.com/go-kratos/kratos/v2/encoding/xml"
	_ "github.com/go-kratos/kratos/v2/encoding/yaml"
	"github.com/go-kratos/kratos/v2/metadata"
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

// Transporter is transport context value interface.
type Transporter interface {
	Kind() string
	Endpoint() string
	// Clone returns a deep copy of Transporter
	Clone() Transporter

	ServiceMethod() string
	SetServiceMethod(string)

	Metadata() metadata.Metadata
	// WithMetadata merge new metadata into transport,
	// it will override old metadata key value if key exists
	WithMetadata(metadata.Metadata)
}

// Defines a set of transport kind
const (
	KindGRPC = "grpc"
	KindHTTP = "http"
)

type serverTransportKey struct{}
type clientTransportKey struct{}

// NewServerContext returns a new Context that carries value.
func NewServerContext(ctx context.Context, tr Transporter) context.Context {
	return context.WithValue(ctx, serverTransportKey{}, tr)
}

// FromServerContext returns the Transport value stored in ctx, if any.
func FromServerContext(ctx context.Context) (tr Transporter, ok bool) {
	tr, ok = ctx.Value(serverTransportKey{}).(Transporter)
	return
}

// NewClientContext returns a new Context that carries value.
func NewClientContext(ctx context.Context, tr Transporter) context.Context {
	return context.WithValue(ctx, clientTransportKey{}, tr)
}

// FromClientContext returns the Transport value stored in ctx, if any.
func FromClientContext(ctx context.Context) (tr Transporter, ok bool) {
	tr, ok = ctx.Value(clientTransportKey{}).(Transporter)
	return
}

// SetServerServiceMethod set serviceMethod into context transport
func SetServerServiceMethod(ctx context.Context, serviceMethod string) {
	tr, ok := FromServerContext(ctx)
	if ok {
		tr.SetServiceMethod(serviceMethod)
	}
}
