package random

import (
	"context"
	"math/rand"

	"github.com/go-kratos/kratos/v2/selector"
	"github.com/go-kratos/kratos/v2/selector/node/direct"
)

const (
	// Name is random balancer name
	Name = "random"
)

var _ selector.Balancer[*direct.Node] = (*Balancer[*direct.Node])(nil) // Name is balancer name

// Option is random builder option.
type Option func(o *options)

// options is random builder options
type options struct{}

// Balancer is a random balancer.
type Balancer[W selector.WeightedNode] struct{}

// New a random selector.
func New(opts ...Option) selector.Selector {
	return NewBuilder(opts...).Build()
}

// Pick is pick a weighted node.
func (p *Balancer[W]) Pick(_ context.Context, nodes []W) (W, selector.DoneFunc, error) {
	if len(nodes) == 0 {
		var zero W
		return zero, nil, selector.ErrNoAvailable
	}
	cur := rand.Intn(len(nodes))
	selected := nodes[cur]
	d := selected.Pick()
	return selected, d, nil
}

// NewBuilder returns a selector builder with random balancer
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

// Builder is random builder
type Builder[W selector.WeightedNode] struct{}

// Build creates Balancer
func (b *Builder[W]) Build() selector.Balancer[W] {
	return &Balancer[W]{}
}
