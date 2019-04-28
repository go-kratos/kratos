package dataflow

import (
	"context"
	"flag"
	"path/filepath"
	"testing"
	"time"

	"go-common/app/interface/main/app-interface/conf"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	s *Service
)

func init() {
	dir, _ := filepath.Abs("../cmd/app-interface-test.toml")
	flag.Set("conf", dir)
	conf.Init()
	s = New(conf.Conf)
	time.Sleep(time.Second)
}

func WithService(f func(s *Service)) func() {
	return func() {
		f(s)
	}
}

func TestReport(t *testing.T) {
	Convey("get app banner", t, WithService(func(s *Service) {
		err := s.Report(context.Background(), "1", "2", "3", "4", "5", time.Now())
		So(err, ShouldBeNil)
	}))
}
