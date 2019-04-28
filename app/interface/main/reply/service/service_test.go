package service

import (
	"context"
	"flag"
	"path/filepath"
	"testing"
	"time"

	"go-common/app/interface/main/reply/conf"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	testSvr *Service
)

func init() {
	dir, _ := filepath.Abs("../cmd/reply-test.toml")
	flag.Set("conf", dir)
	conf.Init()
	testSvr = New(conf.Conf)
	time.Sleep(time.Second)
}

func CleanCache() {
	//c := context.Background()
	//pool := redis.NewPool(conf.Conf.Redis.Config)
	//pool.Get(c).Do("FLUSHDB")
}

func WithService(f func(s *Service)) func() {
	return func() {
		Reset(func() { CleanCache() })
		f(testSvr)
	}
}

func TestFetchFans(t *testing.T) {
	s := &Service{}
	s.FetchFans(context.Background(), []int64{1, 2}, 11)
}
