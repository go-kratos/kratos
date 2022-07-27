package consul

import (
	"context"
	"fmt"
	"net"
	"reflect"
	"testing"
	"time"

	"github.com/go-kratos/kratos/v2/registry"
	"github.com/hashicorp/consul/api"
)

func tcpServer(t *testing.T, lis net.Listener) {
	for {
		conn, err := lis.Accept()
		if err != nil {
			return
		}
		fmt.Println("get tcp")
		conn.Close()
	}
}

func TestRegistry_Register(t *testing.T) {
	opts := []Option{
		WithHealthCheck(false),
	}

	type args struct {
		ctx        context.Context
		serverName string
		server     []*registry.ServiceInstance
	}

	test := []struct {
		name    string
		args    args
		want    []*registry.ServiceInstance
		wantErr bool
	}{
		{
			name: "normal",
			args: args{
				ctx:        context.Background(),
				serverName: "server-1",
				server: []*registry.ServiceInstance{
					{
						ID:        "1",
						Name:      "server-1",
						Version:   "v0.0.1",
						Metadata:  nil,
						Endpoints: []string{"http://127.0.0.1:8000"},
					},
				},
			},
			want: []*registry.ServiceInstance{
				{
					ID:        "1",
					Name:      "server-1",
					Version:   "v0.0.1",
					Metadata:  nil,
					Endpoints: []string{"http://127.0.0.1:8000"},
				},
			},
			wantErr: false,
		},
		{
			name: "registry new service replace old service",
			args: args{
				ctx:        context.Background(),
				serverName: "server-1",
				server: []*registry.ServiceInstance{
					{
						ID:        "1",
						Name:      "server-1",
						Version:   "v0.0.1",
						Metadata:  nil,
						Endpoints: []string{"http://127.0.0.1:8000"},
					},
					{
						ID:        "1",
						Name:      "server-1",
						Version:   "v0.0.2",
						Metadata:  nil,
						Endpoints: []string{"http://127.0.0.1:8000"},
					},
				},
			},
			want: []*registry.ServiceInstance{
				{
					ID:        "1",
					Name:      "server-1",
					Version:   "v0.0.2",
					Metadata:  nil,
					Endpoints: []string{"http://127.0.0.1:8000"},
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range test {
		t.Run(tt.name, func(t *testing.T) {
			cli, err := api.NewClient(&api.Config{Address: "127.0.0.1:8500"})
			if err != nil {
				t.Fatalf("create consul client failed: %v", err)
			}

			r := New(cli, opts...)

			for _, instance := range tt.args.server {
				err = r.Register(tt.args.ctx, instance)
				if err != nil {
					t.Error(err)
				}
			}

			watch, err := r.Watch(tt.args.ctx, tt.args.serverName)
			if err != nil {
				t.Error(err)
			}
			got, err := watch.Next()

			if (err != nil) != tt.wantErr {
				t.Errorf("GetService() error = %v, wantErr %v", err, tt.wantErr)
				t.Errorf("GetService() got = %v", got)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetService() got = %v, want %v", got, tt.want)
			}

			for _, instance := range tt.args.server {
				_ = r.Deregister(tt.args.ctx, instance)
			}
		})
	}
}

func TestRegistry_GetService(t *testing.T) {
	addr := fmt.Sprintf("%s:9091", getIntranetIP())
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		t.Errorf("listen tcp %s failed!", addr)
		t.Fail()
	}
	defer lis.Close()
	go tcpServer(t, lis)
	time.Sleep(time.Millisecond * 100)
	cli, err := api.NewClient(&api.Config{Address: "127.0.0.1:8500"})
	if err != nil {
		t.Fatalf("create consul client failed: %v", err)
	}
	opts := []Option{
		WithHeartbeat(true),
		WithHealthCheck(true),
		WithHealthCheckInterval(5),
	}
	r := New(cli, opts...)

	instance1 := &registry.ServiceInstance{
		ID:        "1",
		Name:      "server-1",
		Version:   "v0.0.1",
		Endpoints: []string{fmt.Sprintf("tcp://%s?isSecure=false", addr)},
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
			name:   "normal",
			fields: fields{r},
			args: args{
				ctx:         context.Background(),
				serviceName: "server-1",
			},
			want:    []*registry.ServiceInstance{instance1},
			wantErr: false,
			preFunc: func(t *testing.T) {
				if err := r.Register(context.Background(), instance1); err != nil {
					t.Error(err)
				}
				watch, err := r.Watch(context.Background(), instance1.Name)
				if err != nil {
					t.Error(err)
				}
				_, err = watch.Next()
				if err != nil {
					t.Error(err)
				}
			},
			deferFunc: func(t *testing.T) {
				err := r.Deregister(context.Background(), instance1)
				if err != nil {
					t.Error(err)
				}
			},
		},
		{
			name:   "can't get any",
			fields: fields{r},
			args: args{
				ctx:         context.Background(),
				serviceName: "server-x",
			},
			want:    nil,
			wantErr: true,
			preFunc: func(t *testing.T) {
				if err := r.Register(context.Background(), instance1); err != nil {
					t.Error(err)
				}
				watch, err := r.Watch(context.Background(), instance1.Name)
				if err != nil {
					t.Error(err)
				}
				_, err = watch.Next()
				if err != nil {
					t.Error(err)
				}
			},
			deferFunc: func(t *testing.T) {
				err := r.Deregister(context.Background(), instance1)
				if err != nil {
					t.Error(err)
				}
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.preFunc != nil {
				test.preFunc(t)
			}
			if test.deferFunc != nil {
				defer test.deferFunc(t)
			}

			service, err := test.fields.registry.GetService(context.Background(), test.args.serviceName)
			if (err != nil) != test.wantErr {
				t.Errorf("GetService() error = %v, wantErr %v", err, test.wantErr)
				t.Errorf("GetService() got = %v", service)
				return
			}
			if !reflect.DeepEqual(service, test.want) {
				t.Errorf("GetService() got = %v, want %v", service, test.want)
			}
		})
	}
}

func TestRegistry_Watch(t *testing.T) {
	addr := fmt.Sprintf("%s:9091", getIntranetIP())

	time.Sleep(time.Millisecond * 100)
	cli, err := api.NewClient(&api.Config{Address: "127.0.0.1:8500", WaitTime: 2 * time.Second})
	if err != nil {
		t.Fatalf("create consul client failed: %v", err)
	}

	instance1 := &registry.ServiceInstance{
		ID:        "1",
		Name:      "server-1",
		Version:   "v0.0.1",
		Endpoints: []string{fmt.Sprintf("tcp://%s?isSecure=false", addr)},
	}

	type args struct {
		ctx      context.Context
		opts     []Option
		instance *registry.ServiceInstance
	}
	tests := []struct {
		name    string
		args    args
		want    []*registry.ServiceInstance
		wantErr bool
		preFunc func(t *testing.T)
	}{
		{
			name: "normal",
			args: args{
				ctx:      context.Background(),
				instance: instance1,
				opts: []Option{
					WithHealthCheck(false),
				},
			},
			want:    []*registry.ServiceInstance{instance1},
			wantErr: false,
			preFunc: func(t *testing.T) {
			},
		},
		{
			name: "register with healthCheck",
			args: args{
				ctx:      context.Background(),
				instance: instance1,
				opts: []Option{
					WithHeartbeat(true),
					WithHealthCheck(true),
					WithHealthCheckInterval(5),
				},
			},
			want:    []*registry.ServiceInstance{instance1},
			wantErr: false,
			preFunc: func(t *testing.T) {
				lis, err := net.Listen("tcp", addr)
				if err != nil {
					t.Errorf("listen tcp %s failed!", addr)
				}
				go tcpServer(t, lis)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.preFunc != nil {
				tt.preFunc(t)
			}

			r := New(cli, tt.args.opts...)

			err := r.Register(tt.args.ctx, tt.args.instance)
			if err != nil {
				t.Error(err)
			}
			defer func() {
				err = r.Deregister(tt.args.ctx, tt.args.instance)
				if err != nil {
					t.Error(err)
				}
			}()

			watch, err := r.Watch(tt.args.ctx, tt.args.instance.Name)
			if err != nil {
				t.Error(err)
			}

			service, err := watch.Next()

			if (err != nil) != tt.wantErr {
				t.Errorf("GetService() error = %v, wantErr %v", err, tt.wantErr)
				t.Errorf("GetService() got = %v", service)
				return
			}
			if !reflect.DeepEqual(service, tt.want) {
				t.Errorf("GetService() got = %v, want %v", service, tt.want)
			}
		})
	}
}

func getIntranetIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "127.0.0.1"
	}

	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return "127.0.0.1"
}
