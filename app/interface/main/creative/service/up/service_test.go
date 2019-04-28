package up

import (
	"context"
	"flag"
	"go-common/app/interface/main/creative/conf"
	"path/filepath"
	"testing"
	"time"

	"go-common/app/interface/main/creative/service"

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

func Test_UpSwitch(t *testing.T) {
	var (
		mid int64 = 27515244
		c         = context.Background()
	)
	Convey("UpSwitch", t, WithService(func(s *Service) {
		res, err := s.UpSwitch(c, mid, "")
		So(err, ShouldBeNil)
		So(res, ShouldNotBeNil)
	}))
}

func TestService_ShowStaff(t *testing.T) {
	var c = context.Background()
	Convey("ShowStaff", t, func(ctx C) {
		_, err := s.ShowStaff(c, 27515244)
		So(err, ShouldBeNil)
	})
}
