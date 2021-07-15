package log

import (
	"context"
	"fmt"
	"github.com/go-kratos/kratos/v2/middleware"
	"log"
	"time"
)

func LogMiddleware(handler middleware.Handler) middleware.Handler {
	return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
		log.Println("log middleware in", req)
		reply, err = handler(ctx, req)
		fmt.Println("log middleware out", reply)
		return
	}
}


func TimeMiddleware(handler middleware.Handler) middleware.Handler {
	return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
		log.Println("time middleware in", req)
		start := time.Now()
		reply, err = handler(ctx, req)
		fmt.Println("time middleware out", reply)
		fmt.Println(time.Since(start))
		return
	}
}
