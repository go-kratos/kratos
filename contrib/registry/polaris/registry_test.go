package polaris

import (
	"context"
	"github.com/go-kratos/kratos/v2/registry"
	"github.com/polarismesh/polaris-go/api"
	"testing"
	"time"
)

// TestRegistry . TestRegistryManyService
func TestRegistry(t *testing.T) {
	provider, err := api.NewProviderAPI()
	defer provider.Destroy()
	if err != nil {
		t.Fatal(err)
	}

	r := New(provider, WithDefaultTimeout(1*time.Second))
	ctx := context.Background()

	schema := "tcp://127.0.0.1?isSecure=false"
	svc := &registry.ServiceInstance{
		Name:      "kratos-provider-",
		Version:   "test",
		Metadata:  map[string]string{"app": "kratos"},
		Endpoints: []string{schema},
	}
	err = r.Register(ctx, svc)
	if err != nil {
		t.Fatal(err)
	}
	// 暂停20秒后Deregister
	time.Sleep(20 * time.Second)
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

	r := New(provider, WithDefaultTimeout(1*time.Second))
	ctx := context.Background()

	schema := "tcp://127.0.0.1:9000?isSecure=false"
	svc := &registry.ServiceInstance{
		Name:      "kratos-provider-",
		Version:   "test",
		Metadata:  map[string]string{"app": "kratos"},
		Endpoints: []string{schema},
	}
	svc1 := &registry.ServiceInstance{
		Name:      "kratos-provider-1-",
		Version:   "test",
		Metadata:  map[string]string{"app": "kratos"},
		Endpoints: []string{schema},
	}
	svc2 := &registry.ServiceInstance{
		Name:      "kratos-provider-2-",
		Version:   "test",
		Metadata:  map[string]string{"app": "kratos"},
		Endpoints: []string{schema},
	}

	err = r.Register(ctx, svc)
	if err != nil {
		t.Fatal(err)
	}
	err = r.Register(ctx, svc1)
	if err != nil {
		t.Fatal(err)
	}
	err = r.Register(ctx, svc2)
	if err != nil {
		t.Fatal(err)
	}

	// 暂停20秒后Deregister
	time.Sleep(20 * time.Second)
	err = r.Deregister(ctx, svc)
	if err != nil {
		t.Fatal(err)
	}

	err = r.Deregister(ctx, svc1)
	if err != nil {
		t.Fatal(err)
	}

	err = r.Deregister(ctx, svc2)
	if err != nil {
		t.Fatal(err)
	}

}
