package zookeeper

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/go-kratos/kratos/v2/registry"
	"github.com/go-zookeeper/zk"
)

func TestRegistry_GetService(t *testing.T) {
	conn, _, err := zk.Connect([]string{"127.0.0.1:2181"}, time.Second*15)
	if err != nil {
		t.Fatal(err)
		return
	}
	r := New(conn)

	svrHello := &registry.ServiceInstance{
		ID:        "1",
		Name:      "hello",
		Version:   "v1.0.0",
		Endpoints: []string{"127.0.0.1:8080"},
	}

	type fields struct {
		registry *Registry
	}
	type args struct {
		ctx         context.Context
		serviceName string
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		want      []*registry.ServiceInstance
		wantErr   bool
		preFunc   func(t *testing.T)
		deferFunc func(t *testing.T)
	}{
		{
			name: "normal",
			preFunc: func(t *testing.T) {
				err = r.Register(context.Background(), svrHello)
				if err != nil {
					t.Error(err)
				}
			},
			deferFunc: func(t *testing.T) {
				err = r.Deregister(context.Background(), svrHello)
				if err != nil {
					t.Error(err)
				}
			},
			fields: fields{
				registry: r,
			},
			args: args{
				ctx:         context.Background(),
				serviceName: svrHello.Name,
			},
			want:    []*registry.ServiceInstance{svrHello},
			wantErr: false,
		},
		{
			name: "can't get any",
			fields: fields{
				registry: r,
			},
			args: args{
				ctx:         context.Background(),
				serviceName: "helloxxx",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "conn close",
			preFunc: func(t *testing.T) {
				conn.Close()
			},
			fields: fields{
				registry: r,
			},
			args: args{
				ctx:         context.Background(),
				serviceName: "hello",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.preFunc != nil {
				tt.preFunc(t)
			}
			if tt.deferFunc != nil {
				defer tt.deferFunc(t)
			}
			r := tt.fields.registry
			got, err := r.GetService(tt.args.ctx, tt.args.serviceName)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetService() error = %v, wantErr %v", err, tt.wantErr)
				t.Errorf("GetService() got = %v", got)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetService() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRegistry_Register(t *testing.T) {
	conn, _, err := zk.Connect([]string{"127.0.0.1:2181"}, time.Second*15)
	if err != nil {
		t.Fatal(err)
		return
	}
	r := New(conn)

	svrHello := &registry.ServiceInstance{
		ID:        "1",
		Name:      "hello",
		Version:   "v1.0.0",
		Endpoints: []string{"127.0.0.1:8080"},
	}

	type fields struct {
		registry *Registry
	}
	type args struct {
		ctx     context.Context
		service *registry.ServiceInstance
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		wantErr   bool
		preFunc   func(t *testing.T)
		deferFunc func(t *testing.T)
	}{
		{
			name: "normal",
			fields: fields{
				registry: r,
			},
			args: args{
				ctx:     context.Background(),
				service: svrHello,
			},
			wantErr: false,
		},
		{
			name: "invalid path",
			fields: fields{
				registry: New(conn, WithRootPath("invalid")),
			},
			args: args{
				ctx: context.Background(),
				service: &registry.ServiceInstance{
					ID:        "1",
					Name:      "hello1",
					Version:   "v1.0.0",
					Endpoints: []string{"127.0.0.1:8080"},
				},
			},
			wantErr: true,
		},
		{
			name: "auth",
			preFunc: func(t *testing.T) {
				err = conn.AddAuth("digest", []byte("test:test"))
				if err != nil {
					t.Error(err)
				}
			},
			fields: fields{
				registry: New(conn, WithRootPath("/tt1"), WithDigestACL("test", "test")),
			},
			args: args{
				ctx: context.Background(),
				service: &registry.ServiceInstance{
					ID:        "1",
					Name:      "hello2",
					Version:   "v1.0.0",
					Endpoints: []string{"127.0.0.1:8080"},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.preFunc != nil {
				tt.preFunc(t)
			}
			if tt.deferFunc != nil {
				defer tt.deferFunc(t)
			}
			r := tt.fields.registry
			if err := r.Register(tt.args.ctx, tt.args.service); (err != nil) != tt.wantErr {
				t.Errorf("Register() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRegistry_Deregister(t *testing.T) {
	conn, _, err := zk.Connect([]string{"127.0.0.1:2181"}, time.Second*15)
	if err != nil {
		t.Fatal(err)
		return
	}
	r := New(conn)

	svrHello := &registry.ServiceInstance{
		ID:        "1",
		Name:      "hello",
		Version:   "v1.0.0",
		Endpoints: []string{"127.0.0.1:8080"},
	}

	cancelCtx, cancel := context.WithCancel(context.Background())
	cancel()
	type fields struct {
		registry *Registry
	}
	type args struct {
		ctx     context.Context
		service *registry.ServiceInstance
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		wantErr   bool
		preFunc   func(t *testing.T)
		deferFunc func(t *testing.T)
	}{
		{
			name: "normal",
			fields: fields{
				registry: r,
			},
			args: args{
				ctx:     context.Background(),
				service: svrHello,
			},
			wantErr: false,
			preFunc: func(t *testing.T) {
				err = r.Register(context.Background(), svrHello)
				if err != nil {
					t.Error(err)
				}
			},
		},
		{
			name: "with ctx cancel",
			fields: fields{
				registry: r,
			},
			args: args{
				ctx:     cancelCtx,
				service: svrHello,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := tt.fields.registry
			if err := r.Deregister(tt.args.ctx, tt.args.service); (err != nil) != tt.wantErr {
				t.Errorf("Deregister() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRegistry_Watch(t *testing.T) {
	conn, _, err := zk.Connect([]string{"127.0.0.1:2181"}, time.Second*15)
	if err != nil {
		t.Fatal(err)
		return
	}
	closeConn, _, err := zk.Connect([]string{"127.0.0.1:2181"}, time.Second*15)
	if err != nil {
		t.Fatal(err)
		return
	}
	r := New(conn)

	svrHello := &registry.ServiceInstance{
		ID:        "1",
		Name:      "hello",
		Version:   "v1.0.0",
		Endpoints: []string{"127.0.0.1:8080"},
	}

	cancelCtx, cancel := context.WithCancel(context.Background())
	type fields struct {
		registry *Registry
	}
	type args struct {
		ctx         context.Context
		serviceName string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
		want    []*registry.ServiceInstance

		preFunc     func(t *testing.T)
		deferFunc   func(t *testing.T)
		processFunc func(t *testing.T, w registry.Watcher)
	}{
		{
			name: "normal",
			fields: fields{
				registry: r,
			},
			args: args{
				ctx:         context.Background(),
				serviceName: svrHello.Name,
			},
			wantErr: false,
			want:    []*registry.ServiceInstance{svrHello},
			deferFunc: func(t *testing.T) {
				err = r.Deregister(context.Background(), svrHello)
				if err != nil {
					t.Error(err)
				}
			},
			processFunc: func(t *testing.T, w registry.Watcher) {
				err = r.Register(context.Background(), svrHello)
				if err != nil {
					t.Error(err)
				}
			},
		},
		{
			name: "ctx cancel",
			fields: fields{
				registry: r,
			},
			args: args{
				ctx:         cancelCtx,
				serviceName: svrHello.Name,
			},
			wantErr: true,
			want:    nil,
			processFunc: func(t *testing.T, w registry.Watcher) {
				cancel()
			},
		},
		{
			name: "disconnect",
			fields: fields{
				registry: New(closeConn),
			},
			args: args{
				ctx:         context.Background(),
				serviceName: svrHello.Name,
			},
			wantErr: true,
			want:    nil,
			processFunc: func(t *testing.T, w registry.Watcher) {
				closeConn.Close()
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.preFunc != nil {
				tt.preFunc(t)
			}
			if tt.deferFunc != nil {
				defer tt.deferFunc(t)
			}
			r := tt.fields.registry
			watcher, err := r.Watch(tt.args.ctx, tt.args.serviceName)
			if err != nil {
				t.Error(err)
				return
			}
			defer func() {
				err = watcher.Stop()
				if err != nil {
					t.Error(err)
				}
			}()
			_, err = watcher.Next()
			if err != nil {
				t.Error(err)
				return
			}

			if tt.processFunc != nil {
				tt.processFunc(t, watcher)
			}

			want, err := watcher.Next()
			if (err != nil) != tt.wantErr {
				t.Errorf("Watch() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(want, tt.want) {
				t.Errorf("Watch() watcher = %v, want %v", watcher, tt.want)
			}
		})
	}
}
