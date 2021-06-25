package random

import (
	"context"
	"fmt"
	"math/rand"
	"sync"

	"github.com/go-kratos/kratos/v2/registry"
	"github.com/go-kratos/kratos/v2/transport/http/balancer"
)

var _ balancer.Balancer = &Balancer{}

type Balancer struct {
	lock  sync.RWMutex
	nodes []*registry.ServiceInstance
}

func New() *Balancer {
	return &Balancer{}
}

func (b *Balancer) Pick(ctx context.Context) (node *registry.ServiceInstance, done func(context.Context, balancer.DoneInfo), err error) {
	b.lock.RLock()
	nodes := b.nodes
	b.lock.RUnlock()

	if len(nodes) == 0 {
		return nil, nil, fmt.Errorf("no instances avaiable")
	}
	if len(nodes) == 1 {
		return nodes[0], func(context.Context, balancer.DoneInfo) {}, nil
	}
	idx := rand.Intn(len(nodes))
	return nodes[idx], func(context.Context, balancer.DoneInfo) {}, nil
}

func (b *Balancer) Update(nodes []*registry.ServiceInstance) {
	b.lock.Lock()
	defer b.lock.Unlock()
	b.nodes = nodes
}
