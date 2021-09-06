package middleware

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

var i int

func TestChain(t *testing.T) {
	next := func(ctx context.Context, req interface{}) (interface{}, error) {
		t.Log(req)
		i += 10
		return "reply", nil
	}

	got, err := Chain(test1Middleware, test2Middleware, test3Middleware)(next)(context.Background(), "hello kratos!")
	assert.Nil(t, err)
	assert.Equal(t, got, "reply")
	assert.Equal(t, i, 16)
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
