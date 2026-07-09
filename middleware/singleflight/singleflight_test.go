package singleflight

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/go-kratos/kratos/v2/transport/grpc"
)

type testVali struct {
	in  string
	out int
}

type Transport2 struct {
	grpc.Transport
	op string
}

func (t Transport2) Operation() string {
	return t.op
}

// 测试使用单飞时
func TestUse(t *testing.T) {
	var mu sync.Mutex
	var callNum int

	var mock middleware.Handler = func(ctx context.Context, req interface{}) (interface{}, error) {
		in := req.(testVali)
		mu.Lock()
		callNum++
		mu.Unlock()
		time.Sleep(1 * time.Second)
		return in.out, nil
	}

	tests := []testVali{
		{"1", 1},
		{"2", 2},
		{"2", 2},
		{"3", 3},
		{"3", 3},
		{"3", 3},
	}
	var wg sync.WaitGroup
	for _, test := range tests {
		wg.Add(1)
		go func(te testVali) {
			t.Run(te.in, func(t *testing.T) {
				v := SingleFlight("test")(mock) //注册
				tr := &Transport2{op: "test"}
				ctx := transport.NewServerContext(context.Background(), tr)
				re, err := v(ctx, te)
				if err != nil {
					t.Error(err)
				}
				if re != te.out {
					t.Errorf("err: %v", te)
				}
				wg.Done()
			})
		}(test)
	}

	wg.Wait()

	//最后计算总调用次数
	t.Run("callNum", func(t *testing.T) {
		if callNum != 3 {
			t.Errorf("callNum err: %v", callNum)
		}
	})
}

// 测试不使用单飞时
func TestNoUse(t *testing.T) {
	var mu sync.Mutex
	var callNum int

	var mock middleware.Handler = func(ctx context.Context, req interface{}) (interface{}, error) {
		in := req.(testVali)
		mu.Lock()
		callNum++
		mu.Unlock()
		time.Sleep(1 * time.Second)
		return in.out, nil
	}

	tests := []testVali{
		{"1", 1},
		{"2", 2},
		{"2", 2},
		{"3", 3},
		{"3", 3},
		{"3", 3},
	}
	var wg sync.WaitGroup
	for _, test := range tests {
		wg.Add(1)
		go func(te testVali) {
			t.Run(te.in, func(t *testing.T) {
				v := SingleFlight()(mock) //移除注册
				tr := &Transport2{op: "test"}
				ctx := transport.NewServerContext(context.Background(), tr)
				re, err := v(ctx, te)
				if err != nil {
					t.Error(err)
				}
				if re != te.out {
					t.Errorf("err: %v", te)
				}
				wg.Done()
			})
		}(test)
	}

	wg.Wait()

	//最后计算总调用次数
	t.Run("callNum", func(t *testing.T) {
		if callNum != 6 {
			t.Errorf("callNum err: %v", callNum)
		}
	})
}
