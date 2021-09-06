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

type Selector struct {
}

func New() *Selector {
	return &Selector{}
}

func (p *Selector) Select(ctx context.Context, nodes []balancer.Node) (selected balancer.Node, done func(ctx context.Context, di balancer.DoneInfo), err error) {
	if len(nodes) == 0 {
		err = balancer.ErrNoAvaliable
		return
	}
	cur := rand.Intn(len(nodes))
	selected = nodes[cur]
	done = func(context.Context, balancer.DoneInfo) {}
	return
}
