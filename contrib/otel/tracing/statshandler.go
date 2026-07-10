package tracing

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/stats"
)

// ClientHandler is tracing ClientHandler
type ClientHandler struct{}

// HandleConn exists to satisfy gRPC stats.Handler.
func (c *ClientHandler) HandleConn(_ context.Context, _ stats.ConnStats) {
	fmt.Println("Handle connection.")
}

// TagConn exists to satisfy gRPC stats.Handler.
func (c *ClientHandler) TagConn(ctx context.Context, _ *stats.ConnTagInfo) context.Context {
	return ctx
}

// HandleRPC implements per-RPC tracing and stats instrumentation.
func (c *ClientHandler) HandleRPC(ctx context.Context, rs stats.RPCStats) {
	if _, ok := rs.(*stats.OutHeader); !ok {
		return
	}
	p, ok := peer.FromContext(ctx)
	if !ok {
		return
	}
	span := trace.SpanFromContext(ctx)
	if span.SpanContext().IsValid() {
		span.SetAttributes(peerAttr(p.Addr.String())...)
	}
}

// TagRPC implements per-RPC context management.
func (c *ClientHandler) TagRPC(ctx context.Context, _ *stats.RPCTagInfo) context.Context {
	return ctx
}
