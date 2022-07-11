package kratos

import (
	"context"
	"fmt"
	"net/url"
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/go-kratos/kratos/v2/registry"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
)

type mockRegistry struct {
	lk      sync.Mutex
	service map[string]*registry.ServiceInstance
}

func (r *mockRegistry) Register(ctx context.Context, service *registry.ServiceInstance) error {
	if service == nil || service.ID == "" {
		return fmt.Errorf("no service id")
	}
	r.lk.Lock()
	defer r.lk.Unlock()
	r.service[service.ID] = service
	return nil
}

// Deregister the registration.
func (r *mockRegistry) Deregister(ctx context.Context, service *registry.ServiceInstance) error {
	r.lk.Lock()
	defer r.lk.Unlock()
	if r.service[service.ID] == nil {
		return fmt.Errorf("deregister service not found")
	}
	delete(r.service, service.ID)
	return nil
}

func TestApp(t *testing.T) {
	hs := http.NewServer()
	gs := grpc.NewServer()
	app := New(
		Name("kratos"),
		Version("v1.0.0"),
		Server(hs, gs),
		Registrar(&mockRegistry{service: make(map[string]*registry.ServiceInstance)}),
	)
	time.AfterFunc(time.Second, func() {
		_ = app.Stop()
	})
	if err := app.Run(); err != nil {
		t.Fatal(err)
	}
}

func TestApp_ID(t *testing.T) {
	v := "123"
	o := New(ID(v))
	if !reflect.DeepEqual(v, o.ID()) {
		t.Fatalf("o.ID():%s is not equal to v:%s", o.ID(), v)
	}
}

func TestApp_Name(t *testing.T) {
	v := "123"
	o := New(Name(v))
	if !reflect.DeepEqual(v, o.Name()) {
		t.Fatalf("o.Name():%s is not equal to v:%s", o.Name(), v)
	}
}

func TestApp_Version(t *testing.T) {
	v := "123"
	o := New(Version(v))
	if !reflect.DeepEqual(v, o.Version()) {
		t.Fatalf("o.Version():%s is not equal to v:%s", o.Version(), v)
	}
}

func TestApp_Metadata(t *testing.T) {
	v := map[string]string{
		"a": "1",
		"b": "2",
	}
	o := New(Metadata(v))
	if !reflect.DeepEqual(v, o.Metadata()) {
		t.Fatalf("o.Metadata():%s is not equal to v:%s", o.Metadata(), v)
	}
}

func TestApp_Endpoint(t *testing.T) {
	v := []string{"https://go-kratos.dev", "localhost"}
	var endpoints []*url.URL
	for _, urlStr := range v {
		if endpoint, err := url.Parse(urlStr); err != nil {
			t.Errorf("invalid endpoint:%v", urlStr)
		} else {
			endpoints = append(endpoints, endpoint)
		}
	}
	o := New(Endpoint(endpoints...))
	if instance, err := o.buildInstance(); err != nil {
		t.Error("build instance failed")
	} else {
		o.instance = instance
	}
	if !reflect.DeepEqual(o.Endpoint(), v) {
		t.Errorf("Endpoint() = %v, want %v", o.Endpoint(), v)
	}
}

func TestApp_buildInstance(t *testing.T) {
	target := struct {
		id        string
		name      string
		version   string
		metadata  map[string]string
		endpoints []string
	}{
		id:      "1",
		name:    "kratos",
		version: "v1.0.0",
		metadata: map[string]string{
			"a": "1",
			"b": "2",
		},
		endpoints: []string{"https://go-kratos.dev", "localhost"},
	}
	var endpoints []*url.URL
	for _, urlStr := range target.endpoints {
		if endpoint, err := url.Parse(urlStr); err != nil {
			t.Errorf("invalid endpoint:%v", urlStr)
		} else {
			endpoints = append(endpoints, endpoint)
		}
	}
	app := New(
		ID(target.id),
		Name(target.name),
		Version(target.version),
		Metadata(target.metadata),
		Endpoint(endpoints...),
	)
	if instance, err := app.buildInstance(); err != nil {
		t.Error("build instance failed")
	} else {
		if instance.ID != target.id {
			t.Errorf("ID() = %v, want %v", instance.ID, target.id)
		}
		if instance.Name != target.name {
			t.Errorf("Name() = %v, want %v", instance.Name, target.name)
		}
		if instance.Version != target.version {
			t.Errorf("Version() = %v, want %v", instance.Version, target.version)
		}
		if !reflect.DeepEqual(instance.Endpoints, target.endpoints) {
			t.Errorf("Endpoint() = %v, want %v", instance.Endpoints, target.endpoints)
		}
		if !reflect.DeepEqual(instance.Metadata, target.metadata) {
			t.Errorf("Metadata() = %v, want %v", instance.Metadata, target.metadata)
		}
	}
}

func TestApp_NewContext(t *testing.T) {
	app := New(
		ID("testId"),
	)
	NewContext(context.Background(), app)
}

func TestApp_FromContext(t *testing.T) {
	target := struct {
		id        string
		version   string
		name      string
		metadata  map[string]string
		endpoints []string
	}{
		id:        "1",
		name:      "kratos-v1",
		metadata:  map[string]string{},
		version:   "v1",
		endpoints: []string{"https://go-kratos.dev", "localhost"},
	}
	a := &App{
		opts:     options{id: target.id, name: target.name, metadata: target.metadata, version: target.version},
		ctx:      context.Background(),
		cancel:   nil,
		instance: &registry.ServiceInstance{Endpoints: target.endpoints},
	}

	ctx := NewContext(context.Background(), a)

	if got, ok := FromContext(ctx); ok {
		if got.ID() != target.id {
			t.Errorf("ID() = %v, want %v", got.ID(), target.id)
		}
		if got.Name() != target.name {
			t.Errorf("Name() = %v, want %v", got.Name(), target.name)
		}
		if got.Version() != target.version {
			t.Errorf("Version() = %v, want %v", got.Version(), target.version)
		}
		if !reflect.DeepEqual(got.Endpoint(), target.endpoints) {
			t.Errorf("Endpoint() = %v, want %v", got.Endpoint(), target.endpoints)
		}
		if !reflect.DeepEqual(got.Metadata(), target.metadata) {
			t.Errorf("Metadata() = %v, want %v", got.Metadata(), target.metadata)
		}
	} else {
		t.Errorf("ok() = %v, want %v", ok, true)
	}
}
