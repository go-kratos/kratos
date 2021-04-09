package http

import (
	"context"
	"net/http"
)

// ServerInfo represent HTTP server information.
type ServerInfo struct {
	Request  *http.Request
	Response http.ResponseWriter
}

type serverKey struct{}

// NewServerContext returns a new Context that carries value.
func NewServerContext(ctx context.Context, info ServerInfo) context.Context {
	return context.WithValue(ctx, serverKey{}, info)
}

// FromServerContext returns the Transport value stored in ctx, if any.
func FromServerContext(ctx context.Context) (info ServerInfo, ok bool) {
	info, ok = ctx.Value(serverKey{}).(ServerInfo)
	return
}

// ClientInfo represent HTTP client information.
type ClientInfo struct {
	Request *http.Request
}

type clientKey struct{}

// NewClientContext returns a new Context that carries value.
func NewClientContext(ctx context.Context, info ClientInfo) context.Context {
	return context.WithValue(ctx, clientKey{}, info)
}

// FromClientContext returns the Transport value stored in ctx, if any.
func FromClientContext(ctx context.Context) (info ClientInfo, ok bool) {
	info, ok = ctx.Value(clientKey{}).(ClientInfo)
	return
}
