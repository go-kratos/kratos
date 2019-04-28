package resource

import (
	"context"
	"flag"
	"path/filepath"
	"testing"
	"time"

	"go-common/app/interface/main/creative/conf"
	"go-common/app/interface/main/creative/service"
	"go-common/library/ecode"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	s *Service
)

func init() {
	dir, _ := filepath.Abs("../../cmd/creative.toml")
	flag.Set("conf", dir)
	conf.Init()
	rpcdaos := service.NewRPCDaos(conf.Conf)
	s = New(conf.Conf, rpcdaos)
	time.Sleep(time.Second)
}

func WithService(f func(s *Service)) func() {
	return func() {
		Reset(func() {})
		f(s)
	}
}

func Test_Resource(t *testing.T) {
	Convey("Resource", t, WithService(func(s *Service) {
		res, err := s.AcademyBanner(context.TODO(), "", "", "", "", "", "", 1, 2, int8(2), int64(123), false)
		So(err, ShouldEqual, ecode.NothingFound)
		So(res, ShouldNotBeEmpty)
	}))
}
