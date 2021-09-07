package http

import (
	"context"
	"sync/atomic"

	"github.com/go-kratos/kratos/v2/balancer"
	"github.com/go-kratos/kratos/v2/metadata"
)

type picker struct {
	selector    balancer.Selector
	nodeBuilder balancer.NodeBuilder

	nodes atomic.Value
}

func newPicker(selector balancer.Selector, nodeBuilder balancer.NodeBuilder) *picker {
	p := &picker{
		selector:    selector,
		nodeBuilder: nodeBuilder,
	}
	p.nodes.Store([]balancer.Node{})
	return p
}

func (p *picker) Update(readys []node) {
	if len(readys) == 0 {
		return
	}
	nodes := []balancer.Node{}
	for _, n := range readys {
		nodes = append(nodes, p.nodeBuilder.Build(n.addr, 100, metadata.New(n.metadata)))
	}
	p.nodes.Store(nodes)
}

func (p *picker) Pick(ctx context.Context) (addr string, done func(ctx context.Context, di balancer.DoneInfo), err error) {
	var n balancer.Node
	n, done, err = p.selector.Select(ctx, p.nodes.Load().([]balancer.Node))
	if err != nil {
		return
	}
	addr = n.Address()
	return
}
