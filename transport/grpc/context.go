package grpc

import "context"

// ServerInfo is HTTP server infomation.
type ServerInfo struct {
	// Server is the service implementation the user provides. This is read-only.
	Server interface{}
	// FullMethod is the full RPC method string, i.e., /package.service/method.
	FullMethod string
}

type serverKey struct{}

// NewContext returns a new Context that carries value.
func NewContext(ctx context.Context, s ServerInfo) context.Context {
	return context.WithValue(ctx, serverKey{}, s)
}

// FromContext returns the Transport value stored in ctx, if any.
func FromContext(ctx context.Context) (s ServerInfo, ok bool) {
	s, ok = ctx.Value(serverKey{}).(ServerInfo)
	return
}
