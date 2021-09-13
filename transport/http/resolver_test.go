package http

import (
	"context"
	"errors"
	"testing"

	"github.com/go-kratos/kratos/v2/registry"
	"github.com/go-kratos/kratos/v2/selector"
	"github.com/stretchr/testify/assert"
)

func TestParseTarget(t *testing.T) {
	target, err := parseTarget("localhost:8000", true)
	assert.Nil(t, err)
	assert.Equal(t, &Target{Scheme: "http", Authority: "localhost:8000"}, target)

	target, err = parseTarget("discovery:///demo", true)
	assert.Nil(t, err)
	assert.Equal(t, &Target{Scheme: "discovery", Authority: "", Endpoint: "demo"}, target)

	target, err = parseTarget("127.0.0.1:8000", true)
	assert.Nil(t, err)
	assert.Equal(t, &Target{Scheme: "http", Authority: "127.0.0.1:8000"}, target)

	target, err = parseTarget("https://127.0.0.1:8000", false)
	assert.Nil(t, err)
	assert.Equal(t, &Target{Scheme: "https", Authority: "127.0.0.1:8000"}, target)

	target, err = parseTarget("127.0.0.1:8000", false)
	assert.Nil(t, err)
	assert.Equal(t, &Target{Scheme: "https", Authority: "127.0.0.1:8000"}, target)
}

type mockRebalancer struct{}

func (m *mockRebalancer) Apply(nodes []selector.Node) {
	return
}

type mockDiscoverys struct{}

func (*mockDiscoverys) GetService(ctx context.Context, serviceName string) ([]*registry.ServiceInstance, error) {
	return nil, nil
}

func (*mockDiscoverys) Watch(ctx context.Context, serviceName string) (registry.Watcher, error) {
	return &mockWatch{}, nil
}

type mockWatch struct {
	count int
}

func (m *mockWatch) Next() ([]*registry.ServiceInstance, error) {
	if m.count == 0 {
		m.count++
		return nil, errors.New("mock test error")
	}
	instance := &registry.ServiceInstance{
		ID:        "1",
		Name:      "kratos",
		Version:   "v1",
		Metadata:  map[string]string{},
		Endpoints: []string{},
	}
	return []*registry.ServiceInstance{instance}, nil
}

func (m *mockWatch) Stop() error {
	return nil
}

func TestResolver(t *testing.T) {
	ta := &Target{
		Scheme:    "http",
		Authority: "",
		Endpoint:  "http://127.0.0.1:9001",
	}
	_, err := newResolver(context.Background(), &mockDiscoverys{}, ta, &mockRebalancer{}, false, false)
	assert.Nil(t, err)
	_, err = newResolver(context.Background(), &mockDiscoverys{}, ta, &mockRebalancer{}, true, false)
	assert.Nil(t, err)
}
