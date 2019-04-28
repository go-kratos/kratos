package service

import (
	"context"
	"flag"
	"path/filepath"
	"time"

	"go-common/app/interface/main/creative/conf"
	"go-common/library/cache/redis"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	pvr *Public
)

func CleanCache() {
	c := context.TODO()
	pool := redis.NewPool(conf.Conf.Redis.Antispam.Config)
	pool.Get(c).Do("FLUSHDB")
}

func init() {
	dir, _ := filepath.Abs("../cmd/creative.toml")
	flag.Set("conf", dir)
	conf.Init()
	pvr = New(conf.Conf, NewRPCDaos(conf.Conf))
	time.Sleep(time.Second)
}
func WithService(f func(p *Public)) func() {
	return func() {
		Reset(func() { CleanCache() })
		f(pvr)
	}
}
