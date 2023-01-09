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
		WithRegistryTimeout(time.Second),
		WithRegistryHealthy(true),
		WithRegistryIsolate(false),
		WithRegistryRetryCount(3),
		WithRegistryWeight(100),
		WithRegistryTTL(1000),
	)

	ins := &registry.ServiceInstance{
		ID:      "test-ut-1",
		Name:    "test-ut-1",
		Version: "v1.0.0",
		Endpoints: []string{
			"grpc://127.0.0.1:8080",
			"http://127.0.0.1:9090",
		},
	}

	err = r.Register(context.Background(), ins)

	t.Cleanup(func() {
		if err = r.Deregister(context.Background(), ins); err != nil {
			t.Fatal(err)
		}
	})

	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(time.Second * 3)
	service, err := r.GetService(context.Background(), "test-ut-1")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(service)
}

func TestWatch(t *testing.T) {
	sdk, err := polaris.NewSDKContextByAddress("127.0.0.1:8091")
	if err != nil {
		t.Fatal(err)
	}

	p := New(sdk)

	r := p.Registry(
		WithRegistryTimeout(time.Second),
		WithRegistryHealthy(true),
		WithRegistryIsolate(false),
		WithRegistryRetryCount(3),
		WithRegistryWeight(100),
		WithRegistryTTL(1000),
	)

	err = r.Register(context.Background(), &registry.ServiceInstance{
		ID:      "test-ut",
		Name:    "test-ut",
		Version: "v1.0.0",
		Endpoints: []string{
			"grpc://127.0.0.1:8080",
			"http://127.0.0.1:9090",
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(time.Second * 3)
	w, err := r.Watch(context.Background(), "test-ut")
	if err != nil {
		t.Fatal(err)
	}
	service, err := w.Next()
	if err != nil {
		t.Fatal(err)
	}

	if len(service) != 1 {
		t.Fatalf("want 1, got %d, service %+v", len(service), service)
	}

	err = r.Register(context.Background(), &registry.ServiceInstance{
		ID:      "test-ut",
		Name:    "test-ut",
		Version: "v1.0.0",
		Endpoints: []string{
			"grpc://127.0.0.2:8080",
			"http://127.0.0.2:9090",
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(time.Second * 3)
	service, err = w.Next()
	if err != nil {
		t.Fatal(err)
	}
	if len(service) != 2 {
		t.Fatalf("want 1, got %d, service %+v", len(service), service)
	}

	err = r.Deregister(context.Background(), &registry.ServiceInstance{
		ID:      "test-ut",
		Name:    "test-ut",
		Version: "v1.0.0",
		Endpoints: []string{
			"grpc://127.0.0.1:8080",
			"http://127.0.0.1:9090",
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(time.Second * 3)
	service, err = w.Next()
	if err != nil {
		t.Fatal(err)
	}
	if len(service) != 1 {
		t.Fatalf("want 1, got %d, service %+v", len(service), service)
	}
	err = r.Deregister(context.Background(), &registry.ServiceInstance{
		ID:      "test-ut",
		Name:    "test-ut",
		Version: "v1.0.0",
		Endpoints: []string{
			"grpc://127.0.0.2:8080",
			"http://127.0.0.2:9090",
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(time.Second * 3)
	service, err = w.Next()
	if err != nil {
		t.Fatal(err)
	}
	if len(service) != 0 {
		t.Fatalf("want 0, got %d", len(service))
	}
}
