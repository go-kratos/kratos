package kratos

import (
	"context"
)

// AppInfo is application context value.
type AppInfo struct {
	ID        string
	Name      string
	Version   string
	Metadata  map[string]string
	Endpoints []string
}

type appKey struct{}

// NewContext returns a new Context that carries value.
func NewContext(ctx context.Context, s AppInfo) context.Context {
	return context.WithValue(ctx, appKey{}, s)
}

// FromContext returns the Transport value stored in ctx, if any.
func FromContext(ctx context.Context) (s AppInfo, ok bool) {
	s, ok = ctx.Value(appKey{}).(AppInfo)
	return
}
