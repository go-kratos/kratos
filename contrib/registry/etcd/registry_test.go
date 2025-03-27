package etcd

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"google.golang.org/grpc"

	clientv3 "go.etcd.io/etcd/client/v3"

	"github.com/go-kratos/kratos/v2/registry"
)

func TestRegistry(t *testing.T) {
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

	r := New(client)
	w, err := r.Watch(ctx, s.Name)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		_ = w.Stop()
	}()
	go func() {
		for {
			res, err1 := w.Next()
			if err1 != nil {
				return
			}
			t.Logf("watch: %d", len(res))
			for _, r := range res {
				t.Logf("next: %+v", r)
			}
		}
	}()
	time.Sleep(time.Second)

	if err1 := r.Register(ctx, s); err1 != nil {
		t.Fatal(err1)
	}
	time.Sleep(time.Second)

	res, err := r.GetService(ctx, s.Name)
	if err != nil {
		t.Fatal(err)
	}
	if len(res) != 1 && res[0].Name != s.Name {
		t.Errorf("not expected: %+v", res)
	}

	if err1 := r.Deregister(ctx, s); err1 != nil {
		t.Fatal(err1)
	}
	time.Sleep(time.Second)

	res, err = r.GetService(ctx, s.Name)
	if err != nil {
		t.Fatal(err)
	}
	if len(res) != 0 {
		t.Errorf("not expected empty")
	}
}

func TestConcurrentRegistry(t *testing.T) {
	client, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"127.0.0.1:2379"},
		DialTimeout: time.Second, DialOptions: []grpc.DialOption{grpc.WithBlock()},
	})
	if err != nil {
		t.Fatal(err)
	}
	defer client.Close()

	ctx := context.Background()
	serviceInstances := make([]*registry.ServiceInstance, 0)
	for i := 0; i < 10; i++ {
		serviceInstances = append(serviceInstances, &registry.ServiceInstance{
			ID:   fmt.Sprintf("%d", i),
			Name: fmt.Sprintf("helloworld_%d", i),
		})
	}

	r := New(client)
	watchers := make([]registry.Watcher, 0, len(serviceInstances))
	for _, s := range serviceInstances {
		w, err := r.Watch(ctx, s.Name)
		if err != nil {
			t.Fatal(err)
		}
		watchers = append(watchers, w)
		go func() {
			for {
				res, err1 := w.Next()
				if err1 != nil {
					return
				}
				t.Logf("watch: %d", len(res))
				for _, r := range res {
					t.Logf("next: %+v", r)
				}
			}
		}()
	}
	defer func() {
		for _, w := range watchers {
			_ = w.Stop()
		}
	}()

	time.Sleep(time.Second)

	var wg sync.WaitGroup
	for _, s := range serviceInstances {
		wg.Add(1)
		go func(s *registry.ServiceInstance) {
			defer wg.Done()
			if err1 := r.Register(ctx, s); err1 != nil {
				t.Error(err1)
				return
			}
		}(s)
	}
	time.Sleep(time.Second)

	for _, s := range serviceInstances {
		res, err := r.GetService(ctx, s.Name)
		if err != nil {
			t.Fatal(err)
		}
		if len(res) != 1 && res[0].Name != s.Name {
			t.Errorf("not expected: %+v", res)
		}
	}

	for _, s := range serviceInstances {
		wg.Add(1)
		go func(s *registry.ServiceInstance) {
			defer wg.Done()
			if err1 := r.Deregister(ctx, s); err1 != nil {
				t.Error(err1)
				return
			}
		}(s)
	}
	time.Sleep(time.Second)

	for _, s := range serviceInstances {
		res, err := r.GetService(ctx, s.Name)
		if err != nil {
			t.Fatal(err)
		}
		if len(res) != 0 {
			t.Errorf("not expected empty")
		}
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
