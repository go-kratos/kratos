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
						ID:        "2",
						Name:      "server-1",
						Version:   "v0.0.1",
						Metadata:  nil,
						Endpoints: []string{"http://127.0.0.1:8000"},
					},
					{
						ID:        "2",
						Name:      "server-1",
						Version:   "v0.0.2",
						Metadata:  nil,
						Endpoints: []string{"http://127.0.0.1:8000"},
					},
				},
			},
			want: []*registry.ServiceInstance{
				{
					ID:        "2",
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
			watchCtx, watchCancel := context.WithCancel(context.Background())
			watch, err := r.Watch(watchCtx, tt.args.serverName)
			if err != nil {
				t.Error(err)
			}
			got, err := watch.Next()

			if (err != nil) != tt.wantErr {
				t.Errorf("GetService() error = %v, wantErr %v", err, tt.wantErr)
				t.Errorf("GetService() got = %v", got)
				watchCancel()
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetService() got = %v, want %v", got, tt.want)
			}

			for _, instance := range tt.args.server {
				_ = r.Deregister(tt.args.ctx, instance)
			}
			err = watch.Stop()
			if err != nil {
				t.Error(err)
			}
			watchCancel()
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

	instance2 := &registry.ServiceInstance{
		ID:        "2",
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
				watchCtx, watchCancel := context.WithCancel(context.Background())
				watch, err := r.Watch(watchCtx, instance1.Name)
				if err != nil {
					t.Error(err)
				}
				_, err = watch.Next()
				if err != nil {
					t.Error(err)
				}
				err = watch.Stop()
				if err != nil {
					t.Error(err)
				}
				watchCancel()
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
				watchCtx, watchCancel := context.WithCancel(context.Background())
				watch, err := r.Watch(watchCtx, instance2.Name)
				if err != nil {
					t.Error(err)
				}
				_, err = watch.Next()
				if err != nil {
					t.Error(err)
				}
				err = watch.Stop()
				if err != nil {
					t.Error(err)
				}
				watchCancel()
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
		Endpoints: []string{fmt.Sprintf("tcp://%s?isSecure=false", addr)},
	}

	instance2 := &registry.ServiceInstance{
		ID:        "2",
		Name:      "server-1",
		Version:   "v0.0.1",
		Endpoints: []string{fmt.Sprintf("tcp://%s?isSecure=false", addr)},
	}

	instance3 := &registry.ServiceInstance{
		ID:        "3",
		Name:      "server-1",
		Version:   "v0.0.1",
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
				},
			},
			want:    []*registry.ServiceInstance{instance1},
			wantErr: false,
			preFunc: func(*testing.T) {},
		},
		{
			name: "ctx has been canceled",
			args: args{
				ctx:      canceledCtx,
				cancel:   cancel,
				instance: instance2,
				opts: []Option{
					WithHealthCheck(false),
				},
			},
			want:    nil,
			wantErr: true,
			preFunc: func(*testing.T) {},
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
			err = watch.Stop()
			if err != nil {
				t.Error(err)
			}
		})
	}
}

func TestRegistry_IdleAndWatch(t *testing.T) {
	addr := fmt.Sprintf("%s:9091", getIntranetIP())

	time.Sleep(time.Millisecond * 100)
	cli, err := api.NewClient(&api.Config{Address: "127.0.0.1:8500", WaitTime: 2 * time.Second})
	if err != nil {
		t.Fatalf("create consul client failed: %v", err)
	}

	r := New(cli, []Option{
		WithHealthCheck(false),
	}...)

	instance1 := &registry.ServiceInstance{
		ID:        "1",
		Name:      "server-1",
		Version:   "v0.0.1",
		Endpoints: []string{fmt.Sprintf("tcp://%s?isSecure=false", addr)},
	}
	instance2 := &registry.ServiceInstance{
		ID:        "1",
		Name:      "server-1",
		Version:   "v0.0.2",
		Endpoints: []string{fmt.Sprintf("tcp://%s?isSecure=false", addr)},
	}

	type args struct {
		ctx            context.Context
		instance       *registry.ServiceInstance
		changeInstance *registry.ServiceInstance
	}

	tests := []struct {
		name    string
		args    args
		want    []*registry.ServiceInstance
		wantErr bool
	}{
		{
			name: "many client, one idle",
			args: args{
				ctx:            context.Background(),
				instance:       instance1,
				changeInstance: instance2,
			},
			want:    []*registry.ServiceInstance{instance1},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err = r.Register(tt.args.ctx, tt.args.instance)
			if err != nil {
				t.Error(err)
			}
			defer func() {
				err = r.Deregister(tt.args.ctx, tt.args.instance)
				if err != nil {
					t.Error(err)
				}
			}()

			for i := 0; i < 10; i++ {
				watchCtx, watchCancel := context.WithCancel(context.Background())
				watch, err1 := r.Watch(watchCtx, tt.args.instance.Name)
				if err1 != nil {
					t.Error(err1)
				}
				if i != 9 {
					go func(i int) {
						// first
						service, err2 := watch.Next()
						if (err2 != nil) != tt.wantErr {
							t.Errorf("GetService() error = %v, wantErr %v", err, tt.wantErr)
							t.Errorf("GetService() got = %v", service)
							watchCancel()
							return
						}
						// instance changes
						service, err2 = watch.Next()
						if i == 9 {
							return
						}
						if (err2 != nil) != tt.wantErr {
							t.Errorf("GetService() error = %v, wantErr %v", err, tt.wantErr)
							t.Errorf("GetService() got = %v", service)
							watchCancel()
							return
						}
						if !reflect.DeepEqual(service, tt.want) {
							t.Errorf("GetService() got = %v, want %v", service, tt.want)
						}
						err2 = watch.Stop()
						if err2 != nil {
							t.Error(err)
						}
						watchCancel()
						// t.Logf("service:%v, i:%d", service, i)
					}(i)
				} else {
					time.Sleep(time.Second * 3)
					// become idle, close watcher
					err1 = watch.Stop()
					if err1 != nil {
						t.Errorf("watch stop err:%v", err)
					}
					watchCancel()
				}
			}
			time.Sleep(2 * time.Second)
			err = r.Register(tt.args.ctx, tt.args.changeInstance)
			if err != nil {
				t.Error(err)
			}
			time.Sleep(1 * time.Second)
		})
	}
}

func TestRegistry_IdleAndWatch2(t *testing.T) {
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
	}{
		{
			name: "all clients are idle, create a new one",
			args: args{
				ctx:      context.Background(),
				instance: instance1,
				opts: []Option{
					WithHealthCheck(false),
				},
			},
			want:    []*registry.ServiceInstance{instance1},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := New(cli, tt.args.opts...)

			err = r.Register(tt.args.ctx, tt.args.instance)
			if err != nil {
				t.Error(err)
			}
			defer func() {
				err = r.Deregister(tt.args.ctx, tt.args.instance)
				if err != nil {
					t.Error(err)
				}
			}()

			ctx, cancel := context.WithCancel(context.Background())
			for i := 0; i < 10; i++ {
				stopCtx, stopCancel := context.WithCancel(ctx)
				watch, err1 := r.Watch(stopCtx, tt.args.instance.Name)
				if err1 != nil {
					t.Error(err1)
				}
				go func(_ int) {
					// first
					service, err2 := watch.Next()
					if (err2 != nil) != tt.wantErr {
						t.Errorf("GetService() error = %v, wantErr %v", err, tt.wantErr)
						t.Errorf("GetService() got = %v", service)
						return
					}
				}(i)
				go func() {
					select {
					case <-stopCtx.Done():
						err1 = watch.Stop()
						if err1 != nil {
							t.Errorf("watch stop err:%v", err)
						}
						return
					case <-time.After(time.Minute):
						stopCancel()
						err1 = watch.Stop()
						if err1 != nil {
							t.Errorf("watch stop err:%v", err)
						}
						return
					}
				}()
			}
			time.Sleep(time.Second * 3)
			cancel()
			time.Sleep(time.Second * 2)
			// Everything is idle. Add new watch.
			watchCtx, watchCancel := context.WithCancel(context.Background())
			watch, err := r.Watch(watchCtx, tt.args.instance.Name)
			if err != nil {
				t.Error(err)
			}
			service, err := watch.Next()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetService() error = %v, wantErr %v", err, tt.wantErr)
				t.Errorf("GetService() got = %v", service)
				watchCancel()
				return
			}
			if !reflect.DeepEqual(service, tt.want) {
				t.Errorf("GetService() got = %v, want %v", service, tt.want)
			}
			watchCancel()
		})
	}
}

func TestRegistry_ExitOldResolverAndReWatch(t *testing.T) {
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
	instance2 := &registry.ServiceInstance{
		ID:        "2",
		Name:      "server-1",
		Version:   "v0.0.2",
		Endpoints: []string{fmt.Sprintf("tcp://%s?isSecure=false", addr)},
	}
	type args struct {
		ctx             context.Context
		opts            []Option
		instance        *registry.ServiceInstance
		initialInstance *registry.ServiceInstance
	}

	tests := []struct {
		name    string
		args    args
		want    []*registry.ServiceInstance
		wantErr bool
	}{
		{
			name: "When it has entered idle mode, but the old resolver has not completely exited, the watch will be re-established due to new requests coming in.",
			args: args{
				ctx:             context.Background(),
				initialInstance: instance1,
				instance:        instance2,
				opts: []Option{
					WithHealthCheck(false),
					WithTimeout(time.Second * 2),
				},
			},
			want:    []*registry.ServiceInstance{instance2},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := New(cli, tt.args.opts...)

			err = r.Register(tt.args.ctx, tt.args.initialInstance)
			if err != nil {
				t.Error(err)
			}
			// first watch
			ctx, cancel := context.WithCancel(context.Background())
			watch, err := r.Watch(ctx, tt.args.instance.Name)
			if err != nil {
				t.Error(err)
			}
			service, err := watch.Next()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetService() error = %v, wantErr %v", err, tt.wantErr)
				t.Errorf("GetService() got = %v", service)
			}

			time.Sleep(time.Second * 3)
			// The simulation entered idle mode first, but the old resolver was not closed yet, and new requests triggered a new Watch.
			watchCtx := context.Background()
			// old resolver cancel
			err = watch.Stop()
			if err != nil {
				t.Errorf("watch stop err:%v", err)
			}
			cancel()
			// If it sleeps for a period of time, the old resolve goroutine will exit before the new Watch is processed, and there will be no problems at this time.
			// time.Sleep(time.Second * 8)
			newWatch, err := r.Watch(watchCtx, tt.args.instance.Name)
			if err != nil {
				t.Error(err)
			}
			service, err = newWatch.Next()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetService() error = %v, wantErr %v", err, tt.wantErr)
				t.Errorf("GetService() got = %v", service)
			}
			// change register info
			time.Sleep(time.Second * 1)
			err = r.Deregister(tt.args.ctx, tt.args.initialInstance)
			if err != nil {
				t.Error(err)
			}
			time.Sleep(time.Second * 5)
			err = r.Register(tt.args.ctx, tt.args.instance)
			if err != nil {
				t.Error(err)
			}

			time.Sleep(time.Second * 2)

			newWatchCtx, newWatchCancel := context.WithCancel(context.Background())
			c := make(chan struct{}, 1)

			go func() {
				service, err = newWatch.Next()
				if (err != nil) != tt.wantErr {
					t.Errorf("GetService() error = %v, wantErr %v", err, tt.wantErr)
					t.Errorf("GetService() got = %v", service)
					return
				}
				if !reflect.DeepEqual(service, tt.want) {
					t.Errorf("GetService() got = %v, want %v", service, tt.want)
				}
				c <- struct{}{}
			}()
			time.AfterFunc(time.Second*10, newWatchCancel)
			select {
			case <-newWatchCtx.Done():
				t.Errorf("Timeout getservice. May be no new resolve goroutine to obtain the latest service information")
			case <-c:
				return
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
