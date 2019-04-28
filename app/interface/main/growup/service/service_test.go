package service

import (
	"context"
	"flag"
	"path/filepath"
	"testing"
	"time"

	"go-common/app/interface/main/growup/conf"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	srv *Service
)

func init() {
	dir, _ := filepath.Abs("../cmd/growup-interface.toml")
	flag.Set("conf", dir)
	conf.Init()
	srv = New(conf.Conf)
	time.Sleep(time.Second)
}

func WithService(f func(s *Service)) func() {
	return func() {
		// Reset(func() { CleanCache() })
		f(srv)
	}
}

func Test_GetUpStatus(t *testing.T) {
	var (
		mid = int64(1011)
	)
	Convey("interface", t, WithService(func(s *Service) {
		_, err := s.GetUpStatus(context.Background(), mid, "127.0.0.1")
		So(err, ShouldBeNil)
	}))
}

func Test_JoinAv(t *testing.T) {
	var (
		accountType = 1
		mid         = int64(1011)
		signType    = 2
	)
	Convey("interface", t, WithService(func(s *Service) {
		err := s.JoinAv(context.Background(), accountType, mid, signType)
		So(err, ShouldBeNil)
	}))
}

func Test_Quit(t *testing.T) {
	var (
		mid    = int64(1011)
		reason = "quit"
	)
	Convey("interface", t, WithService(func(s *Service) {
		err := s.Quit(context.Background(), mid, reason)
		So(err, ShouldBeNil)
	}))
}
