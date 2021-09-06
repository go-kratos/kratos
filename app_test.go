package kratos

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/go-kratos/kratos/v2/registry"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/stretchr/testify/assert"
)

func TestApp(t *testing.T) {
	hs := http.NewServer()
	gs := grpc.NewServer()
	app := New(
		Name("kratos"),
		Version("v1.0.0"),
		Server(hs, gs),
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
	assert.Equal(t, v, o.ID())
}

func TestApp_Name(t *testing.T) {
	v := "123"
	o := New(Name(v))
	assert.Equal(t, v, o.Name())
}

func TestApp_Version(t *testing.T) {
	v := "123"
	o := New(Version(v))
	assert.Equal(t, v, o.Version())
}

func TestApp_Metadata(t *testing.T) {
	v := map[string]string{
		"a": "1",
		"b": "2",
	}
	o := New(Metadata(v))
	assert.Equal(t, v, o.Metadata())
}

func TestApp_Endpoint(t *testing.T) {
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
				id: "3", version: "v3", name: "kratos-v3", endpoint: []string{},
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
