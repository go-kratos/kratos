package http

import (
	"context"

	"github.com/go-kratos/kratos/v2/metadata"
	"github.com/go-kratos/kratos/v2/transport"
)

var (
	_ transport.Transporter = &Transport{}
)

// Transport is an HTTP transport.
type Transport struct {
	endpoint  string
	path      string
	method    string
	operation string
	metadata  metadata.Metadata
}

// Kind returns the transport kind.
func (tr *Transport) Kind() string {
	return "http"
}

// Endpoint returns the transport endpoint.
func (tr *Transport) Endpoint() string {
	return tr.endpoint
}

// Operation returns the transport operation.
func (tr *Transport) Operation() string {
	return tr.operation
}

// SetOperation sets the transport operation.
func (tr *Transport) SetOperation(operation string) {
	tr.operation = operation
}

// Metadata returns the transport metadata.
func (tr *Transport) Metadata() metadata.Metadata {
	return tr.metadata
}

// WithMetadata with a metadata into transport md.
func (tr *Transport) WithMetadata(md metadata.Metadata) {
	for k, v := range md {
		tr.metadata.Set(k, v)
	}
}

// Path returns the Transport path from server context.
func Path(ctx context.Context) string {
	if tr, ok := transport.FromServerContext(ctx); ok {
		if tr, ok := tr.(*Transport); ok {
			return tr.path
		}
	}
	return ""
}

// Method returns the Transport method from server context.
func Method(ctx context.Context) string {
	if tr, ok := transport.FromServerContext(ctx); ok {
		if tr, ok := tr.(*Transport); ok {
			return tr.method
		}
	}
	return ""
}
