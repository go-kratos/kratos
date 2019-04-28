package service

import (
	"context"
	"flag"
	"path/filepath"
	"time"

	"go-common/app/service/main/thumbup/conf"
	"go-common/library/cache/redis"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	svr *Service
	c   = context.Background()
)

func CleanCache() {
	pool := redis.NewPool(conf.Conf.Redis.Config)
	pool.Get(c).Do("FLUSHDB")
}

func init() {
	dir, _ := filepath.Abs("../cmd/thumbup-test.toml")
	flag.Set("conf", dir)
	conf.Init()
	svr = New(conf.Conf)
	time.Sleep(time.Second)
}
func WithService(f func(s *Service)) func() {
	return func() {
		Reset(func() { CleanCache() })
		f(svr)
	}
}
