package wrr

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/go-kratos/kratos/v3/registry"
	"github.com/go-kratos/kratos/v3/selector"
	"github.com/go-kratos/kratos/v3/selector/filter"
)

func TestWrr(t *testing.T) {
	wrr := New()
	nodes := make([]selector.Node, 0, 2)
	nodes = append(nodes, selector.NewNode(
		"http",
		"127.0.0.1:8080",
		&registry.ServiceInstance{
			ID:       "127.0.0.1:8080",
			Version:  "v2.0.0",
			Metadata: map[string]string{"weight": "10"},
		}))
	nodes = append(nodes, selector.NewNode(
		"http",
		"127.0.0.1:9090",
		&registry.ServiceInstance{
			ID:       "127.0.0.1:9090",
			Version:  "v2.0.0",
			Metadata: map[string]string{"weight": "20"},
		}))
	wrr.Apply(nodes)
	var count1, count2 int
	for i := 0; i < 90; i++ {
		n, done, err := wrr.Select(context.Background(), selector.WithNodeFilter(filter.Version("v2.0.0")))
		if err != nil {
			t.Errorf("expect no error, got %v", err)
		}
		if done == nil {
			t.Errorf("expect done callback, got nil")
		}
		if n == nil {
			t.Errorf("expect node, got nil")
		}
		done(context.Background(), selector.DoneInfo{})
		if n.Address() == "127.0.0.1:8080" {
			count1++
		} else if n.Address() == "127.0.0.1:9090" {
			count2++
		}
	}
	if !reflect.DeepEqual(count1, 30) {
		t.Errorf("expect 30, got %d", count1)
	}
	if !reflect.DeepEqual(count2, 60) {
		t.Errorf("expect 60, got %d", count2)
	}
}

// TestCurrentWeightCleanup tests that stale entries in currentWeight map are cleaned up
func TestCurrentWeightCleanup(t *testing.T) {
	balancer := &Balancer{currentWeight: make(map[string]float64)}

	// Create initial nodes
	nodes1 := []selector.WeightedNode{
		&mockWeightedNode{address: "node1", weight: 10},
		&mockWeightedNode{address: "node2", weight: 20},
		&mockWeightedNode{address: "node3", weight: 30},
	}

	// Pick from initial nodes to populate currentWeight
	for i := 0; i < 10; i++ {
		_, _, err := balancer.Pick(context.Background(), nodes1)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	}

	// Verify all 3 nodes are in currentWeight map
	if len(balancer.currentWeight) != 3 {
		t.Errorf("expected 3 entries in currentWeight, got %d", len(balancer.currentWeight))
	}

	// Change to different set of nodes (simulating service discovery update)
	nodes2 := []selector.WeightedNode{
		&mockWeightedNode{address: "node2", weight: 20}, // only node2 remains
		&mockWeightedNode{address: "node4", weight: 40}, // node4 is new
	}

	// Pick from new nodes
	_, _, err := balancer.Pick(context.Background(), nodes2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify that stale entries (node1, node3) are cleaned up
	if len(balancer.currentWeight) != 2 {
		t.Errorf("expected 2 entries in currentWeight after cleanup, got %d", len(balancer.currentWeight))
	}

	// Verify that node2 and node4 are present, but node1 and node3 are not
	if _, exists := balancer.currentWeight["node1"]; exists {
		t.Error("stale entry node1 should have been cleaned up")
	}
	if _, exists := balancer.currentWeight["node3"]; exists {
		t.Error("stale entry node3 should have been cleaned up")
	}
	if _, exists := balancer.currentWeight["node2"]; !exists {
		t.Error("node2 should be present in currentWeight")
	}
	if _, exists := balancer.currentWeight["node4"]; !exists {
		t.Error("node4 should be present in currentWeight")
	}
}

// addrSet returns the set of addresses currently tracked in currentWeight.
func addrSet(b *Balancer) map[string]struct{} {
	b.mu.Lock()
	defer b.mu.Unlock()
	s := make(map[string]struct{}, len(b.currentWeight))
	for addr := range b.currentWeight {
		s[addr] = struct{}{}
	}
	return s
}

// assertTracked checks that currentWeight tracks exactly the given addresses,
// i.e. no stale entries linger and the map never grows beyond the live node set.
func assertTracked(t *testing.T, b *Balancer, want ...string) {
	t.Helper()
	got := addrSet(b)
	if len(got) != len(want) {
		t.Fatalf("currentWeight tracks %d entries, want %d (%v)", len(got), len(want), got)
	}
	for _, addr := range want {
		if _, ok := got[addr]; !ok {
			t.Fatalf("expected %q to be tracked, currentWeight=%v", addr, got)
		}
	}
}

// TestCleanupOnlyWhenNodesChange verifies that repeated picks with the same node
// set never accumulate stale entries, and that changing the node set (including a
// same-length swap and a partial overlap) drops entries for nodes that disappeared.
func TestCleanupOnlyWhenNodesChange(t *testing.T) {
	ctx := context.Background()
	balancer := &Balancer{currentWeight: make(map[string]float64)}

	nodes1 := []selector.WeightedNode{
		&mockWeightedNode{address: "node1", weight: 10},
		&mockWeightedNode{address: "node2", weight: 20},
	}
	// nodes2 is a full swap of nodes1 with the same length - the trickiest case
	// for length-based staleness detection.
	nodes2 := []selector.WeightedNode{
		&mockWeightedNode{address: "node3", weight: 30},
		&mockWeightedNode{address: "node4", weight: 40},
	}
	// nodes3 keeps node3 but replaces node4 with node5 (same length, partial overlap).
	nodes3 := []selector.WeightedNode{
		&mockWeightedNode{address: "node3", weight: 30},
		&mockWeightedNode{address: "node5", weight: 50},
	}

	pick := func(nodes []selector.WeightedNode) {
		t.Helper()
		if _, _, err := balancer.Pick(ctx, nodes); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	}

	// Repeated picks with the same nodes must not grow the map.
	for i := 0; i < 6; i++ {
		pick(nodes1)
		assertTracked(t, balancer, "node1", "node2")
	}

	// Full swap of an equal-length node set must drop the old entries.
	for i := 0; i < 4; i++ {
		pick(nodes2)
		assertTracked(t, balancer, "node3", "node4")
	}

	// Partial overlap swap: node4 leaves, node5 joins, node3 stays.
	for i := 0; i < 4; i++ {
		pick(nodes3)
		assertTracked(t, balancer, "node3", "node5")
	}
}

// mockWeightedNode is a mock implementation for testing
type mockWeightedNode struct {
	address string
	weight  float64
}

func (m *mockWeightedNode) Raw() selector.Node { return nil }
func (m *mockWeightedNode) Weight() float64    { return m.weight }
func (m *mockWeightedNode) Address() string    { return m.address }
func (m *mockWeightedNode) Pick() selector.DoneFunc {
	return func(context.Context, selector.DoneInfo) {}
}
func (m *mockWeightedNode) PickElapsed() time.Duration  { return 0 }
func (m *mockWeightedNode) Scheme() string              { return "http" }
func (m *mockWeightedNode) ServiceName() string         { return "test" }
func (m *mockWeightedNode) InitialWeight() *int64       { return nil }
func (m *mockWeightedNode) Version() string             { return "v1.0.0" }
func (m *mockWeightedNode) Metadata() map[string]string { return nil }

func TestEmpty(t *testing.T) {
	b := &Balancer{}
	_, _, err := b.Pick(context.Background(), []selector.WeightedNode{})
	if err == nil {
		t.Errorf("expect no error, got %v", err)
	}
}

// BenchmarkPickWithSameNodes benchmarks Pick() calls with the same node set
// This demonstrates the performance improvement where cleanup only happens on node changes
func BenchmarkPickWithSameNodes(b *testing.B) {
	balancer := &Balancer{currentWeight: make(map[string]float64)}

	// Create a fixed set of nodes
	nodes := []selector.WeightedNode{
		&mockWeightedNode{address: "node1", weight: 10},
		&mockWeightedNode{address: "node2", weight: 20},
		&mockWeightedNode{address: "node3", weight: 30},
		&mockWeightedNode{address: "node4", weight: 40},
		&mockWeightedNode{address: "node5", weight: 50},
	}

	ctx := context.Background()

	b.ResetTimer()
	b.ReportAllocs()

	// Benchmark Pick() calls with the same nodes
	// After the first call, no cleanup should occur on subsequent calls
	for i := 0; i < b.N; i++ {
		_, _, err := balancer.Pick(ctx, nodes)
		if err != nil {
			b.Fatalf("unexpected error: %v", err)
		}
	}
}

// BenchmarkPickWithChangingNodes benchmarks Pick() calls with changing node sets
// This shows the overhead when nodes actually change (expected to be slower)
func BenchmarkPickWithChangingNodes(b *testing.B) {
	balancer := &Balancer{currentWeight: make(map[string]float64)}

	// Create alternating sets of nodes to simulate node changes
	nodes1 := []selector.WeightedNode{
		&mockWeightedNode{address: "node1", weight: 10},
		&mockWeightedNode{address: "node2", weight: 20},
	}

	nodes2 := []selector.WeightedNode{
		&mockWeightedNode{address: "node3", weight: 30},
		&mockWeightedNode{address: "node4", weight: 40},
	}

	ctx := context.Background()

	b.ResetTimer()
	b.ReportAllocs()

	// Benchmark Pick() calls with alternating node sets
	// This will trigger cleanup on every call
	for i := 0; i < b.N; i++ {
		var nodes []selector.WeightedNode
		if i%2 == 0 {
			nodes = nodes1
		} else {
			nodes = nodes2
		}

		_, _, err := balancer.Pick(ctx, nodes)
		if err != nil {
			b.Fatalf("unexpected error: %v", err)
		}
	}
}
