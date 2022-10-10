package selector

import (
	"context"
	"sync/atomic"
)

var (
	_ Rebalancer = (*Default)(nil)
	_ Builder    = (*DefaultBuilder)(nil)
)

// Default is composite selector.
type Default struct {
	NodeBuilder WeightedNodeBuilder
	Balancer    Balancer

	nodes atomic.Value
}

// Select is select one node.
func (d *Default) Select(ctx context.Context, opts ...SelectOption) (selected Node, done DoneFunc, err error) {
	var (
		options    SelectOptions
		candidates []WeightedNode
	)
	nodes, ok := d.nodes.Load().([]WeightedNode)
	if !ok {
		return nil, nil, ErrNoAvailable
	}
	for _, o := range opts {
		o(&options)
	}
	if len(options.NodeFilters) > 0 {
		newNodes := make([]Node, len(nodes))
		for i, wc := range nodes {
			newNodes[i] = wc
		}
		for _, filter := range options.NodeFilters {
			newNodes = filter(ctx, newNodes)
		}
		candidates = make([]WeightedNode, len(newNodes))
		for i, n := range newNodes {
			candidates[i] = n.(WeightedNode)
		}
	} else {
		candidates = nodes
	}

	if len(candidates) == 0 {
		return nil, nil, ErrNoAvailable
	}
	wn, done, err := d.Balancer.Pick(ctx, candidates)
	if err != nil {
		return nil, nil, err
	}
	p, ok := FromPeerContext(ctx)
	if ok {
		p.Node = wn.Raw()
	}
	return wn.Raw(), done, nil
}

// Apply update nodes info.
func (d *Default) Apply(nodes []Node) {
	weightedNodes := make([]WeightedNode, 0, len(nodes))
	for _, n := range nodes {
		weightedNodes = append(weightedNodes, d.NodeBuilder.Build(n))
	}
	// TODO: Do not delete unchanged nodes
	d.nodes.Store(weightedNodes)
}

// DefaultBuilder is de
type DefaultBuilder struct {
	Node     WeightedNodeBuilder
	Balancer BalancerBuilder
}

// Build create builder
func (db *DefaultBuilder) Build() Selector {
	return &Default{
		NodeBuilder: db.Node,
		Balancer:    db.Balancer.Build(),
	}
}
