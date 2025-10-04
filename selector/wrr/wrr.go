package wrr

import (
	"context"
	"sync"

	"github.com/go-kratos/kratos/v2/selector"
	"github.com/go-kratos/kratos/v2/selector/node/direct"
)

const (
	// Name is wrr(Weighted Round Robin) balancer name
	Name = "wrr"
)

var _ selector.Balancer = (*Balancer)(nil) // Name is balancer name

// Option is wrr builder option.
type Option func(o *options)

// options is wrr builder options
type options struct{}

// Balancer is a wrr balancer.
type Balancer struct {
	mu            sync.Mutex
	currentWeight map[string]float64
	lastNodes     []selector.WeightedNode
}

// equalNodes checks if two slices of WeightedNode contain the same nodes
func equalNodes(a, b []selector.WeightedNode) bool {
	if len(a) != len(b) {
		return false
	}

	// Create a map of addresses from slice a
	aMap := make(map[string]bool)
	for _, node := range a {
		aMap[node.Address()] = true
	}

	// Check if all nodes in slice b exist in slice a
	for _, node := range b {
		if !aMap[node.Address()] {
			return false
		}
	}

	return true
}

// New random a selector.
func New(opts ...Option) selector.Selector {
	return NewBuilder(opts...).Build()
}

// Pick is pick a weighted node.
func (p *Balancer) Pick(_ context.Context, nodes []selector.WeightedNode) (selector.WeightedNode, selector.DoneFunc, error) {
	if len(nodes) == 0 {
		return nil, nil, selector.ErrNoAvailable
	}

	p.mu.Lock()
	defer p.mu.Unlock()

	// Check if the node list has changed
	if len(p.lastNodes) != len(nodes) || !equalNodes(p.lastNodes, nodes) {
		// Update lastNodes
		p.lastNodes = make([]selector.WeightedNode, len(nodes))
		copy(p.lastNodes, nodes)

		// Create a set of current node addresses for cleanup
		currentNodes := make(map[string]bool)
		for _, node := range nodes {
			currentNodes[node.Address()] = true
		}

		// Clean up stale entries from currentWeight map
		for address := range p.currentWeight {
			if !currentNodes[address] {
				delete(p.currentWeight, address)
			}
		}
	}

	var totalWeight float64
	var selected selector.WeightedNode
	var selectWeight float64

	// nginx wrr load balancing algorithm: http://blog.csdn.net/zhangskd/article/details/50194069
	for _, node := range nodes {
		totalWeight += node.Weight()
		cwt := p.currentWeight[node.Address()]
		// current += effectiveWeight
		cwt += node.Weight()
		p.currentWeight[node.Address()] = cwt
		if selected == nil || selectWeight < cwt {
			selectWeight = cwt
			selected = node
		}
	}
	p.currentWeight[selected.Address()] = selectWeight - totalWeight

	d := selected.Pick()
	return selected, d, nil
}

// NewBuilder returns a selector builder with wrr balancer
func NewBuilder(opts ...Option) selector.Builder {
	var option options
	for _, opt := range opts {
		opt(&option)
	}
	return &selector.DefaultBuilder{
		Balancer: &Builder{},
		Node:     &direct.Builder{},
	}
}

// Builder is wrr builder
type Builder struct{}

// Build creates Balancer
func (b *Builder) Build() selector.Balancer {
	return &Balancer{currentWeight: make(map[string]float64)}
}
