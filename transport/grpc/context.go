package grpc

import "context"

// ServerInfo is gRPC server infomation.
type ServerInfo struct {
	// Server is the service implementation the user provides. This is read-only.
	Server interface{}
	// FullMethod is the full RPC method string, i.e., /package.service/method.
	FullMethod string
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

// ClientInfo is gRPC server infomation.
type ClientInfo struct {
	// FullMethod is the full RPC method string, i.e., /package.service/method.
	FullMethod string
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
