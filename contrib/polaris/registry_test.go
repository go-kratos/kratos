package polaris

import (
	"context"
	"testing"
	"time"

	"github.com/polarismesh/polaris-go"

	"github.com/go-kratos/kratos/v2/registry"
)

// TestRegistry
func TestRegistry(t *testing.T) {
	sdk, err := polaris.NewSDKContextByAddress("127.0.0.1:8091")
	if err != nil {
		t.Fatal(err)
	}

	p := New(sdk)

	r := p.Registry(
		WithTimeout(time.Second),
		WithHealthy(true),
		WithIsolate(false),
		WithRegistryNamespace("default"),
		WithRetryCount(0),
		WithWeight(100),
		WithTTL(10),
	)

	ctx := context.Background()

	svc := &registry.ServiceInstance{
		Name:      "kratos-provider",
		Version:   "test",
		Metadata:  map[string]string{"app": "kratos"},
		Endpoints: []string{"grpc://127.0.0.1:9000", "http://127.0.0.1:8000"},
	}

	err = r.Register(ctx, svc)
	if err != nil {
		t.Fatal(err)
	}

	time.Sleep(time.Second)

	result, err := r.GetService(context.Background(), "kratos-provider")
	if err != nil {
		t.Fatal(err)
	}

	for _, instance := range result {
		t.Log(instance)
	}
}

// TestRegistryMany
func TestRegistryMany(t *testing.T) {
	sdk, err := polaris.NewSDKContextByAddress("127.0.0.1:8091")
	if err != nil {
		t.Fatal(err)
	}

	p := New(sdk)
	r := p.Registry(
		WithTimeout(time.Second),
		//WithHealthy(true),
		WithIsolate(true),
		WithRegistryNamespace("default"),
		//WithProtocol("tcp"),
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
