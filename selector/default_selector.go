package selector

import (
	"context"
	"sync/atomic"
)

// Default is composite selector.
type Default struct {
	NodeBuilder WeightedNodeBuilder
	Balancer    Balancer
	Filters     []Filter

	nodes atomic.Value
}

// Select select one node.
func (d *Default) Select(ctx context.Context, opts ...SelectOption) (selected Node, done DoneFunc, err error) {
	nodes, _ := d.nodes.Load().([]Node)
	var options SelectOptions
	for _, o := range opts {
		o(&options)
	}
	if len(d.Filters) > 0 || len(options.Filters) > 0 {
		// TODO: get from pool
		newNodes := make([]Node, len(nodes))
		copy(newNodes, nodes)
		for _, f := range d.Filters {
			f(ctx, &newNodes)
		}
		for _, f := range options.Filters {
			f(ctx, &newNodes)
		}
		nodes = newNodes
	}
	if len(nodes) == 0 {
		return nil, nil, ErrNoAvailable
	}
	// TODO: get from pool
	candidates := make([]WeightedNode, 0, len(nodes))
	for _, n := range nodes {
		candidates = append(candidates, n.(WeightedNode))
	}
	wn, done, err := d.Balancer.Pick(ctx, candidates)
	if err != nil {
		return nil, nil, err
	}
	return wn.Raw(), done, nil
}

// Apply update nodes info.
func (d *Default) Apply(nodes []Node) {
	weightedNodes := make([]Node, 0, len(nodes))
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
	Filters  []Filter
}

// Build create builder
func (db *DefaultBuilder) Build() Selector {
	return &Default{
		NodeBuilder: db.Node,
		Balancer:    db.Balancer.Build(),
		Filters:     db.Filters,
	}
}
