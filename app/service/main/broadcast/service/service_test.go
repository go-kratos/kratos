package service

import (
	"flag"
	"path/filepath"
	"time"

	"go-common/library/conf/paladin"
	"go-common/library/naming/discovery"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	s *Service
)

func init() {
	dir, _ := filepath.Abs("../cmd/local")
	flag.Set("conf", dir)
	if err := paladin.Init(); err != nil {
		panic(err)
	}
	s = New(discovery.New(nil))
	time.Sleep(time.Second)
}

func CleanCache() {
	//c := context.TODO()
	//pool := redis.NewPool(conf.Conf.Redis.Config)
	//pool.Get(c).Do("FLUSHDB")
}

func WithService(f func(s *Service)) func() {
	return func() {
		Reset(func() { CleanCache() })
		f(s)
	}
}
