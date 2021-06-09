package middleware

import "context"

// ServiceInfo represent service information.
type ServiceInfo struct {
	// FullMethod is the full RPC method string, i.e., /package.service/method.
	FullMethod string
}

type serviceKey struct{}

// NewContext returns a new Context that carries value.
func NewContext(ctx context.Context, info ServiceInfo) context.Context {
	return context.WithValue(ctx, serviceKey{}, info)
}

// FromContext returns the Service value stored in ctx, if any.
func FromContext(ctx context.Context) (info ServiceInfo, ok bool) {
	info, ok = ctx.Value(serviceKey{}).(ServiceInfo)
	return
}
