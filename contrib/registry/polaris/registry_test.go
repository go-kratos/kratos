package polaris

import (
	"context"
	"testing"
	"time"

	"github.com/go-kratos/kratos/v2/registry"

	"github.com/polarismesh/polaris-go/pkg/config"
)

// TestRegistry
func TestRegistry(t *testing.T) {
	conf := config.NewDefaultConfiguration([]string{"127.0.0.1:8091"})

	r := NewRegistryWithConfig(
		conf,
		WithTimeout(time.Second),
		WithHeartbeat(true),
		WithHealthy(true),
		WithIsolate(true),
		WithNamespace("default"),
		WithProtocol("tcp"),
		WithRetryCount(0),
		WithWeight(100),
		WithTTL(10),
	)

	ctx := context.Background()

	svc := &registry.ServiceInstance{
		Name:      "kratos-provider-0-",
		Version:   "test",
		Metadata:  map[string]string{"app": "kratos"},
		Endpoints: []string{"tcp://127.0.0.1:9000?isSecure=false"},
	}

	err := r.Register(ctx, svc)
	if err != nil {
		t.Fatal(err)
	}

	time.Sleep(time.Second)

	result, err := r.GetService(context.Background(), "kratos-provider-0-tcp")
	if err != nil {
		t.Fatal(err)
	}

	if len(result) != 1 {
		t.Fatal("register error")
	}

	for _, item := range result {
		if item.Name != "kratos-provider-0-tcp" || item.Endpoints[0] != "tcp://127.0.0.1:9000" {
			t.Fatal("register error")
		}
	}

	watch, err := r.Watch(ctx, "kratos-provider-0-tcp")
	if err != nil {
		t.Fatal(err)
	}

	// Test update
	svc.Version = "release1.0.0"

	if err = r.Register(ctx, svc); err != nil {
		t.Fatal(err)
	}

	result, err = watch.Next()

	if err != nil {
		t.Fatal(err)
	}

	if len(result) != 1 || result[0].Version != "release1.0.0" {
		t.Fatal("register error")
	}
	// Test add instance
	svc1 := &registry.ServiceInstance{
		Name:      "kratos-provider-0-",
		Version:   "test",
		Metadata:  map[string]string{"app": "kratos"},
		Endpoints: []string{"tcp://127.0.0.1:9001?isSecure=false"},
	}

	if err = r.Register(ctx, svc1); err != nil {
		t.Fatal(err)
	}

	if _, err = watch.Next(); err != nil {
		t.Fatal(err)
	}

	result, err = r.GetService(ctx, "kratos-provider-0-tcp")

	if err != nil {
		t.Fatal(err)
	}

	if len(result) != 2 {
		t.Fatal("register error")
	}

	if err = r.Deregister(ctx, svc); err != nil {
		t.Fatal(err)
	}
	if err = r.Deregister(ctx, svc1); err != nil {
		t.Fatal(err)
	}

	result, err = watch.Next()
	if err != nil {
		t.Fatal(err)
	}

	if len(result) != 0 {
		t.Fatal("register error")
	}
}

// TestRegistryMany
func TestRegistryMany(t *testing.T) {
	conf := config.NewDefaultConfiguration([]string{"127.0.0.1:8091"})

	r := NewRegistryWithConfig(
		conf,
		WithTimeout(time.Second),
		WithHeartbeat(true),
		WithHealthy(true),
		WithIsolate(true),
		WithNamespace("default"),
		WithProtocol("tcp"),
		WithRetryCount(0),
		WithWeight(100),
		WithTTL(10),
	)

	ctx := context.Background()

	// Multi endpoint
	svc := &registry.ServiceInstance{
		Name:      "kratos-provider-1-",
		Version:   "test",
		Metadata:  map[string]string{"app": "kratos"},
		Endpoints: []string{"tcp://127.0.0.1:9000?isSecure=false", "tcp://127.0.0.1:9001?isSecure=false"},
	}
	// Normal
	svc1 := &registry.ServiceInstance{
		Name:      "kratos-provider-2-",
		Version:   "test",
		Metadata:  map[string]string{"app": "kratos"},
		Endpoints: []string{"tcp://127.0.0.1:9002?isSecure=false"},
	}
	// Without metadata
	svc2 := &registry.ServiceInstance{
		Name:      "kratos-provider-3-",
		Version:   "test",
		Endpoints: []string{"tcp://127.0.0.1:9003?isSecure=false"},
	}

	if err := r.Register(ctx, svc); err != nil {
		t.Fatal(err)
	}

	if err := r.Register(ctx, svc1); err != nil {
		t.Fatal(err)
	}

	if err := r.Register(ctx, svc2); err != nil {
		t.Fatal(err)
	}

	time.Sleep(3 * time.Second)

	result1, err := r.GetService(ctx, "kratos-provider-1-tcp")

	if err != nil || len(result1) != 2 || result1[0].Name != "kratos-provider-1-tcp" {
		t.Fatal(err)
	}

	result2, err := r.GetService(ctx, "kratos-provider-2-tcp")

	if err != nil || len(result2) != 1 || result2[0].Name != "kratos-provider-2-tcp" || result2[0].Endpoints[0] != "tcp://127.0.0.1:9002" {
		t.Fatal(err)
	}

	result3, err := r.GetService(ctx, "kratos-provider-3-tcp")

	if err != nil || len(result3) != 1 || result3[0].Name != "kratos-provider-3-tcp" || result3[0].Endpoints[0] != "tcp://127.0.0.1:9003" {
		t.Fatal(err)
	}

	watch1, err := r.Watch(ctx, "kratos-provider-1-tcp")
	if err != nil {
		t.Fatal(err)
	}
	watch2, err := r.Watch(ctx, "kratos-provider-2-tcp")
	if err != nil {
		t.Fatal(err)
	}
	watch3, err := r.Watch(ctx, "kratos-provider-3-tcp")
	if err != nil {
		t.Fatal(err)
	}

	if err = r.Deregister(ctx, svc); err != nil {
		t.Fatal(err)
	}

	result1, err = watch1.Next()
	if err != nil || len(result1) != 0 {
		t.Fatal("deregister error")
	}

	err = r.Deregister(ctx, svc1)
	if err != nil {
		t.Fatal(err)
	}

	result2, err = watch2.Next()
	if err != nil || len(result2) != 0 {
		t.Fatal("deregister error")
	}
	err = r.Deregister(ctx, svc2)
	if err != nil {
		t.Fatal(err)
	}

	result3, err = watch3.Next()
	if err != nil || len(result3) != 0 {
		t.Fatal("deregister error")
	}
}
