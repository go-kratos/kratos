package main

import (
	"context"
	"fmt"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/go-kratos/kratos/v2/transport/http"
)

func serviceMiddleware(handler middleware.Handler) middleware.Handler {
	return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
		var (
			kind      string
			operation string
		)
		if info, ok := transport.FromServerContext(ctx); ok {
			kind = info.Kind().String()
			operation = info.Operation()
			fmt.Println("service middleware in", req)
			// You can assert info as *http.Transport/*grpc.Transport
			if kind == "http" {
				if ht, ok := info.(*http.Transport); ok {
					// You can then use the original *http.Request
					host := ht.Request().Host
					fmt.Printf("host: %s, kind: %s, operation: %s\n", host, kind, operation)
				}
			}
		}
		reply, err = handler(ctx, req)
		fmt.Println("service middleware out", reply)
		return
	}
}

func serviceMiddleware2(handler middleware.Handler) middleware.Handler {
	return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
		fmt.Println("service middleware 2 in", req)
		reply, err = handler(ctx, req)
		fmt.Println("service middleware 2 out", reply)
		return
	}
}
