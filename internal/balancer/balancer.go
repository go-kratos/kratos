package balancer

import (
	"context"

	"github.com/go-kratos/kratos/v2/registry"
)

// DoneInfo is callback when rpc done
type DoneInfo struct {
	Err     error
	Trailer map[string]string
}

// Balancer is node pick balancer
type Balancer interface {
	Pick(ctx context.Context, pathPattern string, nodes []registry.Service) (node registry.Service, done func(context.Context, DoneInfo), err error)
}
