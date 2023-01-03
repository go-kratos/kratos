package polaris

import (
	"context"
	"encoding/json"
	"github.com/go-kratos/kratos/v2/registry"
	"os"
	"testing"
	"time"

	"github.com/polarismesh/polaris-go"
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
	service, err := r.GetService(context.Background(), "test-ut")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(service)
}

func TestDeregister(t *testing.T) {
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
}

func TestWatch(t *testing.T) {
	sdk, err := polaris.NewSDKContextByAddress("127.0.0.1:8091")
	if err != nil {
		t.Fatal(err)
	}

	p := New(sdk)

	r := p.Registry(
		WithTimeout(time.Second),
		WithHealthy(false),
		WithIsolate(false),
		WithRegistryNamespace("default"),
		WithRetryCount(0),
		WithWeight(100),
		WithTTL(10),
	)

	w, err := r.Watch(context.Background(), "test-ut")
	if err != nil {
		t.Fatal(err)
	}

	ch := make(chan struct{})
	go func(t *testing.T) {
		for {
			next, err := w.Next()
			if err != nil {
				t.Error(err)
				os.Exit(1)
			}
			bytes, err := json.Marshal(next)
			if err != nil {
				t.Error(err)
				os.Exit(1)
			}
			t.Log(string(bytes))
			if len(next) == 0 {
				ch <- struct{}{}
			}
		}
	}(t)

	err = r.Register(context.Background(), &registry.ServiceInstance{
		ID:      "test-ut",
		Name:    "test-ut",
		Version: "v1.0.0",
		Endpoints: []string{
			"grpc://127.0.0.1:8080",
			"http://127.0.0.1:9090",
		},
	})
	time.Sleep(time.Second * 2)
	err = r.Register(context.Background(), &registry.ServiceInstance{
		ID:      "test-ut",
		Name:    "test-ut",
		Version: "v1.0.0",
		Endpoints: []string{
			"grpc://127.0.0.2:8080",
			"http://127.0.0.2:9090",
		},
	})
	time.Sleep(time.Second * 2)
	err = r.Deregister(context.Background(), &registry.ServiceInstance{
		ID:      "test-ut",
		Name:    "test-ut",
		Version: "v1.0.0",
		Endpoints: []string{
			"grpc://127.0.0.1:8080",
			"http://127.0.0.1:9090",
		},
	})
	time.Sleep(time.Second * 2)
	err = r.Deregister(context.Background(), &registry.ServiceInstance{
		ID:      "test-ut",
		Name:    "test-ut",
		Version: "v1.0.0",
		Endpoints: []string{
			"grpc://127.0.0.2:8080",
			"http://127.0.0.2:9090",
		},
	})
	<-ch
}
