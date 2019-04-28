package service

import (
	"context"
	"flag"
	"path/filepath"
	"time"

	"go-common/app/service/main/location/conf"
	"go-common/library/cache/redis"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	s *Service
)

func CleanCache() {
	c := context.TODO()
	pool := redis.NewPool(conf.Conf.Redis.Zlimit.Config)
	pool.Get(c).Do("FLUSHDB")
}

func init() {
	dir, _ := filepath.Abs("../cmd/convey-test.toml")
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
