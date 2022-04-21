package registry_test

import (
	"context"
	"github.com/go-kratos/kratos/v2/registry"
	"sync/atomic"
	"testing"
)

type mockRegistrar struct {
	count int32
}

func (m *mockRegistrar) Register(ctx context.Context, service *registry.ServiceInstance) error {
	atomic.AddInt32(&m.count, 1)
	return nil
}

func (m *mockRegistrar) Deregister(ctx context.Context, service *registry.ServiceInstance) error {
	atomic.AddInt32(&m.count, -1)
	return nil
}

func TestRegistrarGroup(t *testing.T) {
	r1 := &mockRegistrar{}
	r2 := &mockRegistrar{}
	ins := &registry.ServiceInstance{}
	g := registry.RegistrarGroup(r1, r2)
	err := g.Register(context.Background(), ins)
	if err != nil {
		t.Fatal(err)
	}
	if int(r1.count) != 1 || int(r2.count) != 1 {
		t.Error("failed")
	}
	err = g.Deregister(context.Background(), ins)
	if err != nil {
		t.Fatal(err)
	}
	if int(r1.count) != 0 || int(r2.count) != 0 {
		t.Error("failed")
	}
}
