package random

import (
	"context"
	"math/rand"

	"github.com/go-kratos/kratos/v2/balancer"
)

var (
	_ balancer.Selector = &Selector{}

	// Name is balancer name
	Name = "random"
)

type Selector struct{}

func New() *Selector {
	return &Selector{}
}

func (p *Selector) Select(_ context.Context, nodes []balancer.Node) (balancer.Node, func(context.Context, balancer.DoneInfo), error) {
	if len(nodes) == 0 {
		err := balancer.ErrNoAvaliable
		return nil, nil, err
	}
	cur := rand.Intn(len(nodes))
	selected := nodes[cur]
	d := selected.Pick()
	done := func(ctx context.Context, info balancer.DoneInfo) {
		d(ctx, info)
	}
	return selected, done, nil
}
