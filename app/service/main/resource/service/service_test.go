package service

import (
	"context"
	"flag"
	"path/filepath"
	"time"

	"go-common/app/service/main/resource/conf"
	"go-common/library/cache/redis"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	s *Service
)

func CleanCache() {
	c := context.Background()
	pool := redis.NewPool(conf.Conf.Redis.Ads.Config)
	pool.Get(c).Do("FLUSHDB")
}

func init() {
	dir, _ := filepath.Abs("../cmd/resource-service-test.toml")
	flag.Set("conf", dir)
	conf.Init()
	s = New(conf.Conf)
	time.Sleep(time.Second)
}

func WithService(f func(s *Service)) func() {
	return func() {
		Reset(func() { CleanCache() })
		f(s)
	}
}
