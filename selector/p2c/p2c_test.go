package p2c

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/go-kratos/kratos/v2/registry"
	"github.com/go-kratos/kratos/v2/selector"
	"github.com/go-kratos/kratos/v2/selector/filter"
	"github.com/stretchr/testify/assert"
)

func TestWrr3(t *testing.T) {
	p2c := New(WithFilter(filter.Version("v2.0.0")))
	var nodes []selector.Node
	for i := 0; i < 3; i++ {
		addr := fmt.Sprintf("127.0.0.%d:8080", i)
		nodes = append(nodes, selector.NewNode(
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
			assert.Nil(t, err)
			assert.NotNil(t, done)
			assert.NotNil(t, n)
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
	assert.Greater(t, count1, int64(1500))
	assert.Less(t, count1, int64(4500))
	assert.Greater(t, count2, int64(1500))
	assert.Less(t, count2, int64(4500))
	assert.Greater(t, count3, int64(1500))
	assert.Less(t, count3, int64(4500))
}

func TestEmpty(t *testing.T) {
	b := &Balancer{}
	_, _, err := b.Pick(context.Background(), []selector.WeightedNode{})
	assert.NotNil(t, err)
}

func TestOne(t *testing.T) {
	p2c := New(WithFilter(filter.Version("v2.0.0")))
	var nodes []selector.Node
	for i := 0; i < 1; i++ {
		addr := fmt.Sprintf("127.0.0.%d:8080", i)
		nodes = append(nodes, selector.NewNode(
			addr,
			&registry.ServiceInstance{
				ID:       addr,
				Version:  "v2.0.0",
				Metadata: map[string]string{"weight": "10"},
			}))
	}
	p2c.Apply(nodes)
	n, done, err := p2c.Select(context.Background())
	assert.Nil(t, err)
	assert.NotNil(t, done)
	assert.NotNil(t, n)
	assert.Equal(t, "127.0.0.0:8080", n.Address())
}
