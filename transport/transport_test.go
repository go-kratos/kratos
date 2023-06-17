package transport

import (
	"context"
	"reflect"
	"testing"
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
	if !ok {
		t.Errorf("expected:%v got:%v", true, ok)
	}
	if tr == nil {
		t.Errorf("expected:%v got:%v", nil, tr)
	}
	mtr, ok := tr.(*mockTransport)
	if !ok {
		t.Errorf("expected:%v got:%v", true, ok)
	}
	if mtr == nil {
		t.Fatalf("expected:%v got:%v", nil, mtr)
	}
	if mtr.Kind().String() != KindGRPC.String() {
		t.Errorf("expected:%v got:%v", KindGRPC.String(), mtr.Kind().String())
	}
	if !reflect.DeepEqual(mtr.endpoint, "test_endpoint") {
		t.Errorf("expected:%v got:%v", "test_endpoint", mtr.endpoint)
	}
}

func TestClientTransport(t *testing.T) {
	ctx := context.Background()

	ctx = NewClientContext(ctx, &mockTransport{endpoint: "test_endpoint"})
	tr, ok := FromClientContext(ctx)
	if !ok {
		t.Errorf("expected:%v got:%v", true, ok)
	}
	if tr == nil {
		t.Errorf("expected:%v got:%v", nil, tr)
	}
	mtr, ok := tr.(*mockTransport)
	if !ok {
		t.Errorf("expected:%v got:%v", true, ok)
	}
	if mtr == nil {
		t.Errorf("expected:%v got:%v", nil, mtr)
	}
	if !reflect.DeepEqual(mtr.endpoint, "test_endpoint") {
		t.Errorf("expected:%v got:%v", "test_endpoint", mtr.endpoint)
	}
}
