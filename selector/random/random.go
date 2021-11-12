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

var _ selector.Balancer = &Balancer{} // Name is balancer name

// WithFilter with select filters
func WithFilter(filters ...selector.Filter) Option {
	return func(o *options) {
		o.filters = filters
	}
}

// Option is random builder option.
type Option func(o *options)

// options is random builder options
type options struct {
	filters []selector.Filter
}

// Balancer is a random balancer.
type Balancer struct{}

// New an random selector.
func New(opts ...Option) selector.Selector {
	return NewBuilder(opts...).Build()
}

// Pick is pick a weighted node.
func (p *Balancer) Pick(_ context.Context, nodes []selector.WeightedNode) (selector.WeightedNode, selector.DoneFunc, error) {
	if len(nodes) == 0 {
		return nil, nil, selector.ErrNoAvailable
	}
	cur := rand.Intn(len(nodes))
	selected := nodes[cur]
	d := selected.Pick()
	return selected, d, nil
}

// NewBuilder returns a selector builder with random balancer
func NewBuilder(opts ...Option) selector.Builder {
	var option options
	for _, opt := range opts {
		opt(&option)
	}
	return &selector.DefaultBuilder{
		Filters:  option.filters,
		Balancer: &Builder{},
		Node:     &direct.Builder{},
	}
}

// Builder is random builder
type Builder struct{}

// Build creates Balancer
func (b *Builder) Build() selector.Balancer {
	return &Balancer{}
}
