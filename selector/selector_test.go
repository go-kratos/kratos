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

var errNodeNotMatch = errors.New("node is not match")

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

func mockFilter(version string) NodeFilter {
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

type mockMustErrorBalancerBuilder struct{}

func (b *mockMustErrorBalancerBuilder) Build() Balancer {
	return &mockMustErrorBalancer{}
}

type mockMustErrorBalancer struct{}

func (b *mockMustErrorBalancer) Pick(ctx context.Context, nodes []WeightedNode) (selected WeightedNode, done DoneFunc, err error) {
	return nil, nil, errNodeNotMatch
}

func TestDefault(t *testing.T) {
	builder := DefaultBuilder{
		Node:     &mockWeightedNodeBuilder{},
		Balancer: &mockBalancerBuilder{},
	}
	selector := builder.Build()
	var nodes []Node
	nodes = append(nodes, NewNode(
		"http",
		"127.0.0.1:8080",
		&registry.ServiceInstance{
			ID:        "127.0.0.1:8080",
			Name:      "helloworld",
			Version:   "v2.0.0",
			Endpoints: []string{"http://127.0.0.1:8080"},
			Metadata:  map[string]string{"weight": "10"},
		}))
	nodes = append(nodes, NewNode(
		"http",
		"127.0.0.1:9090",
		&registry.ServiceInstance{
			ID:        "127.0.0.1:9090",
			Name:      "helloworld",
			Version:   "v1.0.0",
			Endpoints: []string{"http://127.0.0.1:9090"},
			Metadata:  map[string]string{"weight": "10"},
		}))

	selector.Apply(nodes)
	n, done, err := selector.Select(context.Background(), WithNodeFilter(mockFilter("v2.0.0")))
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
	if n.Scheme() == "" {
		t.Errorf("expect %v, got %v", "", n.Scheme())
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

	// peer in ctx
	ctx := NewPeerContext(context.Background(), &Peer{
		Node: mockWeightedNode{},
	})
	n, done, err = selector.Select(ctx)
	if err != nil {
		t.Errorf("expect %v, got %v", ErrNoAvailable, err)
	}
	if done == nil {
		t.Errorf("expect %v, got %v", nil, done)
	}
	if n == nil {
		t.Errorf("expect %v, got %v", nil, n)
	}

	// no v3.0.0 instance
	n, done, err = selector.Select(context.Background(), WithNodeFilter(mockFilter("v3.0.0")))
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
	n, done, err = selector.Select(context.Background(), WithNodeFilter(mockFilter("v2.0.0")))
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
	n, done, err = selector.Select(context.Background(), WithNodeFilter(mockFilter("v2.0.0")))
	if !errors.Is(ErrNoAvailable, err) {
		t.Errorf("expect %v, got %v", ErrNoAvailable, err)
	}
	if done != nil {
		t.Errorf("expect %v, got %v", nil, done)
	}
	if n != nil {
		t.Errorf("expect %v, got %v", nil, n)
	}

	// without node_filters
	n, done, err = selector.Select(context.Background())
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

func TestWithoutApply(t *testing.T) {
	builder := DefaultBuilder{
		Node:     &mockWeightedNodeBuilder{},
		Balancer: &mockBalancerBuilder{},
	}
	selector := builder.Build()
	n, done, err := selector.Select(context.Background())
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

func TestNoPick(t *testing.T) {
	builder := DefaultBuilder{
		Node:     &mockWeightedNodeBuilder{},
		Balancer: &mockMustErrorBalancerBuilder{},
	}
	var nodes []Node
	nodes = append(nodes, NewNode(
		"http",
		"127.0.0.1:8080",
		&registry.ServiceInstance{
			ID:        "127.0.0.1:8080",
			Name:      "helloworld",
			Version:   "v2.0.0",
			Endpoints: []string{"http://127.0.0.1:8080"},
			Metadata:  map[string]string{"weight": "10"},
		}))
	nodes = append(nodes, NewNode(
		"http",
		"127.0.0.1:9090",
		&registry.ServiceInstance{
			ID:        "127.0.0.1:9090",
			Name:      "helloworld",
			Version:   "v1.0.0",
			Endpoints: []string{"http://127.0.0.1:9090"},
			Metadata:  map[string]string{"weight": "10"},
		}))
	selector := builder.Build()
	selector.Apply(nodes)
	n, done, err := selector.Select(context.Background())
	if !errors.Is(errNodeNotMatch, err) {
		t.Errorf("expect %v, got %v", errNodeNotMatch, err)
	}
	if done != nil {
		t.Errorf("expect %v, got %v", nil, done)
	}
	if n != nil {
		t.Errorf("expect %v, got %v", nil, n)
	}
}

func TestGlobalSelector(t *testing.T) {
	builder := DefaultBuilder{
		Node:     &mockWeightedNodeBuilder{},
		Balancer: &mockBalancerBuilder{},
	}
	SetGlobalSelector(&builder)

	gBuilder := GlobalSelector()
	if gBuilder == nil {
		t.Errorf("expect %v, got %v", nil, gBuilder)
	}
}
