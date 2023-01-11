package selector

import (
	"context"
	"time"
)

// Balancer is balancer interface
type Balancer[W WeightedNode] interface {
	Pick(ctx context.Context, nodes []W) (selected W, done DoneFunc, err error)
}

// BalancerBuilder build balancer
type BalancerBuilder[W WeightedNode] interface {
	Build() Balancer[W]
}

// WeightedNode calculates scheduling weight in real time
type WeightedNode interface {
	Node

	// Raw returns the original node
	Raw() Node

	// Weight is the runtime calculated weight
	Weight() float64

	// Pick the node
	Pick() DoneFunc

	// PickElapsed is time elapsed since the latest pick
	PickElapsed() time.Duration
}

// WeightedNodeBuilder is WeightedNode Builder
type WeightedNodeBuilder[W WeightedNode] interface {
	Build(Node) W
}
