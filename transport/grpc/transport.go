package grpc

import (
	"github.com/go-kratos/kratos/v2/metadata"
	"github.com/go-kratos/kratos/v2/transport"
)

var (
	_ transport.Transporter = &Transport{}
)

// Transport is a gRPC transport.
type Transport struct {
	endpoint string
	method   string
	metadata metadata.Metadata
}

// Kind returns the transport kind.
func (tr *Transport) Kind() string {
	return "grpc"
}

// Endpoint returns the transport endpoint.
func (tr *Transport) Endpoint() string {
	return tr.endpoint
}

// Method returns the transport method.
func (tr *Transport) Method() string {
	return tr.method
}

// SetMethod sets the transport method.
func (tr *Transport) SetMethod(method string) {
	tr.method = method
}

// Metadata returns the transport metadata.
func (tr *Transport) Metadata() metadata.Metadata {
	return tr.metadata
}

// WithMetadata with a metadata into transport md.
func (tr *Transport) WithMetadata(md metadata.Metadata) {
	if tr.metadata == nil {
		tr.metadata = md
		return
	}
	for k, v := range md {
		tr.metadata.Set(k, v)
	}
}
