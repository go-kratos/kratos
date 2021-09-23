package p2c

import (
	"context"
	"math/rand"
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

var _ selector.Balancer = &Balancer{}

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

// New creates a p2c selector.
func New(opts ...Option) selector.Selector {
	var option options
	for _, opt := range opts {
		opt(&option)
	}
	return &selector.Default{
		NodeBuilder: &ewma.Builder{},
		Balancer: &Balancer{
			r: rand.New(rand.NewSource(time.Now().UnixNano())),
		},
		Filters: option.filters,
	}
}

// Balancer is p2c selector.
type Balancer struct {
	r  *rand.Rand
	lk int64
}

// choose two distinct nodes.
func (s *Balancer) prePick(nodes []selector.WeightedNode) (nodeA selector.WeightedNode, nodeB selector.WeightedNode) {
	a := s.r.Intn(len(nodes))
	b := s.r.Intn(len(nodes) - 1)
	if b >= a {
		b = b + 1
	}
	nodeA, nodeB = nodes[a], nodes[b]
	return
}

// Pick pick a node.
func (s *Balancer) Pick(ctx context.Context, nodes []selector.WeightedNode) (selector.WeightedNode, selector.DoneFunc, error) {
	if len(nodes) == 0 {
		return nil, nil, selector.ErrNoAvailable
	} else if len(nodes) == 1 {
		done := nodes[0].Pick()
		return nodes[0], done, nil
	}

	var pc, upc selector.WeightedNode
	nodeA, nodeB := s.prePick(nodes)
	// meta.Weight为服务发布者在discovery中设置的权重
	if nodeB.Weight() > nodeA.Weight() {
		pc, upc = nodeB, nodeA
	} else {
		pc, upc = nodeA, nodeB
	}

	// 如果落选节点在forceGap期间内从来没有被选中一次，则强制选一次
	// 利用强制的机会，来触发成功率、延迟的更新
	if upc.PickElapsed() > forcePick && atomic.CompareAndSwapInt64(&s.lk, 0, 1) {
		pc = upc
		atomic.StoreInt64(&s.lk, 0)
	}
	done := pc.Pick()
	return pc, done, nil
}
