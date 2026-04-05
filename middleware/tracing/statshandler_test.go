package tracing

import (
	"context"
	"net"
	"testing"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.opentelemetry.io/otel/trace/noop"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/stats"
)

type ctxKey string

const testKey ctxKey = "MY_TEST_KEY"

type recordingSpan struct {
	trace.Span
	spanCtx    trace.SpanContext
	attributes []attribute.KeyValue
}

func (r *recordingSpan) SpanContext() trace.SpanContext {
	return r.spanCtx
}

func (r *recordingSpan) SetAttributes(attrs ...attribute.KeyValue) {
	r.attributes = append(r.attributes, attrs...)
}

func newValidSpanContext() trace.SpanContext {
	sc := trace.SpanContext{}
	sc = sc.WithTraceID([16]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16})
	sc = sc.WithSpanID([8]byte{1, 2, 3, 4, 5, 6, 7, 8})
	return sc
}

func TestClient_HandleConn(t *testing.T) {
	client := &ClientHandler{}
	// HandleConn should not panic with nil stats
	client.HandleConn(context.Background(), nil)
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

func TestClient_HandleRPC(t *testing.T) {
	client := &ClientHandler{}

	t.Run("non-OutHeader is ignored", func(t *testing.T) {
		span := &recordingSpan{
			Span:    noop.Span{},
			spanCtx: newValidSpanContext(),
		}
		ctx := trace.ContextWithSpan(context.Background(), span)
		ctx = peer.NewContext(ctx, &peer.Peer{Addr: &net.TCPAddr{IP: net.ParseIP("1.1.1.1"), Port: 8080}})
		// pass nil instead of *stats.OutHeader
		client.HandleRPC(ctx, nil)
		if len(span.attributes) != 0 {
			t.Errorf("expected no attributes for non-OutHeader, got %v", span.attributes)
		}
	})

	t.Run("no peer in context", func(t *testing.T) {
		span := &recordingSpan{
			Span:    noop.Span{},
			spanCtx: newValidSpanContext(),
		}
		ctx := trace.ContextWithSpan(context.Background(), span)
		client.HandleRPC(ctx, &stats.OutHeader{})
		if len(span.attributes) != 0 {
			t.Errorf("expected no attributes without peer, got %v", span.attributes)
		}
	})

	t.Run("invalid span context", func(t *testing.T) {
		span := &recordingSpan{
			Span:    noop.Span{},
			spanCtx: trace.SpanContext{}, // invalid: zero trace/span IDs
		}
		ctx := trace.ContextWithSpan(context.Background(), span)
		ctx = peer.NewContext(ctx, &peer.Peer{Addr: &net.TCPAddr{IP: net.ParseIP("1.1.1.1"), Port: 8080}})
		client.HandleRPC(ctx, &stats.OutHeader{})
		if len(span.attributes) != 0 {
			t.Errorf("expected no attributes for invalid span context, got %v", span.attributes)
		}
	})

	t.Run("valid span sets peer attributes", func(t *testing.T) {
		span := &recordingSpan{
			Span:    noop.Span{},
			spanCtx: newValidSpanContext(),
		}
		ctx := trace.ContextWithSpan(context.Background(), span)
		ctx = peer.NewContext(ctx, &peer.Peer{Addr: &net.TCPAddr{IP: net.ParseIP("1.1.1.1"), Port: 8080}})
		client.HandleRPC(ctx, &stats.OutHeader{})
		if len(span.attributes) == 0 {
			t.Fatal("expected peer attributes to be set, got none")
		}
		var foundIP, foundPort bool
		for _, attr := range span.attributes {
			if attr.Key == "net.peer.ip" && attr.Value.AsString() == "1.1.1.1" {
				foundIP = true
			}
			if attr.Key == "net.peer.port" && attr.Value.AsString() == "8080" {
				foundPort = true
			}
		}
		if !foundIP {
			t.Error("expected net.peer.ip attribute with value 1.1.1.1")
		}
		if !foundPort {
			t.Error("expected net.peer.port attribute with value 8080")
		}
	})
}
