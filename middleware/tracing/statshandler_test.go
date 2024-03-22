package tracing

import (
	"context"
	"net"
	"testing"

	"go.opentelemetry.io/otel/trace"
	"go.opentelemetry.io/otel/trace/noop"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/stats"
)

type ctxKey string

const testKey ctxKey = "MY_TEST_KEY"

func TestClient_HandleConn(_ *testing.T) {
	(&ClientHandler{}).HandleConn(context.Background(), nil)
}

func TestClient_TagConn(t *testing.T) {
	client := &ClientHandler{}
	ctx := context.WithValue(context.Background(), testKey, 123)

	if client.TagConn(ctx, nil).Value(testKey) != 123 {
		t.Errorf(`The context value must be 123 for the "MY_KEY_TEST" key, %v given.`, client.TagConn(ctx, nil).Value(testKey))
	}
}

func TestClient_TagRPC(t *testing.T) {
	client := &ClientHandler{}
	ctx := context.WithValue(context.Background(), testKey, 123)

	if client.TagRPC(ctx, nil).Value(testKey) != 123 {
		t.Errorf(`The context value must be 123 for the "MY_KEY_TEST" key, %v given.`, client.TagConn(ctx, nil).Value(testKey))
	}
}

type mockSpan struct {
	trace.Span
	mockSpanCtx *trace.SpanContext
}

func (m *mockSpan) SpanContext() trace.SpanContext {
	return *m.mockSpanCtx
}

func TestClient_HandleRPC(_ *testing.T) {
	client := &ClientHandler{}
	ctx := context.Background()
	rs := stats.OutHeader{}

	// Handle stats.RPCStats is not type of stats.OutHeader case
	client.HandleRPC(context.TODO(), nil)

	// Handle context doesn't have the peerkey filled with a Peer instance
	client.HandleRPC(ctx, &rs)

	// Handle context with the peerkey filled with a Peer instance
	ip, _ := net.ResolveIPAddr("ip", "1.1.1.1")
	ctx = peer.NewContext(ctx, &peer.Peer{
		Addr: ip,
	})
	client.HandleRPC(ctx, &rs)

	// Handle context with Span
	_, span := noop.NewTracerProvider().Tracer("Tracer").Start(ctx, "Spanname")
	spanCtx := trace.SpanContext{}
	spanID := [8]byte{12, 12, 12, 12, 12, 12, 12, 12}
	traceID := [16]byte{12, 12, 12, 12, 12, 12, 12, 12, 12, 12, 12, 12, 12, 12, 12, 12}
	spanCtx = spanCtx.WithTraceID(traceID)
	spanCtx = spanCtx.WithSpanID(spanID)
	mSpan := mockSpan{
		Span:        span,
		mockSpanCtx: &spanCtx,
	}
	ctx = trace.ContextWithSpan(ctx, &mSpan)
	client.HandleRPC(ctx, &rs)
}
