package eureka

import (
	"context"
	"fmt"
	"log"
	"sync"
	"testing"
	"time"

	"github.com/go-kratos/kratos/v2/registry"
)

func TestRegistry(_ *testing.T) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	s1 := &registry.ServiceInstance{
		ID:        "0",
		Name:      "helloworld",
		Endpoints: []string{"http://127.0.0.1:1111"},
	}
	s2 := &registry.ServiceInstance{
		ID:        "0",
		Name:      "helloworld2",
		Endpoints: []string{"http://127.0.0.1:222"},
	}

	r, _ := New([]string{"https://127.0.0.1:18761"}, WithContext(ctx), WithHeartbeat(time.Second), WithRefresh(time.Second), WithEurekaPath("eureka"))

	go do(r, s1)
	go do(r, s2)

	time.Sleep(time.Second * 20)
	cancel()
	time.Sleep(time.Second * 1)
}

func do(r *Registry, s *registry.ServiceInstance) {
	w, err := r.Watch(context.Background(), s.Name)
	if err != nil {
		log.Fatal(err)
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
			log.Printf("watch: %d", len(res))
			for _, r := range res {
				log.Printf("next: %+v", r)
			}
		}
	}()

	ctx, cancel := context.WithCancel(context.Background())
	if err = r.Register(ctx, s); err != nil {
		log.Fatal(err)
	}

	time.Sleep(time.Second * 10)
	res, err := r.GetService(ctx, s.Name)
	if err != nil {
		log.Fatal(err)
	}
	for i, re := range res {
		log.Printf("first %d re:%v\n", i, re)
	}

	if len(res) != 1 && res[0].Name != s.Name {
		log.Fatalf("not expected: %+v", res)
	}

	if err = r.Deregister(ctx, s); err != nil {
		log.Fatal(err)
	}
	cancel()
	time.Sleep(time.Second * 10)

	res, err = r.GetService(ctx, s.Name)
	if err != nil {
		log.Fatal(err)
	}
	for i, re := range res {
		log.Printf("second %d re:%v\n", i, re)
	}
	if len(res) != 0 {
		log.Fatalf("not expected empty")
	}
}

func TestLock(_ *testing.T) {
	type me struct {
		lock sync.Mutex
	}

	a := &me{}
	go func() {
		defer a.lock.Unlock()
		a.lock.Lock()
		fmt.Println("This is fmt first.")
		time.Sleep(time.Second * 5)
	}()
	go func() {
		defer a.lock.Unlock()
		a.lock.Lock()
		fmt.Println("This is fmt second.")
		time.Sleep(time.Second * 5)
	}()
	time.Sleep(time.Second * 10)
}
