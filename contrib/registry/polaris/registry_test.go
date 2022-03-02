package polaris

import (
	"context"
	"testing"
	"time"

	"github.com/go-kratos/kratos/v2/log"

	"github.com/polarismesh/polaris-go/pkg/config"

	"github.com/go-kratos/kratos/v2/registry"
)

// TestRegistry . TestRegistryManyService
func TestRegistry(t *testing.T) {
	conf := config.NewDefaultConfiguration([]string{"127.0.0.1:8091"})

	r := NewRegistryWithConfig(
		conf,
		WithTimeout(time.Second*10),
		WithTTL(100),
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

	err = r.Deregister(ctx, svc)
	if err != nil {
		t.Fatal(err)
	}
}

// TestRegistryMany . TestRegistryManyService
func TestRegistryMany(t *testing.T) {
	conf := config.NewDefaultConfiguration([]string{"127.0.0.1:8091"})

	r := NewRegistryWithConfig(
		conf,
		WithTimeout(time.Second*10),
		WithTTL(100),
	)

	svc := &registry.ServiceInstance{
		Name:      "kratos-provider-1-",
		Version:   "test",
		Metadata:  map[string]string{"app": "kratos"},
		Endpoints: []string{"tcp://127.0.0.1:9000?isSecure=false"},
	}
	svc1 := &registry.ServiceInstance{
		Name:      "kratos-provider-2-",
		Version:   "test",
		Metadata:  map[string]string{"app": "kratos"},
		Endpoints: []string{"tcp://127.0.0.1:9001?isSecure=false"},
	}
	svc2 := &registry.ServiceInstance{
		Name:      "kratos-provider-3-",
		Version:   "test",
		Metadata:  map[string]string{"app": "kratos"},
		Endpoints: []string{"tcp://127.0.0.1:9002?isSecure=false"},
	}

	err := r.Register(context.Background(), svc)
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

// TestGetService . TestGetService
func TestGetService(t *testing.T) {
	conf := config.NewDefaultConfiguration([]string{"127.0.0.1:8091"})

	r := NewRegistryWithConfig(
		conf,
		WithTimeout(time.Second*10),
		WithTTL(100),
	)

	ctx := context.Background()

	svc := &registry.ServiceInstance{
		Name:      "kratos-provider-4-",
		Version:   "test",
		Metadata:  map[string]string{"app": "kratos"},
		Endpoints: []string{"tcp://127.0.0.1:9000?isSecure=false"},
	}

	err := r.Register(ctx, svc)
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(time.Second * 1)
	serviceInstances, err := r.GetService(ctx, "kratos-provider-4-tcp")
	if err != nil {
		t.Fatal(err)
	}
	for _, instance := range serviceInstances {
		log.Info(instance)
	}

	err = r.Deregister(ctx, svc)
	if err != nil {
		t.Fatal(err)
	}
}

// TestWatch . TestWatch
func TestWatch(t *testing.T) {
	conf := config.NewDefaultConfiguration([]string{"127.0.0.1:8091"})

	r := NewRegistryWithConfig(
		conf,
		WithTimeout(time.Second*10),
		WithTTL(100),
	)

	svc := &registry.ServiceInstance{
		Name:      "kratos-provider-4-",
		Version:   "test",
		Metadata:  map[string]string{"app": "kratos"},
		Endpoints: []string{"tcp://127.0.0.1:9000?isSecure=false"},
	}

	watch, err := r.Watch(context.Background(), "kratos-provider-4-tcp")
	if err != nil {
		t.Fatal(err)
	}

	err = r.Register(context.Background(), svc)
	if err != nil {
		t.Fatal(err)
	}
	// watch svc
	time.Sleep(time.Second * 1)

	// svc register, AddEvent
	next, err := watch.Next()
	if err != nil {
		t.Fatal(err)
	}
	for _, instance := range next {
		// it will output one instance
		log.Info(instance)
	}

	err = r.Deregister(context.Background(), svc)
	if err != nil {
		t.Fatal(err)
	}

	// svc deregister, DeleteEvent
	next, err = watch.Next()
	if err != nil {
		t.Fatal(err)
	}
	for _, instance := range next {
		// it will output nothing
		log.Info(instance)
	}

	err = watch.Stop()
	if err != nil {
		t.Fatal(err)
	}
	_, err = watch.Next()
	if err == nil {
		// if nil, stop failed
		t.Fatal()
	}
}
