package log

import "context"

type loggerKey struct{}

// NewContext returns a new Context that carries value.
func NewContext(ctx context.Context, l Logger) context.Context {
	return context.WithValue(ctx, loggerKey{}, l)
}

// FromContext returns the Transport value stored in ctx, if any.
func FromContext(ctx context.Context) (l Logger, ok bool) {
	l, ok = ctx.Value(loggerKey{}).(Logger)
	return
}
