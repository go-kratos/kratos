package ewma

import (
	"context"
	"net"
	"reflect"
	"testing"
	"time"

	"github.com/go-kratos/kratos/v3/registry"
	"github.com/go-kratos/kratos/v3/selector"
)

func TestDirect(t *testing.T) {
	b := &Builder{}
	wn := b.Build(selector.NewNode(
		"http",
		"127.0.0.1:9090",
		&registry.ServiceInstance{
			ID:        "127.0.0.1:9090",
			Name:      "helloworld",
			Version:   "v1.0.0",
			Endpoints: []string{"http://127.0.0.1:9090"},
			Metadata:  map[string]string{"weight": "10"},
		}))

	if !reflect.DeepEqual(float64(100), wn.Weight()) {
		t.Errorf("expect %v, got %v", 100, wn.Weight())
	}
	done := wn.Pick()
	if done == nil {
		t.Errorf("done is equal to nil")
	}
	done2 := wn.Pick()
	if done2 == nil {
		t.Errorf("done2 is equal to nil")
	}

	time.Sleep(time.Millisecond * 15)
	done(context.Background(), selector.DoneInfo{})
	if float64(70) >= wn.Weight() {
		t.Errorf("float64(30000) >= wn.Weight()(%v)", wn.Weight())
	}
	if float64(1200) <= wn.Weight() {
		t.Errorf("float64(1000) <= wn.Weight()(%v)", wn.Weight())
	}
	if time.Millisecond*30 <= wn.PickElapsed() {
		t.Errorf("time.Millisecond*30 <= wn.PickElapsed()(%v)", wn.PickElapsed())
	}
	if time.Millisecond*5 >= wn.PickElapsed() {
		t.Errorf("time.Millisecond*5 >= wn.PickElapsed()(%v)", wn.PickElapsed())
	}
}

func TestDirectError(t *testing.T) {
	b := &Builder{}
	wn := b.Build(selector.NewNode(
		"http",
		"127.0.0.1:9090",
		&registry.ServiceInstance{
			ID:        "127.0.0.1:9090",
			Name:      "helloworld",
			Version:   "v1.0.0",
			Endpoints: []string{"http://127.0.0.1:9090"},
			Metadata:  map[string]string{"weight": "10"},
		}))

	for i := 0; i < 5; i++ {
		var err error
		if i != 0 {
			err = context.DeadlineExceeded
		}
		done := wn.Pick()
		if done == nil {
			t.Errorf("expect not nil, got nil")
		}
		time.Sleep(time.Millisecond * 20)
		done(context.Background(), selector.DoneInfo{Err: err})
	}
	if float64(1000) >= wn.Weight() {
		t.Errorf("float64(1000) >= wn.Weight()(%v)", wn.Weight())
	}
	if float64(2000) <= wn.Weight() {
		t.Errorf("float64(2000) <= wn.Weight()(%v)", wn.Weight())
	}
}

// TestCanceledDoesNotDegradeNode verifies that context.Canceled does not
// lower a node's health score. Canceled means the caller gave up — it is
// not evidence of backend failure. context.DeadlineExceeded (backend too
// slow) should still degrade the node as before.
func TestCanceledDoesNotDegradeNode(t *testing.T) {
	newNode := func() selector.WeightedNode {
		return (&Builder{}).Build(selector.NewNode(
			"http",
			"127.0.0.1:9090",
			&registry.ServiceInstance{
				ID:        "127.0.0.1:9090",
				Name:      "helloworld",
				Version:   "v1.0.0",
				Endpoints: []string{"http://127.0.0.1:9090"},
				Metadata:  map[string]string{"weight": "10"},
			}))
	}

	// --- context.Canceled must NOT degrade health ---
	wn := newNode()
	n := wn.(*Node)
	// One successful pick to initialise EWMA state.
	done := wn.Pick()
	time.Sleep(time.Millisecond * 20)
	done(context.Background(), selector.DoneInfo{})
	baselineHealth := n.success.Load()

	// Several picks that all report context.Canceled.
	for i := 0; i < 4; i++ {
		done = wn.Pick()
		time.Sleep(time.Millisecond * 20)
		done(context.Background(), selector.DoneInfo{Err: context.Canceled})
	}
	afterHealth := n.success.Load()
	// Health must not have dropped significantly — Canceled is not a backend fault.
	// Allow a tiny EWMA drift but it must stay near 1000 (full health).
	if afterHealth < baselineHealth*9/10 {
		t.Errorf("context.Canceled should not degrade node health: before=%d after=%d",
			baselineHealth, afterHealth)
	}

	// --- context.DeadlineExceeded still degrades health ---
	wn2 := newNode()
	n2 := wn2.(*Node)
	done = wn2.Pick()
	time.Sleep(time.Millisecond * 20)
	done(context.Background(), selector.DoneInfo{})
	baselineHealth2 := n2.success.Load()

	for i := 0; i < 4; i++ {
		done = wn2.Pick()
		time.Sleep(time.Millisecond * 20)
		done(context.Background(), selector.DoneInfo{Err: context.DeadlineExceeded})
	}
	afterHealth2 := n2.success.Load()
	if afterHealth2 >= baselineHealth2 {
		t.Errorf("context.DeadlineExceeded should degrade node health: before=%d after=%d",
			baselineHealth2, afterHealth2)
	}
}

func TestDirectErrorHandler(t *testing.T) {
	b := &Builder{
		ErrHandler: func(err error) bool {
			return err != nil
		},
	}
	wn := b.Build(selector.NewNode(
		"http",
		"127.0.0.1:9090",
		&registry.ServiceInstance{
			ID:        "127.0.0.1:9090",
			Name:      "helloworld",
			Version:   "v1.0.0",
			Endpoints: []string{"http://127.0.0.1:9090"},
			Metadata:  map[string]string{"weight": "10"},
		}))
	errs := []error{
		context.DeadlineExceeded,
		context.Canceled,
		net.ErrClosed,
	}
	for i := 0; i < 5; i++ {
		var err error
		if i != 0 {
			err = errs[i%len(errs)]
		}
		done := wn.Pick()
		if done == nil {
			t.Errorf("expect not nil, got nil")
		}
		time.Sleep(time.Millisecond * 20)
		done(context.Background(), selector.DoneInfo{Err: err})
	}
	if float64(1000) >= wn.Weight() {
		t.Errorf("float64(100) >= wn.Weight()(%v)", wn.Weight())
	}
	if float64(2000) <= wn.Weight() {
		t.Errorf("float64(200) <= wn.Weight()(%v)", wn.Weight())
	}
}

func BenchmarkPickAndWeight(b *testing.B) {
	bu := &Builder{}
	node := bu.Build(selector.NewNode(
		"http",
		"127.0.0.1:9090",
		&registry.ServiceInstance{
			ID:        "127.0.0.1:9090",
			Name:      "helloworld",
			Version:   "v1.0.0",
			Endpoints: []string{"http://127.0.0.1:9090"},
			Metadata:  map[string]string{"weight": "10"},
		}))
	di := selector.DoneInfo{}
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			done := node.Pick()
			node.Weight()
			done(context.Background(), di)
		}
	})
}
