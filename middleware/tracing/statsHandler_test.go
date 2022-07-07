package tracing

import (
	"context"
	"net"
	"testing"

	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/stats"
)

func Test_Client_HandleConn(t *testing.T) {
	(&ClientHandler{}).HandleConn(context.Background(), nil)
}

func Test_Client_TagConn(t *testing.T) {
	client := &ClientHandler{}
	ctx := context.WithValue(context.Background(), "MY_TEST_KEY", 123)

	if client.TagConn(ctx, nil).Value("MY_TEST_KEY") != 123 {
		t.Errorf(`The context value must be 123 for the "MY_KEY_TEST" key, %v given.`, client.TagConn(ctx, nil).Value("MY_TEST_KEY"))
	}
}

func Test_Client_TagRPC(t *testing.T) {
	client := &ClientHandler{}
	ctx := context.WithValue(context.Background(), "MY_TEST_KEY", 123)

	if client.TagRPC(ctx, nil).Value("MY_TEST_KEY") != 123 {
		t.Errorf(`The context value must be 123 for the "MY_KEY_TEST" key, %v given.`, client.TagConn(ctx, nil).Value("MY_TEST_KEY"))
	}
}

type (
	mockSpan struct {
		trace.Span
		mockSpanCtx *trace.SpanContext
	}
)

func (m *mockSpan) SpanContext() trace.SpanContext {
	return *m.mockSpanCtx
}

func Test_Client_HandleRPC(t *testing.T) {
	client := &ClientHandler{}
	ctx := context.Background()
	rs := stats.OutHeader{}

	// Handle stats.RPCStats is not type of stats.OutHeader case
	client.HandleRPC(nil, nil)

	// Handle context doesn't have the peerkey filled with a Peer instance
	client.HandleRPC(ctx, &rs)

	// Handle context with the peerkey filled with a Peer instance
	ip, _ := net.ResolveIPAddr("ip", "1.1.1.1")
	ctx = peer.NewContext(ctx, &peer.Peer{
		Addr: ip,
	})
	client.HandleRPC(ctx, &rs)

	// Handle context with Span
	_, span := trace.NewNoopTracerProvider().Tracer("Tracer").Start(ctx, "Spanname")
	spanCtx := trace.SpanContext{}
	spanId := [8]byte{12, 12, 12, 12, 12, 12, 12, 12}
	traceId := [16]byte{12, 12, 12, 12, 12, 12, 12, 12, 12, 12, 12, 12, 12, 12, 12, 12}
	spanCtx = spanCtx.WithTraceID(traceId)
	spanCtx = spanCtx.WithSpanID(spanId)
	mSpan := mockSpan{
		Span:        span,
		mockSpanCtx: &spanCtx,
	}
	ctx = trace.ContextWithSpan(ctx, &mSpan)
	client.HandleRPC(ctx, &rs)
}
