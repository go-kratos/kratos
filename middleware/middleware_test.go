package middleware

import (
	"context"
	"fmt"
	"reflect"
	"testing"
)

var i int

func TestChain(t *testing.T) {
	next := func(_ context.Context, req interface{}) (interface{}, error) {
		if req != "hello kratos!" {
			t.Errorf("expect %v, got %v", "hello kratos!", req)
		}
		i += 10
		return "reply", nil
	}

	got, err := Chain(test1Middleware, test2Middleware, test3Middleware)(next)(context.Background(), "hello kratos!")
	if err != nil {
		t.Errorf("expect %v, got %v", nil, err)
	}
	if !reflect.DeepEqual(got, "reply") {
		t.Errorf("expect %v, got %v", "reply", got)
	}
	if !reflect.DeepEqual(i, 16) {
		t.Errorf("expect %v, got %v", 16, i)
	}
}

func test1Middleware(handler Handler) Handler {
	return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
		fmt.Println("test1 before")
		i++
		reply, err = handler(ctx, req)
		fmt.Println("test1 after")
		return
	}
}

func test2Middleware(handler Handler) Handler {
	return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
		fmt.Println("test2 before")
		i += 2
		reply, err = handler(ctx, req)
		fmt.Println("test2 after")
		return
	}
}

func test3Middleware(handler Handler) Handler {
	return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
		fmt.Println("test3 before")
		i += 3
		reply, err = handler(ctx, req)
		fmt.Println("test3 after")
		return
	}
}
