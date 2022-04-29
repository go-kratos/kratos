package zookeeper

import (
	"context"
	"testing"
	"time"

	"github.com/go-kratos/kratos/v2/registry"
)

func TestRegistry(t *testing.T) {
	ctx := context.Background()
	s := &registry.ServiceInstance{
		ID:        "0",
		Name:      "helloworld",
		Endpoints: []string{"http://127.0.0.1:1111"},
	}

	r, _ := New([]string{"127.0.0.1:2181"}, WithDigestACL("username", "password"))

	w, err := r.Watch(ctx, s.Name)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		_ = w.Stop()
	}()
	go func() {
		for {
			res, nextErr := w.Next()
			if nextErr != nil {
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
