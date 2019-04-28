package service

import (
	"context"
	"flag"
	"os"
	"path/filepath"
	"testing"

	"go-common/app/interface/main/tag/conf"
	"go-common/library/cache/redis"

	. "github.com/smartystreets/goconvey/convey"
)

var testSvc *Service

func TestMain(m *testing.M) {
	flag.Parse()
	dir, _ := filepath.Abs("../cmd/tag-example.toml")
	flag.Set("conf", dir)
	conf.Init()
	testSvc = New(conf.Conf)
	os.Exit(m.Run())
}
func CleanCache() {
	c := context.Background()
	pool := redis.NewPool(conf.Conf.Redis.Tag.Redis)
	pool.Get(c).Do("FLUSHDB")
}

func WithService(f func(s *Service)) func() {
	return func() {
		Reset(func() { CleanCache() })
		f(testSvc)
	}
}
