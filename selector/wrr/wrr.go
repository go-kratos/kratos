package wrr

import (
	"context"
	"sync"

	"github.com/go-kratos/kratos/v3/selector"
	"github.com/go-kratos/kratos/v3/selector/node/direct"
)

const (
	// Name is wrr(Weighted Round Robin) balancer name
	Name = "wrr"
)

var _ selector.Balancer = (*Balancer)(nil)

// Option is wrr builder option.
type Option func(o *options)

// options is wrr builder options
type options struct{}

// Balancer is a wrr balancer.
type Balancer struct {
	mu            sync.Mutex
	currentWeight map[string]float64
}

// New random a selector.
func New(opts ...Option) selector.Selector {
	return NewBuilder(opts...).Build()
}

// Pick is pick a weighted node.
func (p *Balancer) Pick(_ context.Context, nodes []selector.WeightedNode) (selector.WeightedNode, selector.DoneFunc, error) {
	if len(nodes) == 0 {
		return nil, nil, selector.ErrNoAvailable
	}

	p.mu.Lock()
	defer p.mu.Unlock()

	var totalWeight float64
	var selected selector.WeightedNode
	var selectWeight float64

	// nginx wrr load balancing algorithm: http://blog.csdn.net/zhangskd/article/details/50194069
	for _, node := range nodes {
		totalWeight += node.Weight()
		cwt := p.currentWeight[node.Address()]
		// current += effectiveWeight
		cwt += node.Weight()
		p.currentWeight[node.Address()] = cwt
		if selected == nil || selectWeight < cwt {
			selectWeight = cwt
			selected = node
		}
	}
	p.currentWeight[selected.Address()] = selectWeight - totalWeight

	// After the loop, currentWeight has an entry for every current node, plus any
	// leftover entries for nodes that have disappeared from service discovery. So
	// len(currentWeight) > len(nodes) exactly when stale entries exist: drop them
	// to keep the map from growing without bound as nodes churn. When the node set
	// is unchanged (the common case) the sizes match and cleanup is skipped, so the
	// per-pick cost is just the algorithm itself with no extra bookkeeping.
	if len(p.currentWeight) > len(nodes) {
		p.cleanupStaleEntries(nodes)
	}

	d := selected.Pick()
	return selected, d, nil
}

// cleanupStaleEntries removes currentWeight entries whose node is no longer present.
func (p *Balancer) cleanupStaleEntries(nodes []selector.WeightedNode) {
	current := make(map[string]struct{}, len(nodes))
	for _, node := range nodes {
		current[node.Address()] = struct{}{}
	}
	for address := range p.currentWeight {
		if _, ok := current[address]; !ok {
			delete(p.currentWeight, address)
		}
	}
}

// NewBuilder returns a selector builder with wrr balancer
func NewBuilder(opts ...Option) selector.Builder {
	var option options
	for _, opt := range opts {
		opt(&option)
	}
	return &selector.DefaultBuilder{
		Balancer: &Builder{},
		Node:     &direct.Builder{},
	}
}

// Builder is wrr builder
type Builder struct{}

// Build creates Balancer
func (b *Builder) Build() selector.Balancer {
	return &Balancer{currentWeight: make(map[string]float64)}
}
