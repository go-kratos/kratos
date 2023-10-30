package consul

import (
	"context"
	"fmt"
	"net"
	"reflect"
	"testing"
	"time"

	"github.com/hashicorp/consul/api"

	"github.com/go-kratos/kratos/v2/registry"
)

func tcpServer(lis net.Listener) {
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
		WithHeartbeat(false),
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
						Metadata:  map[string]string{"cluster": "dc1"},
						Endpoints: []string{"http://127.0.0.1:8000"},
					},
				},
			},
			want: []*registry.ServiceInstance{
				{
					ID:        "1",
					Name:      "server-1",
					Version:   "v0.0.1",
					Metadata:  map[string]string{"cluster": "dc1"},
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
						ID:        "2",
						Name:      "server-1",
						Version:   "v0.0.1",
						Metadata:  map[string]string{"cluster": "dc1"},
						Endpoints: []string{"http://127.0.0.1:8000"},
					},
					{
						ID:        "2",
						Name:      "server-1",
						Version:   "v0.0.2",
						Metadata:  map[string]string{"cluster": "dc1"},
						Endpoints: []string{"http://127.0.0.1:8000"},
					},
				},
			},
			want: []*registry.ServiceInstance{
				{
					ID:        "2",
					Name:      "server-1",
					Version:   "v0.0.2",
					Metadata:  map[string]string{"cluster": "dc1"},
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
	go tcpServer(lis)
	time.Sleep(time.Millisecond * 100)
	cli, err := api.NewClient(&api.Config{Address: "127.0.0.1:8500"})
	if err != nil {
		t.Fatalf("create consul client failed: %v", err)
	}
	opts := []Option{
		WithHeartbeat(false),
		WithHealthCheck(false),
	}
	r := New(cli, opts...)

	instance1 := &registry.ServiceInstance{
		ID:        "1",
		Name:      "server-1",
		Version:   "v0.0.1",
		Metadata:  map[string]string{"cluster": "dc1"},
		Endpoints: []string{fmt.Sprintf("tcp://%s?isSecure=false", addr)},
	}

	instance2 := &registry.ServiceInstance{
		ID:        "2",
		Name:      "server-1",
		Version:   "v0.0.1",
		Metadata:  map[string]string{"cluster": "dc1"},
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
				if err := r.Register(context.Background(), instance2); err != nil {
					t.Error(err)
				}
				watch, err := r.Watch(context.Background(), instance2.Name)
				if err != nil {
					t.Error(err)
				}
				_, err = watch.Next()
				if err != nil {
					t.Error(err)
				}
			},
			deferFunc: func(t *testing.T) {
				err := r.Deregister(context.Background(), instance2)
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
		Metadata:  map[string]string{"cluster": "dc1"},
		Endpoints: []string{fmt.Sprintf("tcp://%s?isSecure=false", addr)},
	}

	instance2 := &registry.ServiceInstance{
		ID:        "2",
		Name:      "server-1",
		Version:   "v0.0.1",
		Metadata:  map[string]string{"cluster": "dc1"},
		Endpoints: []string{fmt.Sprintf("tcp://%s?isSecure=false", addr)},
	}

	instance3 := &registry.ServiceInstance{
		ID:        "3",
		Name:      "server-1",
		Version:   "v0.0.1",
		Metadata:  map[string]string{"cluster": "dc1"},
		Endpoints: []string{fmt.Sprintf("tcp://%s?isSecure=false", addr)},
	}

	type args struct {
		ctx      context.Context
		cancel   func()
		opts     []Option
		instance *registry.ServiceInstance
	}
	canceledCtx, cancel := context.WithCancel(context.Background())

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
					WithHeartbeat(false),
				},
			},
			want:    []*registry.ServiceInstance{instance1},
			wantErr: false,
			preFunc: func(t *testing.T) {
			},
		},
		{
			name: "ctx has been cancelled",
			args: args{
				ctx:      canceledCtx,
				cancel:   cancel,
				instance: instance2,
				opts: []Option{
					WithHealthCheck(false),
					WithHeartbeat(false),
				},
			},
			want:    nil,
			wantErr: true,
			preFunc: func(t *testing.T) {
			},
		},
		{
			name: "register with healthCheck",
			args: args{
				ctx:      context.Background(),
				instance: instance3,
				opts: []Option{
					WithHeartbeat(true),
					WithHealthCheck(true),
					WithHealthCheckInterval(5),
				},
			},
			want:    []*registry.ServiceInstance{instance3},
			wantErr: false,
			preFunc: func(t *testing.T) {
				lis, err := net.Listen("tcp", addr)
				if err != nil {
					t.Errorf("listen tcp %s failed!", addr)
					return
				}
				go tcpServer(lis)
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

			if tt.args.cancel != nil {
				tt.args.cancel()
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

func TestEstablishPeering(t *testing.T) {
	cluster1, err := api.NewClient(&api.Config{Address: "127.0.0.1:8500", WaitTime: 2 * time.Second})
	if err != nil {
		t.Fatalf("create consul client failed: %v", err)
	}
	res, _, err := cluster1.Peerings().GenerateToken(context.Background(), api.PeeringGenerateTokenRequest{PeerName: "cluster02"}, nil)
	if err != nil {
		t.Fatal(err)
	}
	cluster2, err := api.NewClient(&api.Config{Address: "127.0.0.1:8501", WaitTime: 2 * time.Second})
	if err != nil {
		t.Fatalf("create consul client failed: %v", err)
	}
	establish, _, err := cluster2.Peerings().Establish(
		context.Background(),
		api.PeeringEstablishRequest{
			PeerName:     "cluster01",
			PeeringToken: res.PeeringToken,
		}, nil)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(establish)

	peerings, _, err := cluster1.Peerings().List(context.Background(), nil)
	if err != nil {
		t.Fatal(err)
	}

	if len(peerings) != 1 {
		t.Fatal("peerings len != 1")
	}

	// cluster01 -> cluster 02
	// cluster02 x<-x cluster01
	ok, _, err := cluster1.ConfigEntries().Set(
		&api.ExportedServicesConfigEntry{
			Name: "default",
			Services: []api.ExportedService{
				{
					Name: "*",
					Consumers: []api.ServiceConsumer{
						{Peer: "cluster02"},
					},
				},
			},
		}, nil)
	if err != nil || !ok {
		t.Fatal(err)
	}
}

func TestPeeringGetService(t *testing.T) {
	cluster1, err := api.NewClient(&api.Config{Address: "127.0.0.1:8500", WaitTime: 2 * time.Second})
	if err != nil {
		t.Fatalf("create consul client failed: %v", err)
	}
	cluster2, err := api.NewClient(&api.Config{Address: "127.0.0.1:8501", WaitTime: 2 * time.Second})
	if err != nil {
		t.Fatalf("create consul client failed: %v", err)
	}

	c1 := New(cluster1, WithMultiClusterMode(Peering), WithHealthCheck(false), WithHeartbeat(false))
	c2 := New(cluster2, WithMultiClusterMode(Peering), WithHealthCheck(false), WithHeartbeat(false))

	for i := 0; i < 5; i++ {
		id := fmt.Sprintf("ci-test-%d", i)
		if err = c1.Register(context.Background(), &registry.ServiceInstance{
			Name: "ci-test",
			ID:   id,
			Endpoints: []string{
				"grpc://123.123.123.123",
			},
		}); err != nil {
			t.Fatal(err)
		}

		if err = c2.Register(context.Background(), &registry.ServiceInstance{
			Name: "ci-test",
			ID:   id,
			Endpoints: []string{
				"grpc://123.123.123.123",
			},
		}); err != nil {
			t.Fatal(err)
		}

		// cluster01 - want len 5, cluster want len 10
		t.Cleanup(func() {
			_ = c1.Deregister(context.Background(), &registry.ServiceInstance{
				ID: id,
			})
			_ = c2.Deregister(context.Background(), &registry.ServiceInstance{
				ID: id,
			})
		})
	}

	cluster1Services, err := c1.GetService(context.Background(), "ci-test")
	if err != nil {
		t.Fatal(err)
	}

	if len(cluster1Services) != 5 {
		t.Fatal("cluster1Services len != 5")
	}

	cluster2Services, err := c2.GetService(context.Background(), "ci-test")
	if err != nil {
		t.Fatal(err)
	}

	if len(cluster2Services) != 10 {
		t.Fatal("cluster2 get service len != 10")
	}
}

func TestPeeringWatch(t *testing.T) {
	cluster1, err := api.NewClient(&api.Config{Address: "127.0.0.1:8500", WaitTime: 2 * time.Second})
	if err != nil {
		t.Fatalf("create consul client failed: %v", err)
	}
	cluster2, err := api.NewClient(&api.Config{Address: "127.0.0.1:8501", WaitTime: 2 * time.Second})
	if err != nil {
		t.Fatalf("create consul client failed: %v", err)
	}

	c1 := New(cluster1, WithMultiClusterMode(Peering), WithHealthCheck(false), WithHeartbeat(false))
	c2 := New(cluster2, WithMultiClusterMode(Peering), WithHealthCheck(false), WithHeartbeat(false))

	cw1, err := c1.Watch(context.Background(), "ci-test")
	if err != nil {
		t.Fatal(err)
	}

	cw2, err := c2.Watch(context.Background(), "ci-test")
	if err != nil {
		t.Fatal(err)
	}

	for i := 0; i < 5; i++ {
		id := fmt.Sprintf("ci-test-%d", i)

		err = c1.Register(context.Background(), &registry.ServiceInstance{
			ID: fmt.Sprintf("ci-test-%d", i),
			Endpoints: []string{
				"grpc://123.123.123.123",
			},
			Name: "ci-test",
		})

		if err != nil {
			t.Fatal(err)
		}

		err = c2.Register(context.Background(), &registry.ServiceInstance{
			ID: fmt.Sprintf("ci-test-%d", i),
			Endpoints: []string{
				"grpc://123.123.123.123",
			},
			Name: "ci-test",
		})

		if err != nil {
			t.Fatal(err)
		}

		t.Cleanup(func() {
			_ = c1.Deregister(context.Background(), &registry.ServiceInstance{
				ID: id,
			})
			_ = c2.Deregister(context.Background(), &registry.ServiceInstance{
				ID: id,
			})
		})

		time.Sleep(time.Second * 2)

		var res []*registry.ServiceInstance
		if res, err = cw1.Next(); err != nil || len(res) != i+1 {
			t.Errorf("cluster1 watch failed, len %d != %d or err=%v", len(res), i+1, err)
		}

		if res, err = cw2.Next(); err != nil || len(res) != (i+1)*2 {
			t.Errorf("cluster2 watch failed, len %d != %d or err=%v", len(res), (i+1)*2, err)
		}
	}

	// when the obtained instance is 0, the broadcast will not be triggered to prevent all nodes from being removed, so there is one less test here
	for i := 4; i > 0; i-- {
		err = c1.Deregister(context.Background(), &registry.ServiceInstance{
			ID: fmt.Sprintf("ci-test-%d", i-1),
			Endpoints: []string{
				"grpc://123.123.123.123",
			},
			Name: "ci-test",
		})

		if err != nil {
			t.Fatal(err)
		}

		err = c2.Deregister(context.Background(), &registry.ServiceInstance{
			ID: fmt.Sprintf("ci-test-%d", i-1),
			Endpoints: []string{
				"grpc://123.123.123.123",
			},
			Name: "ci-test",
		})

		if err != nil {
			t.Fatal(err)
		}

		time.Sleep(time.Second * 2)

		var res []*registry.ServiceInstance
		if res, err = cw1.Next(); err != nil || len(res) != i {
			t.Errorf("cluster1 watch failed, len %d != %d or err=%v", len(res), i, err)
		}

		if res, err = cw2.Next(); err != nil || len(res) != i*2 {
			t.Errorf("cluster2 watch failed, len %d != %d or err=%v", len(res), i*2, err)
		}
	}
}
