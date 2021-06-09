package random

import (
	"context"
	"fmt"
	"math/rand"

	"github.com/go-kratos/kratos/v2/registry"
	"github.com/go-kratos/kratos/v2/transport/http/balancer"
)

var _ balancer.Balancer = &Balancer{}

type Balancer struct {
}

func New() *Balancer {
	return &Balancer{}
}

func (b *Balancer) Pick(ctx context.Context, nodes []*registry.ServiceInstance) (node *registry.ServiceInstance, done func(context.Context, balancer.DoneInfo), err error) {
	if len(nodes) == 0 {
		return nil, nil, fmt.Errorf("no instances avaiable")
	}
	if len(nodes) == 1 {
		return nodes[0], func(context.Context, balancer.DoneInfo) {}, nil
	}
	idx := rand.Intn(len(nodes))
	return nodes[idx], func(context.Context, balancer.DoneInfo) {}, nil
}
