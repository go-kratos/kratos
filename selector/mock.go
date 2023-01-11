package selector

import (
	"context"
	"sync/atomic"
	"time"
)

type mockWeightedNode struct {
	Node

	lastPick int64
}

// Raw returns the original node
func (n *mockWeightedNode) Raw() Node {
	return n.Node
}

// Weight is the runtime calculated weight
func (n *mockWeightedNode) Weight() float64 {
	if n.InitialWeight() != nil {
		return float64(*n.InitialWeight())
	}
	return 100
}

// Pick the node
func (n *mockWeightedNode) Pick() DoneFunc {
	now := time.Now().UnixNano()
	atomic.StoreInt64(&n.lastPick, now)
	return func(ctx context.Context, di DoneInfo) {}
}

// PickElapsed is time elapsed since the latest pick
func (n *mockWeightedNode) PickElapsed() time.Duration {
	return time.Duration(time.Now().UnixNano() - atomic.LoadInt64(&n.lastPick))
}

type mockWeightedNodeBuilder struct{}

func (b *mockWeightedNodeBuilder) Build(n Node) *mockWeightedNode {
	return &mockWeightedNode{Node: n}
}
