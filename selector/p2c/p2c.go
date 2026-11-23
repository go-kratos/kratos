package p2c

import (
	"context"
	"math/rand/v2"
	"sync/atomic"
	"time"

	"github.com/go-kratos/kratos/v3/selector"
	"github.com/go-kratos/kratos/v3/selector/node/ewma"
)

const (
	forcePick = time.Second * 3
	// Name is p2c(Pick of 2 choices) balancer name
	Name = "p2c"
)

var _ selector.Balancer = (*Balancer)(nil)

// Option is p2c builder option.
type Option func(o *options)

// options is p2c builder options
type options struct{}

// New creates a p2c selector.
func New(opts ...Option) selector.Selector {
	return NewBuilder(opts...).Build()
}

// Balancer is p2c selector.
type Balancer struct {
	picked atomic.Bool
}

// choose two distinct nodes.
func (s *Balancer) prePick(nodes []selector.WeightedNode) (nodeA selector.WeightedNode, nodeB selector.WeightedNode) {
	// rand/v2 top-level functions are safe for concurrent use and lock-free, so
	// no per-balancer mutex/Rand is needed - picks no longer serialize on a lock.
	a := rand.IntN(len(nodes))
	b := rand.IntN(len(nodes) - 1)
	if b >= a {
		b = b + 1
	}
	nodeA, nodeB = nodes[a], nodes[b]
	return
}

// Pick pick a node.
func (s *Balancer) Pick(_ context.Context, nodes []selector.WeightedNode) (selector.WeightedNode, selector.DoneFunc, error) {
	if len(nodes) == 0 {
		return nil, nil, selector.ErrNoAvailable
	}
	if len(nodes) == 1 {
		done := nodes[0].Pick()
		return nodes[0], done, nil
	}

	var pc, upc selector.WeightedNode
	nodeA, nodeB := s.prePick(nodes)
	// meta.Weight is the weight set by the service publisher in discovery
	if nodeB.Weight() > nodeA.Weight() {
		pc, upc = nodeB, nodeA
	} else {
		pc, upc = nodeA, nodeB
	}

	// If the failed node has never been selected once during forceGap, it is forced to be selected once
	// Take advantage of forced opportunities to trigger updates of success rate and delay
	if upc.PickElapsed() > forcePick && s.picked.CompareAndSwap(false, true) {
		defer s.picked.Store(false)
		pc = upc
	}
	done := pc.Pick()
	return pc, done, nil
}

// NewBuilder returns a selector builder with p2c balancer
func NewBuilder(opts ...Option) selector.Builder {
	var option options
	for _, opt := range opts {
		opt(&option)
	}
	return &selector.DefaultBuilder{
		Balancer: &Builder{},
		Node:     &ewma.Builder{},
	}
}

// Builder is p2c builder
type Builder struct{}

// Build creates Balancer
func (b *Builder) Build() selector.Balancer {
	return &Balancer{}
}
