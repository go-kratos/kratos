package resolver

import (
	"fmt"
	"testing"
	"time"

	"github.com/bilibili/kratos/pkg/naming"
)

func Test_FilterLittle(t *testing.T) {
	var backs []*naming.Instance
	for i := 0; i < 3; i++ {
		backs = append(backs, &naming.Instance{
			Zone:     "sh1",
			Env:      "prod",
			AppID:    "2233",
			Hostname: fmt.Sprintf("linux-%d", i),
			Addrs:    []string{fmt.Sprintf("grpc://127.0.0.%d:9000", i)},
			LastTs:   time.Now().Unix(),
		})
	}
	r := &Resolver{
		quit:       make(chan struct{}, 1),
		zone:       "sh1",
		subsetSize: 50,
	}

	if len(r.filter(backs)) != 3 {
		t.Fatalf("backends length must be 3")
	}
}

func Test_FilterBig(t *testing.T) {
	var backs []*naming.Instance
	for i := 0; i < 100; i++ {
		backs = append(backs, &naming.Instance{
			Zone:     "sh1",
			Env:      "prod",
			AppID:    "2233",
			Hostname: fmt.Sprintf("linux-%d", i),
			Addrs:    []string{fmt.Sprintf("grpc://127.0.0.%d:9000", i)},
			LastTs:   time.Now().Unix(),
		})
	}
	r := &Resolver{
		quit:       make(chan struct{}, 1),
		zone:       "sh1",
		subsetSize: 50,
	}

	if len(r.filter(backs)) != 50 {
		t.Fatalf("backends length must be 50")
	}
}

func Test_FilterNone(t *testing.T) {
	var backs []*naming.Instance
	for i := 0; i < 100; i++ {
		backs = append(backs, &naming.Instance{
			Zone:     "sh1",
			Env:      "prod",
			AppID:    "2233",
			Metadata: map[string]string{naming.MetaCluster: "c1"},
			Hostname: fmt.Sprintf("linux-%d", i),
			Addrs:    []string{fmt.Sprintf("grpc://127.0.0.%d:9000", i)},
			LastTs:   time.Now().Unix(),
		})
	}
	r := &Resolver{
		quit:       make(chan struct{}, 1),
		zone:       "sh1",
		subsetSize: 50,
		clusters:   map[string]struct{}{"c2": struct{}{}},
	}

	if len(r.filter(backs)) != 0 {
		t.Fatalf("backends length must be 0")
	}
}

func Test_FilterSome(t *testing.T) {
	var backs []*naming.Instance
	for i := 0; i < 40; i++ {
		backs = append(backs, &naming.Instance{
			Zone:     "sh1",
			Env:      "prod",
			AppID:    "2233",
			Metadata: map[string]string{naming.MetaCluster: "c1"},
			Hostname: fmt.Sprintf("linux-%d", i),
			Addrs:    []string{fmt.Sprintf("grpc://127.0.0.%d:9000", i)},
			LastTs:   time.Now().Unix(),
		})
	}
	for i := 50; i < 150; i++ {
		backs = append(backs, &naming.Instance{
			Zone:     "sh1",
			Env:      "prod",
			AppID:    "2233",
			Metadata: map[string]string{naming.MetaCluster: "c2"},
			Hostname: fmt.Sprintf("linux-%d", i),
			Addrs:    []string{fmt.Sprintf("grpc://127.0.0.%d:9000", i)},
			LastTs:   time.Now().Unix(),
		})
	}
	r := &Resolver{
		quit:       make(chan struct{}, 1),
		zone:       "sh1",
		subsetSize: 50,
		clusters:   map[string]struct{}{"c2": struct{}{}},
	}
	if len(r.filter(backs)) != 50 {
		t.Fatalf("backends length must be 0")
	}

	r2 := &Resolver{
		quit:       make(chan struct{}, 1),
		zone:       "sh1",
		subsetSize: 50,
		clusters:   map[string]struct{}{"c1": struct{}{}},
	}
	if len(r2.filter(backs)) != 40 {
		t.Fatalf("backends length must be 0")
	}
}
