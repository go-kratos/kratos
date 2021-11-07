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
	} else {
		// TODO: get from pool
		candidates = make([]WeightedNode, len(nodes))
		copy(candidates, nodes)
	}
	for _, f := range options.Filters {
		candidates = d.nodeFilter(f, candidates)
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

func (d *Default) nodeFilter(filter NodeFilter, nodes []WeightedNode) []WeightedNode {
	length := len(nodes)
	for i := 0; i < length; i++ {
		if !filter((nodes)[i]) {
			if i == length-1 {
				length--
				break
			}
			for ; length > i; length-- {
				if filter((nodes)[length-1]) {
					(nodes)[i] = (nodes)[length-1]
					length--
					break
				}
			}
		}
	}
	nodes = (nodes)[:length]
	return nodes
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
