package http

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/go-kratos/kratos/v3/registry"
	"github.com/go-kratos/kratos/v3/selector"
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

func (m *mockRebalancer) Apply(_ []selector.Node) {}

type mockDiscoveries struct {
	isSecure bool
	nextErr  bool
	stopErr  bool
}

func (d *mockDiscoveries) GetService(_ context.Context, _ string) ([]*registry.ServiceInstance, error) {
	return nil, nil
}

const errServiceName = "needErr"

func (d *mockDiscoveries) Watch(ctx context.Context, serviceName string) (registry.Watcher, error) {
	if serviceName == errServiceName {
		return nil, errors.New("mock test service name watch err")
	}
	return &mockWatch{ctx: ctx, isSecure: d.isSecure, nextErr: d.nextErr, stopErr: d.stopErr}, nil
}

type mockWatch struct {
	ctx context.Context

	isSecure bool
	count    int

	nextErr bool
	stopErr bool

	lock sync.Mutex
}

func (m *mockWatch) Next() ([]*registry.ServiceInstance, error) {
	select {
	case <-m.ctx.Done():
		return nil, m.ctx.Err()
	default:
	}
	m.lock.Lock()
	defer m.lock.Unlock()
	if m.nextErr {
		return nil, errors.New("mock test error")
	}
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
	m.lock.Lock()
	defer m.lock.Unlock()
	if m.stopErr {
		return errors.New("mock test error")
	}
	// Mark the next call to return an error.
	m.nextErr = true
	return nil
}

func TestResolver(t *testing.T) {
	ta, err := parseTarget("discovery://helloworld", true)
	if err != nil {
		t.Errorf("parse err %v", err)
		return
	}

	cancelCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Async mode should not return an error.
	r, err := newResolver(cancelCtx, &mockDiscoveries{true, false, false}, ta, &mockRebalancer{}, false, false, 25)
	if err != nil {
		t.Errorf("expect %v, got %v", nil, err)
	}
	if r != nil {
		_ = r.Close()
	}

	// Sync mode should run successfully.
	r, err = newResolver(cancelCtx, &mockDiscoveries{false, false, false}, ta, &mockRebalancer{}, true, true, 25)
	if err != nil {
		t.Errorf("expect %v, got %v", nil, err)
	}
	if r != nil {
		_ = r.Close()
	}

	// Sync mode should return an error when Next and Stop fail.
	r, err = newResolver(cancelCtx, &mockDiscoveries{false, true, true}, ta, &mockRebalancer{}, true, true, 25)
	if err == nil {
		t.Errorf("expect err, got nil")
	}
	if r != nil {
		_ = r.Close()
	}

	// Sync mode should return an error when service name watch fails.
	r, err = newResolver(cancelCtx, &mockDiscoveries{false, true, true}, &Target{
		Scheme:   "discovery",
		Endpoint: errServiceName,
	}, &mockRebalancer{}, true, true, 25)
	if err == nil {
		t.Errorf("expect err, got nil")
	}
	if r != nil {
		_ = r.Close()
	}

	cancel()

	// This should log context.Canceled.
	r, err = newResolver(cancelCtx, &mockDiscoveries{false, false, false}, ta, &mockRebalancer{}, false, false, 25)
	if err != nil {
		t.Errorf("expect %v, got %v", nil, err)
	}
	if r != nil {
		_ = r.Close()
	}

	// Sync mode should return an error when the service is canceled.
	r, err = newResolver(cancelCtx, &mockDiscoveries{false, false, true}, ta, &mockRebalancer{}, true, true, 25)
	if err == nil {
		t.Errorf("expect ctx cancel err, got nil")
	}
	if r != nil {
		_ = r.Close()
	}

	time.Sleep(100 * time.Millisecond)
}
