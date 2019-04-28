package context

import (
	"context"
	"time"

	"go-common/library/net/rpc/liverpc"
)

// WithHeader returns new context with header
// Deprecated: Use HeaderOption instead
func WithHeader(ctx context.Context, header *liverpc.Header) (ret context.Context) {
	ret = context.WithValue(ctx, liverpc.KeyHeader, header)
	return
}

// WithTimeout set timeout to rpc request
// Notice this is nothing related to to built-in context.WithTimeout
// Deprecated: Use TimeoutOption instead
func WithTimeout(ctx context.Context, time time.Duration) (ret context.Context) {
	ret = context.WithValue(ctx, liverpc.KeyTimeout, time)
	return
}
