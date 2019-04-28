package service

import (
	"context"
	"flag"
	"fmt"
	"path/filepath"
	"testing"
	"time"

	"go-common/app/interface/main/credit/conf"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	s *Service
)

// func CleanCache() {
// 	c := context.TODO()
// 	pool := redis.NewPool(conf.Conf.Redis.Config)
// 	pool.Get(c).Do("FLUSHDB")
// }

func init() {
	dir, _ := filepath.Abs("../cmd/convey-test.toml")
	flag.Set("conf", dir)
	conf.Init()
	s = New(conf.Conf)
	time.Sleep(time.Second)
}

// func WithService(f func(s *Service)) func() {
// 	return func() {
// 		Reset(func() { CleanCache() })
// 		f(s)
// 	}
// }

func Test_LoadConf(t *testing.T) {
	Convey("should return err be nil", t, func() {
		s.loadConf()
		fmt.Printf("%+v", s.c.Judge)
	})
}

func Test_BatchBLKCases(t *testing.T) {
	ids := []int64{111, 22, 333}
	Convey("return someting", t, func() {
		cas, err := s.BatchBLKCases(context.TODO(), ids)
		So(err, ShouldBeNil)
		So(cas, ShouldNotBeNil)
	})
}
