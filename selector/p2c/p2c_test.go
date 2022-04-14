package p2c

import (
	"context"
	"fmt"
	"math/rand"
	"reflect"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/go-kratos/kratos/v2/registry"
	"github.com/go-kratos/kratos/v2/selector"
	"github.com/go-kratos/kratos/v2/selector/filter"
)

func TestWrr3(t *testing.T) {
	p2c := New(WithFilter(filter.Version("v2.0.0")))
	var nodes []selector.Node
	for i := 0; i < 3; i++ {
		addr := fmt.Sprintf("127.0.0.%d:8080", i)
		nodes = append(nodes, selector.NewNode(
			"http",
			addr,
			&registry.ServiceInstance{
				ID:       addr,
				Version:  "v2.0.0",
				Metadata: map[string]string{"weight": "10"},
			}))
	}
	p2c.Apply(nodes)
	var count1, count2, count3 int64
	group := &sync.WaitGroup{}
	var lk sync.Mutex
	for i := 0; i < 9000; i++ {
		group.Add(1)
		go func() {
			defer group.Done()
			lk.Lock()
			d := time.Duration(rand.Intn(500)) * time.Millisecond
			lk.Unlock()
			time.Sleep(d)
			n, done, err := p2c.Select(context.Background())
			if err != nil {
				t.Errorf("expect %v, got %v", nil, err)
			}
			if n == nil {
				t.Errorf("expect %v, got %v", nil, n)
			}
			if done == nil {
				t.Errorf("expect %v, got %v", nil, done)
			}
			time.Sleep(time.Millisecond * 10)
			done(context.Background(), selector.DoneInfo{})
			if n.Address() == "127.0.0.0:8080" {
				atomic.AddInt64(&count1, 1)
			} else if n.Address() == "127.0.0.1:8080" {
				atomic.AddInt64(&count2, 1)
			} else if n.Address() == "127.0.0.2:8080" {
				atomic.AddInt64(&count3, 1)
			}
		}()
	}
	group.Wait()
	if count1 <= int64(1500) {
		t.Errorf("count1(%v) <= int64(1500)", count1)
	}
	if count1 >= int64(4500) {
		t.Errorf("count1(%v) >= int64(4500),", count1)
	}
	if count2 <= int64(1500) {
		t.Errorf("count2(%v) <= int64(1500)", count1)
	}
	if count2 >= int64(4500) {
		t.Errorf("count2(%v) >= int64(4500),", count2)
	}
	if count3 <= int64(1500) {
		t.Errorf("count3(%v) <= int64(1500)", count3)
	}
	if count3 >= int64(4500) {
		t.Errorf("count3(%v) >= int64(4500),", count3)
	}
}

func TestEmpty(t *testing.T) {
	b := &Balancer{}
	_, _, err := b.Pick(context.Background(), []selector.WeightedNode{})
	if err == nil {
		t.Errorf("expect %v, got %v", nil, err)
	}
}

func TestOne(t *testing.T) {
	p2c := New(WithFilter(filter.Version("v2.0.0")))
	var nodes []selector.Node
	for i := 0; i < 1; i++ {
		addr := fmt.Sprintf("127.0.0.%d:8080", i)
		nodes = append(nodes, selector.NewNode(
			"http",
			addr,
			&registry.ServiceInstance{
				ID:       addr,
				Version:  "v2.0.0",
				Metadata: map[string]string{"weight": "10"},
			}))
	}
	p2c.Apply(nodes)
	n, done, err := p2c.Select(context.Background())
	if err != nil {
		t.Errorf("expect %v, got %v", nil, err)
	}
	if n == nil {
		t.Errorf("expect %v, got %v", nil, n)
	}
	if done == nil {
		t.Errorf("expect %v, got %v", nil, done)
	}
	if !reflect.DeepEqual("127.0.0.0:8080", n.Address()) {
		t.Errorf("expect %v, got %v", "127.0.0.0:8080", n.Address())
	}
}
