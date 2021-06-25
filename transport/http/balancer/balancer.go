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
	// Pick one node
	Pick(ctx context.Context) (node *registry.ServiceInstance, done func(context.Context, DoneInfo), err error)
	// Update nodes when nodes removed or added
	Update(nodes []*registry.ServiceInstance)
}
