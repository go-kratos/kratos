package random

import (
	"context"
	"testing"

	"github.com/go-kratos/kratos/v2/registry"
	"github.com/go-kratos/kratos/v2/selector"
	"github.com/go-kratos/kratos/v2/selector/filter"
)

func TestWrr(t *testing.T) {
	random := New()
	var nodes []selector.Node
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
	random.Apply(nodes)
	var count1, count2 int
	randomNumber := 2000
	for i := 0; i < randomNumber; i++ {
		n, done, err := random.Select(context.Background(), selector.WithNodeFilter(filter.Version("v2.0.0")))
		if err != nil {
			t.Errorf("expect no error, got %v", err)
		}
		if done == nil {
			t.Errorf("expect not nil, got:%v", done)
		}
		if n == nil {
			t.Errorf("expect not nil, got:%v", n)
		}
		done(context.Background(), selector.DoneInfo{})
		if n.Address() == "127.0.0.1:8080" {
			count1++
		} else if n.Address() == "127.0.0.1:9090" {
			count2++
		}
	}
	// output
	percentage1 := float64(count1) / float64(randomNumber) * 100
	percentage2 := float64(count2) / float64(randomNumber) * 100
	if percentage1 > 60 {
		t.Errorf("percentage1(%v) > 60", percentage1)
	}
	if percentage1 < 40 {
		t.Errorf("percentage1(%v) < 40", percentage1)
	}
	if percentage2 > 60 {
		t.Errorf("percentage2(%v) > 60", percentage2)
	}
	if percentage2 < 40 {
		t.Errorf("percentage2(%v) < 40", percentage2)
	}
}

func TestEmpty(t *testing.T) {
	b := &Balancer{}
	_, _, err := b.Pick(context.Background(), []selector.WeightedNode{})
	if err == nil {
		t.Errorf("expect nil, got %v", err)
	}
}
