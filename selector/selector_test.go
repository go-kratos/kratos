package selector

import (
	"context"
	"errors"
	"math/rand"
	"reflect"
	"sync/atomic"
	"testing"
	"time"

	"github.com/go-kratos/kratos/v2/registry"
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

func (b *mockWeightedNodeBuilder) Build(n Node) WeightedNode {
	return &mockWeightedNode{Node: n}
}

func mockFilter(version string) Filter {
	return func(_ context.Context, nodes []Node) []Node {
		newNodes := nodes[:0]
		for _, n := range nodes {
			if n.Version() == version {
				newNodes = append(newNodes, n)
			}
		}
		return newNodes
	}
}

type mockBalancerBuilder struct{}

func (b *mockBalancerBuilder) Build() Balancer {
	return &mockBalancer{}
}

type mockBalancer struct{}

func (b *mockBalancer) Pick(ctx context.Context, nodes []WeightedNode) (selected WeightedNode, done DoneFunc, err error) {
	if len(nodes) == 0 {
		err = ErrNoAvailable
		return
	}
	cur := rand.Intn(len(nodes))
	selected = nodes[cur]
	done = selected.Pick()
	return
}

func TestDefault(t *testing.T) {
	builder := DefaultBuilder{
		Node:     &mockWeightedNodeBuilder{},
		Filters:  []Filter{mockFilter("v2.0.0")},
		Balancer: &mockBalancerBuilder{},
	}
	selector := builder.Build()
	var nodes []Node
	nodes = append(nodes, NewNode(
		"127.0.0.1:8080",
		&registry.ServiceInstance{
			ID:        "127.0.0.1:8080",
			Name:      "helloworld",
			Version:   "v2.0.0",
			Endpoints: []string{"http://127.0.0.1:8080"},
			Metadata:  map[string]string{"weight": "10"},
		}))
	nodes = append(nodes, NewNode(
		"127.0.0.1:9090",
		&registry.ServiceInstance{
			ID:        "127.0.0.1:9090",
			Name:      "helloworld",
			Version:   "v1.0.0",
			Endpoints: []string{"http://127.0.0.1:9090"},
			Metadata:  map[string]string{"weight": "10"},
		}))
	selector.Apply(nodes)
	n, done, err := selector.Select(context.Background(), WithFilter(mockFilter("v2.0.0")))
	if err != nil {
		t.Errorf("expect %v, got %v", nil, err)
	}
	if n == nil {
		t.Errorf("expect %v, got %v", nil, n)
	}
	if done == nil {
		t.Errorf("expect %v, got %v", nil, done)
	}
	if !reflect.DeepEqual("v2.0.0", n.Version()) {
		t.Errorf("expect %v, got %v", "v2.0.0", n.Version())
	}
	if n.Address() == "" {
		t.Errorf("expect %v, got %v", "", n.Address())
	}
	if !reflect.DeepEqual(int64(10), *n.InitialWeight()) {
		t.Errorf("expect %v, got %v", 10, *n.InitialWeight())
	}
	if n.Metadata() == nil {
		t.Errorf("expect %v, got %v", nil, n.Metadata())
	}
	if !reflect.DeepEqual("helloworld", n.ServiceName()) {
		t.Errorf("expect %v, got %v", "helloworld", n.ServiceName())
	}
	done(context.Background(), DoneInfo{})

	// no v3.0.0 instance
	n, done, err = selector.Select(context.Background(), WithFilter(mockFilter("v3.0.0")))
	if !errors.Is(ErrNoAvailable, err) {
		t.Errorf("expect %v, got %v", ErrNoAvailable, err)
	}
	if done != nil {
		t.Errorf("expect %v, got %v", nil, done)
	}
	if n != nil {
		t.Errorf("expect %v, got %v", nil, n)
	}

	// apply zero instance
	selector.Apply([]Node{})
	n, done, err = selector.Select(context.Background(), WithFilter(mockFilter("v2.0.0")))
	if !errors.Is(ErrNoAvailable, err) {
		t.Errorf("expect %v, got %v", ErrNoAvailable, err)
	}
	if done != nil {
		t.Errorf("expect %v, got %v", nil, done)
	}
	if n != nil {
		t.Errorf("expect %v, got %v", nil, n)
	}

	// apply zero instance
	selector.Apply(nil)
	n, done, err = selector.Select(context.Background(), WithFilter(mockFilter("v2.0.0")))
	if !errors.Is(ErrNoAvailable, err) {
		t.Errorf("expect %v, got %v", ErrNoAvailable, err)
	}
	if done != nil {
		t.Errorf("expect %v, got %v", nil, done)
	}
	if n != nil {
		t.Errorf("expect %v, got %v", nil, n)
	}
}
