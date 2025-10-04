package kratos

import (
	"context"
	"errors"
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

func (r *mockRegistry) Register(_ context.Context, service *registry.ServiceInstance) error {
	if service == nil || service.ID == "" {
		return errors.New("no service id")
	}
	r.lk.Lock()
	defer r.lk.Unlock()
	r.service[service.ID] = service
	return nil
}

// Deregister the registration.
func (r *mockRegistry) Deregister(_ context.Context, service *registry.ServiceInstance) error {
	r.lk.Lock()
	defer r.lk.Unlock()
	if r.service[service.ID] == nil {
		return errors.New("deregister service not found")
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
		BeforeStart(func(_ context.Context) error {
			t.Log("BeforeStart...")
			return nil
		}),
		BeforeStop(func(_ context.Context) error {
			t.Log("BeforeStop...")
			return nil
		}),
		AfterStart(func(_ context.Context) error {
			t.Log("AfterStart...")
			return nil
		}),
		AfterStop(func(_ context.Context) error {
			t.Log("AfterStop...")
			return nil
		}),
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
	want := struct {
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
	for _, urlStr := range want.endpoints {
		if endpoint, err := url.Parse(urlStr); err != nil {
			t.Errorf("invalid endpoint:%v", urlStr)
		} else {
			endpoints = append(endpoints, endpoint)
		}
	}
	app := New(
		ID(want.id),
		Name(want.name),
		Version(want.version),
		Metadata(want.metadata),
		Endpoint(endpoints...),
	)
	if got, err := app.buildInstance(); err != nil {
		t.Error("build got failed")
	} else {
		if got.ID != want.id {
			t.Errorf("ID() = %v, want %v", got.ID, want.id)
		}
		if got.Name != want.name {
			t.Errorf("Name() = %v, want %v", got.Name, want.name)
		}
		if got.Version != want.version {
			t.Errorf("Version() = %v, want %v", got.Version, want.version)
		}
		if !reflect.DeepEqual(got.Endpoints, want.endpoints) {
			t.Errorf("Endpoint() = %v, want %v", got.Endpoints, want.endpoints)
		}
		if !reflect.DeepEqual(got.Metadata, want.metadata) {
			t.Errorf("Metadata() = %v, want %v", got.Metadata, want.metadata)
		}
	}
}

func TestApp_Context(t *testing.T) {
	type fields struct {
		id       string
		version  string
		name     string
		instance *registry.ServiceInstance
		metadata map[string]string
		want     struct {
			id       string
			version  string
			name     string
			endpoint []string
			metadata map[string]string
		}
	}
	tests := []fields{
		{
			id:       "1",
			name:     "kratos-v1",
			instance: &registry.ServiceInstance{Endpoints: []string{"https://go-kratos.dev", "localhost"}},
			metadata: map[string]string{},
			version:  "v1",
			want: struct {
				id       string
				version  string
				name     string
				endpoint []string
				metadata map[string]string
			}{
				id: "1", version: "v1", name: "kratos-v1", endpoint: []string{"https://go-kratos.dev", "localhost"},
				metadata: map[string]string{},
			},
		},
		{
			id:       "2",
			name:     "kratos-v2",
			instance: &registry.ServiceInstance{Endpoints: []string{"test"}},
			metadata: map[string]string{"kratos": "https://github.com/go-kratos/kratos"},
			version:  "v2",
			want: struct {
				id       string
				version  string
				name     string
				endpoint []string
				metadata map[string]string
			}{
				id: "2", version: "v2", name: "kratos-v2", endpoint: []string{"test"},
				metadata: map[string]string{"kratos": "https://github.com/go-kratos/kratos"},
			},
		},
		{
			id:       "3",
			name:     "kratos-v3",
			instance: nil,
			metadata: make(map[string]string),
			version:  "v3",
			want: struct {
				id       string
				version  string
				name     string
				endpoint []string
				metadata map[string]string
			}{
				id: "3", version: "v3", name: "kratos-v3", endpoint: nil,
				metadata: map[string]string{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &App{
				opts:     options{id: tt.id, name: tt.name, metadata: tt.metadata, version: tt.version},
				ctx:      context.Background(),
				cancel:   nil,
				instance: tt.instance,
			}

			ctx := NewContext(context.Background(), a)

			if got, ok := FromContext(ctx); ok {
				if got.ID() != tt.want.id {
					t.Errorf("ID() = %v, want %v", got.ID(), tt.want.id)
				}
				if got.Name() != tt.want.name {
					t.Errorf("Name() = %v, want %v", got.Name(), tt.want.name)
				}
				if got.Version() != tt.want.version {
					t.Errorf("Version() = %v, want %v", got.Version(), tt.want.version)
				}
				if !reflect.DeepEqual(got.Endpoint(), tt.want.endpoint) {
					t.Errorf("Endpoint() = %v, want %v", got.Endpoint(), tt.want.endpoint)
				}
				if !reflect.DeepEqual(got.Metadata(), tt.want.metadata) {
					t.Errorf("Metadata() = %v, want %v", got.Metadata(), tt.want.metadata)
				}
			} else {
				t.Errorf("ok() = %v, want %v", ok, true)
			}
		})
	}
}

func TestApp_ContextCanceled(t *testing.T) {
	ctx, stop := context.WithCancel(context.Background())
	stopFn := func(ctx context.Context) error {
		select {
		case <-ctx.Done():
			t.Fatal("context should not be done yet")
		default:
		}
		return nil
	}
	app := New(Context(ctx), Server(&mockServer{stopFn: stopFn}), StopTimeout(time.Hour))
	time.AfterFunc(time.Millisecond*10, stop)
	_ = app.Run()
}
