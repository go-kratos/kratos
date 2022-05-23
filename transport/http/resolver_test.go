package http

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"testing"
	"time"

	"github.com/go-kratos/kratos/v2/registry"
	"github.com/go-kratos/kratos/v2/selector"
)

func TestParseTarget(t *testing.T) {
	target, err := parseTarget("localhost:8000", true)
	if err != nil {
		t.Errorf("expect %v, got %v", nil, err)
	}
	if !reflect.DeepEqual(&Target{Scheme: "http", Authority: "localhost:8000"}, target) {
		t.Errorf("expect %v, got %v", &Target{Scheme: "http", Authority: "localhost:8000"}, target)
	}

	target, err = parseTarget("discovery:///demo", true)
	if err != nil {
		t.Errorf("expect %v, got %v", nil, err)
	}
	if !reflect.DeepEqual(&Target{Scheme: "discovery", Authority: "", Endpoint: "demo"}, target) {
		t.Errorf("expect %v, got %v", &Target{Scheme: "discovery", Authority: "", Endpoint: "demo"}, target)
	}

	target, err = parseTarget("127.0.0.1:8000", true)
	if err != nil {
		t.Errorf("expect %v, got %v", nil, err)
	}
	if !reflect.DeepEqual(&Target{Scheme: "http", Authority: "127.0.0.1:8000"}, target) {
		t.Errorf("expect %v, got %v", &Target{Scheme: "http", Authority: "127.0.0.1:8000"}, target)
	}

	target, err = parseTarget("https://127.0.0.1:8000", false)
	if err != nil {
		t.Errorf("expect %v, got %v", nil, err)
	}
	if !reflect.DeepEqual(&Target{Scheme: "https", Authority: "127.0.0.1:8000"}, target) {
		t.Errorf("expect %v, got %v", &Target{Scheme: "https", Authority: "127.0.0.1:8000"}, target)
	}

	target, err = parseTarget("127.0.0.1:8000", false)
	if err != nil {
		t.Errorf("expect %v, got %v", nil, err)
	}
	if !reflect.DeepEqual(&Target{Scheme: "https", Authority: "127.0.0.1:8000"}, target) {
		t.Errorf("expect %v, got %v", &Target{Scheme: "https", Authority: "127.0.0.1:8000"}, target)
	}
}

type mockRebalancer struct{}

func (m *mockRebalancer) Apply(nodes []selector.Node) {}

type mockDiscoveries struct {
	isSecure bool
}

func (d *mockDiscoveries) GetService(ctx context.Context, serviceName string) ([]*registry.ServiceInstance, error) {
	return nil, nil
}

func (d *mockDiscoveries) Watch(ctx context.Context, serviceName string) (registry.Watcher, error) {
	return &mockWatch{isSecure: d.isSecure}, nil
}

type mockWatch struct {
	isSecure bool
	count    int
}

func (m *mockWatch) Next() ([]*registry.ServiceInstance, error) {
	if m.count == 1 {
		return nil, errors.New("mock test error")
	}
	m.count++
	instance := &registry.ServiceInstance{
		ID:        "1",
		Name:      "kratos",
		Version:   "v1",
		Metadata:  map[string]string{},
		Endpoints: []string{fmt.Sprintf("http://127.0.0.1:9001?isSecure=%s", strconv.FormatBool(m.isSecure))},
	}
	if m.count > 3 {
		time.Sleep(time.Millisecond * 500)
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
		Endpoint:  "discovery://helloworld",
	}
	_, err := newResolver(context.Background(), &mockDiscoveries{true}, ta, &mockRebalancer{}, false, false)
	if err != nil {
		t.Errorf("expect %v, got %v", nil, err)
	}
	_, err = newResolver(context.Background(), &mockDiscoveries{false}, ta, &mockRebalancer{}, true, true)
	if err != nil {
		t.Errorf("expect %v, got %v", nil, err)
	}
}
