package service

import (
	"flag"
	"path/filepath"
	"testing"
	"time"

	"go-common/app/job/main/creative/conf"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	s *Service
)

func init() {
	dir, _ := filepath.Abs("../cmd/creative-job.toml")
	flag.Set("conf", dir)
	conf.Init()
	s = New(conf.Conf)
	time.Sleep(time.Second)
}

func WithService(f func(s *Service)) func() {
	return func() {
		Reset(func() {})
		f(s)
	}
}

func Test_Pub(t *testing.T) {
	Convey("pub", t, WithService(func(s *Service) {
		Convey("pub", func() {
			err := s.pub(int64(2089809), 0, 1)
			So(err, ShouldBeNil)
		})
	}))
}
