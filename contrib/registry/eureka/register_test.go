package eureka

import (
	"context"
	"fmt"
	"log"
	"sync"
	"testing"
	"time"

	"github.com/go-kratos/kratos/v3/registry"
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
		log.Fatalf("Failed to watch service %q: %v", s.Name, err)
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
		log.Fatalf("Failed to register service %q: %v", s.Name, err)
	}

	time.Sleep(time.Second * 10)
	res, err := r.GetService(ctx, s.Name)
	if err != nil {
		log.Fatalf("Failed to get service %q: %v", s.Name, err)
	}
	for i, re := range res {
		log.Printf("first %d re:%v\n", i, re)
	}

	if len(res) != 1 && res[0].Name != s.Name {
		log.Fatalf("not expected: %+v", res)
	}

	if err = r.Deregister(ctx, s); err != nil {
		log.Fatalf("Failed to deregister service %q: %v", s.Name, err)
	}
	cancel()
	time.Sleep(time.Second * 10)

	res, err = r.GetService(ctx, s.Name)
	if err != nil {
		log.Fatalf("Failed to get service %q after deregister: %v", s.Name, err)
	}
	for i, re := range res {
		log.Printf("second %d re:%v\n", i, re)
	}
	if len(res) != 0 {
		log.Fatalf("not expected empty")
	}
}

func TestEndpoints_MetadataNotShared(t *testing.T) {
	r, err := New([]string{"http://127.0.0.1:8761"})
	if err != nil {
		t.Fatal(err)
	}

	svc := &registry.ServiceInstance{
		ID:        "test-id",
		Name:      "test-svc",
		Version:   "v1.0.0",
		Endpoints: []string{"grpc://192.168.1.100:9000", "http://192.168.1.100:8000"},
		Metadata:  map[string]string{"weight": "10"},
	}

	eps := r.Endpoints(svc)
	if len(eps) != 2 {
		t.Fatalf("expected 2 endpoints, got %d", len(eps))
	}

	// Each endpoint must have the correct "Endpoints" value
	wantEndpoints := []string{"grpc://192.168.1.100:9000", "http://192.168.1.100:8000"}
	for i, ep := range eps {
		if ep.MetaData["Endpoints"] != wantEndpoints[i] {
			t.Logf("endpoint[%d] Endpoints = %q, want %q", i, ep.MetaData["Endpoints"], wantEndpoints[i])
			t.Errorf("endpoint[%d] Endpoints = %q, want %q", i, ep.MetaData["Endpoints"], wantEndpoints[i])
		}
	}

	// Mutating one metadata must not affect the other
	eps[0].MetaData["weight"] = "999"
	if eps[1].MetaData["weight"] != "10" {
		t.Errorf("metadata map is shared between endpoints: endpoint[1].weight = %q, want %q", eps[1].MetaData["weight"], "10")
	}
}

func TestEndpoints_EmptyMetadata(t *testing.T) {
	r, err := New([]string{"http://127.0.0.1:8761"})
	if err != nil {
		t.Fatal(err)
	}

	svc := &registry.ServiceInstance{
		ID:        "test-id",
		Name:      "test-svc",
		Version:   "v1.0.0",
		Endpoints: []string{"grpc://192.168.1.100:9000", "http://192.168.1.100:8000"},
	}

	eps := r.Endpoints(svc)
	if len(eps) != 2 {
		t.Fatalf("expected 2 endpoints, got %d", len(eps))
	}
	for i, ep := range eps {
		if ep.MetaData["Endpoints"] == "" {
			t.Errorf("endpoint[%d] Endpoints is empty", i)
		}
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
