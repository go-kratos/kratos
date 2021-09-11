package random

import (
	"context"
	"math/rand"

	"github.com/go-kratos/kratos/v2/selector"
	"github.com/go-kratos/kratos/v2/selector/node/direct"
)

var (
	_ selector.Balancer = &Balancer{}

	// Name is balancer name
	Name = "random"
)

// Balancer is random balancer
type Balancer struct{}

// New random selector
func New() selector.Selector {
	return &selector.Default{
		Balancer:    &Balancer{},
		NodeBuilder: &direct.Builder{},
	}
}

// Pick one node
func (p *Balancer) Pick(_ context.Context, nodes []selector.WeightedNode) (selector.WeightedNode, selector.DoneFunc, error) {
	if len(nodes) == 0 {
		err := selector.ErrNoAvailable
		return nil, nil, err
	}
	cur := rand.Intn(len(nodes))
	selected := nodes[cur]
	d := selected.Pick()
	return selected, d, nil
}
