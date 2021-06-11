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

	Method() string
	SetMethod(string)

	Metadata() metadata.Metadata
	// WithMetadata merge new metadata into transport,
	// it will override old metadata key value if key exists
	WithMetadata(metadata.Metadata)
}

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

// SetServerMethod set serviceMethod into context transport.
func SetServerMethod(ctx context.Context, method string) {
	if tr, ok := FromServerContext(ctx); ok {
		tr.SetMethod(method)
	}
}

// SetClientMethod set serviceMethod into context transport.
func SetClientMethod(ctx context.Context, method string) {
	if tr, ok := FromClientContext(ctx); ok {
		tr.SetMethod(method)
	}
}

// Metadata returns incoming metadata from server transport.
func Metadata(ctx context.Context) metadata.Metadata {
	if tr, ok := FromServerContext(ctx); ok {
		return tr.Metadata()
	}
	return metadata.Metadata{}
}

// SetMetadata sets outgoing metadata into client transport.
func SetMetadata(ctx context.Context, md metadata.Metadata) {
	if tr, ok := FromClientContext(ctx); ok {
		tr.WithMetadata(md)
	}
}
