package display

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

func Test_Zone(t *testing.T) {
	Convey("get Zone", t, WithService(func(s *Service) {
		zone := s.Zone(context.Background(), time.Now())
		So(zone, ShouldNotBeNil)
	}))
}
