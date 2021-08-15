package direct

import (
	"context"
	"sync/atomic"
	"time"

	"github.com/go-kratos/kratos/v2/balancer"
)

var _ balancer.Node = &node{}
var _ balancer.NodeBuilder = &Builder{}

// node is endpoint instance
type node struct {
	addr     string
	metadata balancer.Metadata
	weight   float64

	//last lastPick timestamp
	lastPick int64
}

// Builder is direct node builder
type Builder struct {
}

// Build create node
func (*Builder) Build(addr string, initWeight float64, metadata balancer.Metadata) balancer.Node {
	return &node{
		addr:     addr,
		metadata: metadata,
		weight:   initWeight,
	}
}

func (n *node) Pick() balancer.Done {
	now := time.Now().UnixNano()
	atomic.StoreInt64(&n.lastPick, now)

	return func(ctx context.Context, di balancer.DoneInfo) {}
}

// Weight is node effective weigth
func (n *node) Weight() (weight float64) {
	return n.weight
}

func (n *node) LastPick() time.Time {
	return time.Unix(0, atomic.LoadInt64(&n.lastPick))
}

func (n *node) Address() string {
	return n.Address()
}

func (n *node) Metadata() balancer.Metadata {
	return n.metadata
}
