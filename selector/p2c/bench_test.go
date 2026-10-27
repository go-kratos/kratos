package p2c

import (
	"context"
	"fmt"
	"testing"

	"github.com/go-kratos/kratos/v3/registry"
	"github.com/go-kratos/kratos/v3/selector"
)

func benchSelector(nodeCount int) selector.Selector {
	s := New()
	nodes := make([]selector.Node, 0, nodeCount)
	for i := 0; i < nodeCount; i++ {
		addr := fmt.Sprintf("127.0.0.%d:8080", i)
		nodes = append(nodes, selector.NewNode("http", addr, &registry.ServiceInstance{
			ID:       addr,
			Version:  "v1.0.0",
			Metadata: map[string]string{"weight": "10"},
		}))
	}
	s.Apply(nodes)
	return s
}

// BenchmarkSelectParallel hammers a single shared balancer from GOMAXPROCS
// goroutines, which is how a balancer is used under real client load. The
// per-pick mutex around the *rand.Rand serializes every pick here.
func BenchmarkSelectParallel(b *testing.B) {
	s := benchSelector(10)
	ctx := context.Background()
	b.ResetTimer()
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, done, err := s.Select(ctx)
			if err != nil {
				b.Fatal(err)
			}
			done(ctx, selector.DoneInfo{})
		}
	})
}

// BenchmarkSelectSerial measures the single-goroutine cost (uncontended lock).
func BenchmarkSelectSerial(b *testing.B) {
	s := benchSelector(10)
	ctx := context.Background()
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, done, err := s.Select(ctx)
		if err != nil {
			b.Fatal(err)
		}
		done(ctx, selector.DoneInfo{})
	}
}
