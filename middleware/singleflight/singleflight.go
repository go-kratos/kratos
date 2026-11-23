package singleflight

import (
	"context"
	"fmt"

	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
	"golang.org/x/sync/singleflight"
)

var singleflightGroup singleflight.Group

// 单飞 middleware.
/*
	//服务端下使用，通过传入op名称，使用单飞：
	singleflight.SingleFlight(
			"/service.test1/GetCityName",
			"/service.test2/GetAllCityName",
			)
*/
func SingleFlight(ops ...string) middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
			if tr, ok := transport.FromServerContext(ctx); ok {
				if Contains(ops, tr.Operation()) {
					cacheKey := fmt.Sprintf("%s %s", tr.Operation(), req)
					reply, err, _ = singleflightGroup.Do(cacheKey, func() (interface{}, error) {
						return handler(ctx, req)
					})
					return reply, err
				}
			}
			return handler(ctx, req)
		}
	}
}

func Contains(elems []string, elem string) bool {
	for _, e := range elems {
		if elem == e {
			return true
		}
	}
	return false
}
