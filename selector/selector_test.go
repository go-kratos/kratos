package selector

import (
	"context"
	"errors"
	"math/rand"
	"reflect"
	"testing"

	"github.com/go-kratos/kratos/v2/registry"
)

var errNodeNotMatch = errors.New("node is not match")

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

type mockBalancerBuilder[W WeightedNode] struct {
}

func (b *mockBalancerBuilder[W]) Build() Balancer[W] {
	return &mockBalancer[W]{}
}

type mockBalancer[W WeightedNode] struct{}

func (b *mockBalancer[W]) Pick(ctx context.Context, nodes []W) (selected W, done DoneFunc, err error) {
	if len(nodes) == 0 {
		err = ErrNoAvailable
		return
	}
	cur := rand.Intn(len(nodes))
	selected = nodes[cur]
	done = selected.Pick()
	return
}

type mockMustErrorBalancerBuilder[W WeightedNode] struct{}

func (b *mockMustErrorBalancerBuilder[W]) Build() Balancer[W] {
	return &mockMustErrorBalancer[W]{}
}

type mockMustErrorBalancer[W WeightedNode] struct{}

func (b *mockMustErrorBalancer[W]) Pick(ctx context.Context, nodes []W) (selected W, done DoneFunc, err error) {
	var zero W
	return zero, nil, errNodeNotMatch
}

func TestDefault(t *testing.T) {
	builder := DefaultBuilder[*mockWeightedNode]{
		Node:     &mockWeightedNodeBuilder{},
		Balancer: &mockBalancerBuilder[*mockWeightedNode]{},
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
	builder := DefaultBuilder[*mockWeightedNode]{
		Node:     &mockWeightedNodeBuilder{},
		Balancer: &mockBalancerBuilder[*mockWeightedNode]{},
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
	builder := DefaultBuilder[*mockWeightedNode]{
		Node:     &mockWeightedNodeBuilder{},
		Balancer: &mockMustErrorBalancerBuilder[*mockWeightedNode]{},
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
	builder := DefaultBuilder[*mockWeightedNode]{
		Node:     &mockWeightedNodeBuilder{},
		Balancer: &mockBalancerBuilder[*mockWeightedNode]{},
	}
	SetGlobalSelector(&builder)

	gBuilder := GlobalSelector()
	if gBuilder == nil {
		t.Errorf("expect %v, got %v", nil, gBuilder)
	}
}
