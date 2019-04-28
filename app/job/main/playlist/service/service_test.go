package service

import (
	"context"
	"flag"
	"path/filepath"
	"testing"
	"time"

	"go-common/app/job/main/playlist/conf"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	svr *Service
)

func init() {
	dir, _ := filepath.Abs("../cmd/playlist-job-test.toml")
	flag.Set("conf", dir)
	conf.Init()
	svr = New(conf.Conf)
	time.Sleep(time.Second)

}

func WithService(f func(s *Service)) func() {
	return func() {
		f(svr)
	}
}

func TestService_Ping(t *testing.T) {
	Convey("test ping", t, WithService(func(s *Service) {
		err := s.Ping(context.TODO())
		So(err, ShouldBeNil)
	}))
}
