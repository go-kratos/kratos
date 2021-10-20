package transport

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

// mockTransport is a gRPC transport.
type mockTransport struct {
	endpoint  string
	operation string
}

// Kind returns the transport kind.
func (tr *mockTransport) Kind() Kind {
	return KindGRPC
}

// Endpoint returns the transport endpoint.
func (tr *mockTransport) Endpoint() string {
	return tr.endpoint
}

// Operation returns the transport operation.
func (tr *mockTransport) Operation() string {
	return tr.operation
}

// RequestHeader returns the request header.
func (tr *mockTransport) RequestHeader() Header {
	return nil
}

// ReplyHeader returns the reply header.
func (tr *mockTransport) ReplyHeader() Header {
	return nil
}

func TestServerTransport(t *testing.T) {
	ctx := context.Background()

	ctx = NewServerContext(ctx, &mockTransport{endpoint: "test_endpoint"})
	tr, ok := FromServerContext(ctx)

	assert.Equal(t, true, ok)
	assert.NotNil(t, tr)
	mtr, ok := tr.(*mockTransport)
	assert.Equal(t, true, ok)
	assert.NotNil(t, mtr)
	assert.Equal(t, mtr.endpoint, "test_endpoint")
}

func TestClientTransport(t *testing.T) {
	ctx := context.Background()

	ctx = NewClientContext(ctx, &mockTransport{endpoint: "test_endpoint"})
	tr, ok := FromClientContext(ctx)

	assert.Equal(t, true, ok)
	assert.NotNil(t, tr)
	mtr, ok := tr.(*mockTransport)
	assert.Equal(t, true, ok)
	assert.NotNil(t, mtr)
	assert.Equal(t, mtr.endpoint, "test_endpoint")
}
