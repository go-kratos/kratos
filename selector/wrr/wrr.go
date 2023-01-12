package wrr

import (
	"context"
	"sync"

	"github.com/go-kratos/kratos/v2/selector"
	"github.com/go-kratos/kratos/v2/selector/node/direct"
)

const (
	// Name is wrr balancer name
	Name = "wrr"
)

var _ selector.Balancer[*direct.Node] = (*Balancer[*direct.Node])(nil) // Name is balancer name

// Option is random builder option.
type Option func(o *options)

// options is random builder options
type options struct{}

// Balancer is a random balancer.
type Balancer[W selector.WeightedNode] struct {
	mu            sync.Mutex
	currentWeight map[string]float64
}

// New random a selector.
func New(opts ...Option) selector.Selector {
	return NewBuilder(opts...).Build()
}

// Pick is pick a weighted node.
func (p *Balancer[W]) Pick(_ context.Context, nodes []W) (W, selector.DoneFunc, error) {
	if len(nodes) == 0 {
		var zero W
		return zero, nil, selector.ErrNoAvailable
	}
	var totalWeight float64
	var selected W
	var selectWeight float64

	// nginx wrr load balancing algorithm: http://blog.csdn.net/zhangskd/article/details/50194069
	p.mu.Lock()
	for _, node := range nodes {
		totalWeight += node.Weight()
		cwt := p.currentWeight[node.Address()]
		// current += effectiveWeight
		cwt += node.Weight()
		p.currentWeight[node.Address()] = cwt
		if selector.IsNil(selected) || selectWeight < cwt {
			selectWeight = cwt
			selected = node
		}
	}
	p.currentWeight[selected.Address()] = selectWeight - totalWeight
	p.mu.Unlock()

	d := selected.Pick()
	return selected, d, nil
}

// NewBuilder returns a selector builder with wrr balancer
func NewBuilder(opts ...Option) selector.Builder {
	return NewWithBuilder[*direct.Node](&direct.Builder{}, opts...)
}

func NewWithBuilder[W selector.WeightedNode](weightedNodeBuilder selector.WeightedNodeBuilder[W], opts ...Option) selector.Builder {
	var option options
	for _, opt := range opts {
		opt(&option)
	}
	return &selector.DefaultBuilder[W]{
		Balancer: &Builder[W]{},
		Node:     weightedNodeBuilder,
	}
}

// Builder is wrr builder
type Builder[W selector.WeightedNode] struct{}

// Build creates Balancer
func (b *Builder[W]) Build() selector.Balancer[W] {
	return &Balancer[W]{currentWeight: make(map[string]float64)}
}
