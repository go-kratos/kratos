package random

import (
	"context"
	"math/rand"

	"github.com/go-kratos/kratos/v2/balancer"
)

var (
	_ balancer.Selector = &Selector{}
)

type Selector struct {
}

func New() *Selector {
	return &Selector{}
}

func (p *Selector) Select(ctx context.Context, nodes []balancer.Node) (selected balancer.Node, err error) {
	if len(nodes) == 0 {
		err = balancer.ErrNoAvaliable
		return
	}
	cur := rand.Intn(len(nodes))
	selected = nodes[cur]
	return
}
