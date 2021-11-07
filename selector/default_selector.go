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
	nodes, _ := d.nodes.Load().([]WeightedNode)
	if nodes == nil {
		return nil, nil, ErrNoAvailable
	}
	var options SelectOptions
	for _, o := range opts {
		o(&options)
	}

	var candidates []WeightedNode
	if len(d.Filters) > 0 {
		newNodes := make([]Node, len(nodes))
		for i, wc := range nodes {
			newNodes[i] = wc
		}
		for _, f := range d.Filters {
			newNodes = f(ctx, newNodes)
		}
		// TODO: get from pool
		candidates = make([]WeightedNode, len(newNodes))
		for i, n := range newNodes {
			candidates[i] = n.(WeightedNode)
		}
	}
	if len(options.Filters) > 0 {
		candidates = d.nodeFilter(options.Filters, candidates)
	}
	if len(candidates) == 0 {
		return nil, nil, ErrNoAvailable
	}
	wn, done, err := d.Balancer.Pick(ctx, candidates)
	if err != nil {
		return nil, nil, err
	}
	return wn.Raw(), done, nil
}

func (d *Default) nodeFilter(filters []NodeFilter, nodes []WeightedNode) []WeightedNode {
	// TODO: get from pool
	newNodes := make([]WeightedNode, 0, len(nodes))

	for _, n := range nodes {
		remove := false
		for _, f := range filters {
			if !f(n) {
				remove = true
				break
			}
		}
		if !remove {
			newNodes = append(newNodes, n)
		}

	}
	return newNodes
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
