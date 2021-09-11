package direct

import (
	"context"
	"sync/atomic"
	"time"

	"github.com/go-kratos/kratos/v2/selector"
)

const (
	defaultWeight = 100
)

var (
	_ selector.WeightedNode        = &node{}
	_ selector.WeightedNodeBuilder = &Builder{}
)

// node is endpoint instance
type node struct {
	selector.Node

	// last lastPick timestamp
	lastPick int64
}

// Builder is direct node builder
type Builder struct{}

// Build create node
func (*Builder) Build(n selector.Node) selector.WeightedNode {
	return &node{Node: n, lastPick: 0}
}

func (n *node) Pick() selector.DoneFunc {
	now := time.Now().UnixNano()
	atomic.StoreInt64(&n.lastPick, now)
	return func(ctx context.Context, di selector.DoneInfo) {}
}

// Weight is node effective weight
func (n *node) Weight() float64 {
	if n.InitialWeight() != nil {
		return float64(*n.InitialWeight())
	}
	return defaultWeight
}

func (n *node) PickElapsed() time.Duration {
	return time.Duration(time.Now().UnixNano() - atomic.LoadInt64(&n.lastPick))
}
