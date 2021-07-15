package auth

import (
	"context"
	"fmt"
	"github.com/go-kratos/kratos/v2/middleware"
	"log"
)

func AuthMiddleware(handler middleware.Handler) middleware.Handler {
	return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
		log.Println("auth middleware in", req)
		reply, err = handler(ctx, req)
		fmt.Println("auth middleware out", reply)
		return
	}
}
