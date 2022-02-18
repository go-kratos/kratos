package polaris

import (
	"context"
	"testing"
	"time"

	"github.com/go-kratos/kratos/v2/registry"
	"github.com/polarismesh/polaris-go/api"
)

// TestRegistry . TestRegistryManyService
func TestRegistry(t *testing.T) {
	provider, err := api.NewProviderAPI()
	defer provider.Destroy()
	if err != nil {
		t.Fatal(err)
	}

	r := NewRegistry(provider, WithTimeout(1*time.Second))
	ctx := context.Background()

	schema := "tcp://127.0.0.1:9000?isSecure=false"
	svc := &registry.ServiceInstance{
		Name:      "kratos-provider-0-",
		Version:   "test",
		Metadata:  map[string]string{"app": "kratos"},
		Endpoints: []string{schema},
	}
	err = r.Register(ctx, svc)
	if err != nil {
		t.Fatal(err)
	}
	err = r.Deregister(ctx, svc)
	if err != nil {
		t.Fatal(err)
	}
}

// TestRegistryMany . TestRegistryManyService
func TestRegistryMany(t *testing.T) {
	provider, err := api.NewProviderAPI()
	defer provider.Destroy()
	if err != nil {
		t.Fatal(err)
	}

	r := NewRegistry(provider, WithTimeout(1*time.Second))

	// schema := "tcp://127.0.0.1:9000?isSecure=false"
	svc := &registry.ServiceInstance{
		Name:      "kratos-provider-0-",
		Version:   "test",
		Metadata:  map[string]string{"app": "kratos"},
		Endpoints: []string{"tcp://127.0.0.1:9000?isSecure=false"},
	}
	svc1 := &registry.ServiceInstance{
		Name:      "kratos-provider-1-",
		Version:   "test",
		Metadata:  map[string]string{"app": "kratos"},
		Endpoints: []string{"tcp://127.0.0.1:9001?isSecure=false"},
	}
	svc2 := &registry.ServiceInstance{
		Name:      "kratos-provider-2-",
		Version:   "test",
		Metadata:  map[string]string{"app": "kratos"},
		Endpoints: []string{"tcp://127.0.0.1:9002?isSecure=false"},
	}

	err = r.Register(context.Background(), svc)
	if err != nil {
		t.Fatal(err)
	}

	err = r.Register(context.Background(), svc1)
	if err != nil {
		t.Fatal(err)
	}

	err = r.Register(context.Background(), svc2)
	if err != nil {
		t.Fatal(err)
	}

	err = r.Deregister(context.Background(), svc)
	if err != nil {
		t.Fatal(err)
	}

	err = r.Deregister(context.Background(), svc1)
	if err != nil {
		t.Fatal(err)
	}

	err = r.Deregister(context.Background(), svc2)
	if err != nil {
		t.Fatal(err)
	}
}
