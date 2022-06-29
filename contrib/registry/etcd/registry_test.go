package etcd

import (
	"context"
	"fmt"
	"reflect"
	"testing"
	"time"

	"google.golang.org/grpc"

	"github.com/go-kratos/kratos/v2/registry"
	clientv3 "go.etcd.io/etcd/client/v3"
)

func TestRegistry_GetService(t *testing.T) {
	client, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"127.0.0.1:2379"},
		DialTimeout: time.Second, DialOptions: []grpc.DialOption{grpc.WithBlock()},
	})
	if err != nil {
		t.Fatal(err)
	}
	defer client.Close()
	s := &registry.ServiceInstance{
		ID:        "1",
		Name:      "hello",
		Version:   "v1.0.0",
		Endpoints: []string{"127.0.0.1:8080"},
	}

	r := New(client)
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
				err = r.Register(context.Background(), s)
				if err != nil {
					t.Error(err)
				}
			},
			deferFunc: func(t *testing.T) {
				err = r.Deregister(context.Background(), s)
				if err != nil {
					t.Error(err)
				}
			},
			fields: fields{
				registry: r,
			},
			args: args{
				ctx:         context.Background(),
				serviceName: s.Name,
			},
			want:    []*registry.ServiceInstance{s},
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
			want:    []*registry.ServiceInstance{},
			wantErr: false,
		},
		{
			name: "conn close",
			preFunc: func(t *testing.T) {
				client.Close()
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
	client, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"127.0.0.1:2379"},
		DialTimeout: time.Second, DialOptions: []grpc.DialOption{grpc.WithBlock()},
	})
	if err != nil {
		t.Fatal(err)
	}
	defer client.Close()
	s := &registry.ServiceInstance{
		ID:        "1",
		Name:      "hello",
		Version:   "v1.0.0",
		Endpoints: []string{"127.0.0.1:8080"},
	}

	r := New(client)
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
		want      []*registry.ServiceInstance
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
				service: s,
			},
			wantErr: false,
		},
		{
			name: "namespace",
			fields: fields{
				registry: New(client, Namespace("invalid")),
			},
			args: args{
				ctx: context.Background(),
				service: &registry.ServiceInstance{
					ID:        "2",
					Name:      "hello1",
					Version:   "v1.0.0",
					Endpoints: []string{"127.0.0.1:8080"},
				},
			},
			wantErr: false,
		},
		{
			name: "ttl",
			fields: fields{
				registry: New(client, RegisterTTL(3*time.Second)),
			},
			args: args{
				ctx: context.Background(),
				service: &registry.ServiceInstance{
					ID:        "3",
					Name:      "hello_ttl",
					Version:   "v1.0.0",
					Endpoints: []string{"127.0.0.1:8080"},
				},
			},
			wantErr: false,
			want:    []*registry.ServiceInstance{},
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
			if tt.name == "ttl" {
				time.Sleep(5 * time.Second)
				got, err := r.GetService(tt.args.ctx, tt.args.service.Name)
				if (err != nil) != tt.wantErr {
					t.Errorf("GetService() error = %v, wantErr %v", err, tt.wantErr)
					t.Errorf("GetService() got = %v", got)
					return
				}
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("GetService() got = %v, want %v", got, tt.want)
				}
			}
		})
	}
}
func TestHeartBeat(t *testing.T) {
	client, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"127.0.0.1:2379"},
		DialTimeout: time.Second, DialOptions: []grpc.DialOption{grpc.WithBlock()},
	})
	if err != nil {
		t.Fatal(err)
	}
	defer client.Close()

	ctx := context.Background()
	s := &registry.ServiceInstance{
		ID:   "0",
		Name: "helloworld",
	}

	go func() {
		r := New(client)
		w, err1 := r.Watch(ctx, s.Name)
		if err1 != nil {
			return
		}
		defer func() {
			_ = w.Stop()
		}()
		for {
			res, err2 := w.Next()
			if err2 != nil {
				return
			}
			t.Logf("watch: %d", len(res))
			for _, r := range res {
				t.Logf("next: %+v", r)
			}
		}
	}()
	time.Sleep(time.Second)

	// new a server
	r := New(client,
		RegisterTTL(2*time.Second),
		MaxRetry(5),
	)

	key := fmt.Sprintf("%s/%s/%s", r.opts.namespace, s.Name, s.ID)
	value, _ := marshal(s)
	r.lease = clientv3.NewLease(r.client)
	leaseID, err := r.registerWithKV(ctx, key, value)
	if err != nil {
		t.Fatal(err)
	}

	// wait for lease expired
	time.Sleep(3 * time.Second)

	res, err := r.GetService(ctx, s.Name)
	if err != nil {
		t.Fatal(err)
	}
	if len(res) != 0 {
		t.Errorf("not expected empty")
	}

	go r.heartBeat(ctx, leaseID, key, value)

	time.Sleep(time.Second)
	res, err = r.GetService(ctx, s.Name)
	if err != nil {
		t.Fatal(err)
	}
	if len(res) == 0 {
		t.Errorf("reconnect failed")
	}
}
func TestRegistry_Watch(t *testing.T) {

	client, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"127.0.0.1:2379"},
		DialTimeout: time.Second, DialOptions: []grpc.DialOption{grpc.WithBlock()},
	})
	if err != nil {
		t.Fatal(err)
	}
	closeClient, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"127.0.0.1:2379"},
		DialTimeout: time.Second, DialOptions: []grpc.DialOption{grpc.WithBlock()},
	})
	if err != nil {
		t.Fatal(err)
	}
	defer closeClient.Close()
	defer client.Close()
	s := &registry.ServiceInstance{
		ID:        "1",
		Name:      "hello",
		Version:   "v1.0.0",
		Endpoints: []string{"127.0.0.1:8080"},
	}

	r := New(client)

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
				serviceName: s.Name,
			},
			wantErr: false,
			want:    []*registry.ServiceInstance{s},
			deferFunc: func(t *testing.T) {
				err = r.Deregister(context.Background(), s)
				if err != nil {
					t.Error(err)
				}
			},
			processFunc: func(t *testing.T, w registry.Watcher) {
				err = r.Register(context.Background(), s)
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
				serviceName: s.Name,
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
				registry: New(closeClient),
			},
			args: args{
				ctx:         context.Background(),
				serviceName: s.Name,
			},
			wantErr: true,
			want:    nil,
			processFunc: func(t *testing.T, w registry.Watcher) {
				closeClient.Close()
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
