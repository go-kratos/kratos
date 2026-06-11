package wrr

import (
	"context"
	"fmt"
	"math/rand/v2"
	"testing"

	"github.com/go-kratos/kratos/v3/selector"
)

// TestNoStaleEntriesUnderRandomChurn is a differential/property check: after any
// sequence of random node-set changes, currentWeight must track exactly the live
// node addresses and nothing else. This validates that the O(1) length-based
// staleness detection is equivalent to a full per-pick set comparison.
func TestNoStaleEntriesUnderRandomChurn(t *testing.T) {
	balancer := &Balancer{currentWeight: make(map[string]float64)}
	ctx := context.Background()
	rng := rand.New(rand.NewPCG(1, 2))

	makeNode := func(i int) selector.WeightedNode {
		return &mockWeightedNode{address: fmt.Sprintf("node-%d", i), weight: float64(1 + i%10)}
	}

	for iter := 0; iter < 20000; iter++ {
		// Build a random live set of 1..12 distinct nodes drawn from a pool of 20.
		size := 1 + rng.IntN(12)
		perm := rng.Perm(20)
		live := make(map[string]struct{}, size)
		nodes := make([]selector.WeightedNode, 0, size)
		for _, idx := range perm[:size] {
			n := makeNode(idx)
			nodes = append(nodes, n)
			live[n.Address()] = struct{}{}
		}

		// A few picks against this set (steady state between changes).
		for p := 0; p < 1+rng.IntN(3); p++ {
			if _, _, err := balancer.Pick(ctx, nodes); err != nil {
				t.Fatalf("iter %d: unexpected error: %v", iter, err)
			}
		}

		// currentWeight must equal the live set exactly: no stale, no missing.
		if len(balancer.currentWeight) != len(live) {
			t.Fatalf("iter %d: currentWeight has %d entries, want %d (%v)",
				iter, len(balancer.currentWeight), len(live), keysOf(balancer.currentWeight))
		}
		for addr := range balancer.currentWeight {
			if _, ok := live[addr]; !ok {
				t.Fatalf("iter %d: stale entry %q lingered in currentWeight", iter, addr)
			}
		}
	}
}

func keysOf(m map[string]float64) []string {
	ks := make([]string, 0, len(m))
	for k := range m {
		ks = append(ks, k)
	}
	return ks
}
