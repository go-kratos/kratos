package zookeeper

import (
	"context"
	"testing"
	"time"

	"github.com/go-kratos/kratos/v2/registry"
	"github.com/go-zookeeper/zk"
)

func TestRegistry(t *testing.T) {
	ctx := context.Background()
	s := &registry.ServiceInstance{
		ID:        "0",
		Name:      "helloworld",
		Endpoints: []string{"http://127.0.0.1:1111"},
	}

	conn, _, err := zk.Connect([]string{"127.0.0.1:2181"}, time.Second*15)
	if err != nil {
		t.Fatal(err)
		return
	}
	r := New(conn)
	err = r.Register(ctx, s)
	if err != nil {
		t.Fatal(err)
		return
	}
	time.Sleep(time.Second)

	w, err := r.Watch(ctx, s.Name)
	if err != nil {
		t.Fatal(err)
		return
	}
	defer func() {
		_ = w.Stop()
	}()
	go func() {
		for {
			res, nextErr := w.Next()
			if nextErr != nil {
				t.Errorf("watch next error: %s", nextErr)
				return
			}
			t.Logf("watch: %d", len(res))
			for _, r := range res {
				t.Logf("next: %+v", r)
			}
		}
	}()
	time.Sleep(time.Second)

	if err = r.Register(ctx, s); err != nil {
		t.Fatal(err)
	}
	time.Sleep(time.Second)

	res, err := r.GetService(ctx, s.Name)
	if err != nil {
		t.Fatal(err)
	}
	for i, re := range res {
		t.Logf("first %d re:%v\n", i, re)
	}
	if len(res) != 1 && res[0].Name != s.Name {
		t.Errorf("not expected: %+v", res)
	}

	if err = r.Deregister(ctx, s); err != nil {
		t.Fatal(err)
	}
	time.Sleep(time.Second)

	res, err = r.GetService(ctx, s.Name)
	if err != nil {
		t.Fatal(err)
	}
	for i, re := range res {
		t.Logf("second %d re:%v\n", i, re)
	}
	if len(res) != 0 {
		t.Errorf("not expected empty")
	}
}
