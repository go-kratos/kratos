package grpc

import (
	"github.com/go-kratos/kratos/v2/transport"
)

var (
	_ transport.Transporter = &Transport{}
)

// Transport is a gRPC transport.
type Transport struct {
	endpoint  string
	operation string
	header    transport.Header
}

// Kind returns the transport kind.
func (tr *Transport) Kind() string {
	return "grpc"
}

// Endpoint returns the transport endpoint.
func (tr *Transport) Endpoint() string {
	return tr.endpoint
}

// Operation returns the transport operation.
func (tr *Transport) Operation() string {
	return tr.operation
}

// Header returns the transport header.
func (tr *Transport) Header() transport.Header {
	return tr.header
}
