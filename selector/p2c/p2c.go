package p2c

import (
	"context"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"

	"github.com/go-kratos/kratos/v2/selector"
	"github.com/go-kratos/kratos/v2/selector/node/ewma"
)

const (
	forcePick = time.Second * 3
	// Name is balancer name
	Name = "p2c"
)

var _ selector.Balancer[*ewma.Node] = (*Balancer[*ewma.Node])(nil)

// Option is random builder option.
type Option func(o *options)

// options is random builder options
type options struct{}

// New creates a p2c selector.
func New(opts ...Option) selector.Selector {
	return NewBuilder(opts...).Build()
}

// Balancer is p2c selector.
type Balancer[W selector.WeightedNode] struct {
	mu     sync.Mutex
	r      *rand.Rand
	picked int64
}

// choose two distinct nodes.
func (s *Balancer[W]) prePick(nodes []W) (nodeA W, nodeB W) {
	s.mu.Lock()
	a := s.r.Intn(len(nodes))
	b := s.r.Intn(len(nodes) - 1)
	s.mu.Unlock()
	if b >= a {
		b = b + 1
	}
	nodeA, nodeB = nodes[a], nodes[b]
	return
}

// Pick pick a node.
func (s *Balancer[W]) Pick(ctx context.Context, nodes []W) (W, selector.DoneFunc, error) {
	if len(nodes) == 0 {
		var zero W
		return zero, nil, selector.ErrNoAvailable
	}
	if len(nodes) == 1 {
		done := nodes[0].Pick()
		return nodes[0], done, nil
	}

	var pc, upc W
	nodeA, nodeB := s.prePick(nodes)
	// meta.Weight is the weight set by the service publisher in discovery
	if nodeB.Weight() > nodeA.Weight() {
		pc, upc = nodeB, nodeA
	} else {
		pc, upc = nodeA, nodeB
	}

	// If the failed node has never been selected once during forceGap, it is forced to be selected once
	// Take advantage of forced opportunities to trigger updates of success rate and delay
	if upc.PickElapsed() > forcePick && atomic.CompareAndSwapInt64(&s.picked, 0, 1) {
		pc = upc
		atomic.StoreInt64(&s.picked, 0)
	}
	done := pc.Pick()
	return pc, done, nil
}

// NewBuilder returns a selector builder with p2c balancer
func NewBuilder(opts ...Option) selector.Builder {
	return NewWithBuilder[*ewma.Node](&ewma.Builder{}, opts...)
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

// Builder is p2c builder
type Builder[W selector.WeightedNode] struct{}

// Build creates Balancer
func (b *Builder[W]) Build() selector.Balancer[W] {
	return &Balancer[W]{r: rand.New(rand.NewSource(time.Now().UnixNano()))}
}
