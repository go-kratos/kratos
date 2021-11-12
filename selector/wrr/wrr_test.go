package wrr

import (
	"context"
	"testing"

	"github.com/go-kratos/kratos/v2/registry"
	"github.com/go-kratos/kratos/v2/selector"
	"github.com/go-kratos/kratos/v2/selector/filter"
	"github.com/stretchr/testify/assert"
)

func TestWrr(t *testing.T) {
	wrr := New(WithFilter(filter.Version("v2.0.0")))
	var nodes []selector.Node
	nodes = append(nodes, selector.NewNode(
		"127.0.0.1:8080",
		&registry.ServiceInstance{
			ID:       "127.0.0.1:8080",
			Version:  "v2.0.0",
			Metadata: map[string]string{"weight": "10"},
		}))
	nodes = append(nodes, selector.NewNode(
		"127.0.0.1:9090",
		&registry.ServiceInstance{
			ID:       "127.0.0.1:9090",
			Version:  "v2.0.0",
			Metadata: map[string]string{"weight": "20"},
		}))
	wrr.Apply(nodes)
	var count1, count2 int
	for i := 0; i < 90; i++ {
		n, done, err := wrr.Select(context.Background())
		assert.Nil(t, err)
		assert.NotNil(t, done)
		assert.NotNil(t, n)
		done(context.Background(), selector.DoneInfo{})
		if n.Address() == "127.0.0.1:8080" {
			count1++
		} else if n.Address() == "127.0.0.1:9090" {
			count2++
		}
	}
	assert.Equal(t, 30, count1)
	assert.Equal(t, 60, count2)
}

func TestEmpty(t *testing.T) {
	b := &Balancer{}
	_, _, err := b.Pick(context.Background(), []selector.WeightedNode{})
	assert.NotNil(t, err)
}
