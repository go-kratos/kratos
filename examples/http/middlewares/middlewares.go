package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-kratos/kratos/v2/middleware"
)

func globalFilter(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("global filter in")
		next.ServeHTTP(w, r)
		fmt.Println("global filter out")
	})
}

func globalFilter2(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("global filter 2 in")
		next.ServeHTTP(w, r)
		fmt.Println("global filter 2 out")
	})
}

func routeFilter(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("route filter in")
		next.ServeHTTP(w, r)
		fmt.Println("route filter out")
	})
}

func routeFilter2(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("route filter 2 in")
		next.ServeHTTP(w, r)
		fmt.Println("route filter 2 out")
	})
}

func pathFilter(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("path filter in")
		next.ServeHTTP(w, r)
		fmt.Println("path filter out")
	})
}

func pathFilter2(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("path filter 2 in")
		next.ServeHTTP(w, r)
		fmt.Println("path filter 2 out")
	})
}

func serviceMiddleware(handler middleware.Handler) middleware.Handler {
	return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
		fmt.Println("service middleware in")
		reply, err = handler(ctx, req)
		fmt.Println("service middleware out")
		return
	}
}

func serviceMiddleware2(handler middleware.Handler) middleware.Handler {
	return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
		fmt.Println("service middleware 2 in")
		reply, err = handler(ctx, req)
		fmt.Println("service middleware 2 out")
		return
	}
}
