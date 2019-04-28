package app

import (
	"context"
	"flag"
	"go-common/app/interface/main/creative/conf"
	"go-common/app/interface/main/creative/model/app"
	"path/filepath"
	"testing"
	"time"

	"go-common/app/interface/main/creative/service"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	s *Service
	p *service.Public
)

func init() {
	dir, _ := filepath.Abs("../../cmd/creative.toml")
	flag.Set("conf", dir)
	conf.Init()
	rpcdaos := service.NewRPCDaos(conf.Conf)
	p = service.New(conf.Conf, rpcdaos)
	s = New(conf.Conf, rpcdaos, p)
	time.Sleep(time.Second)
}

func WithService(f func(s *Service)) func() {
	return func() {
		Reset(func() {})
		f(s)
	}
}

func Test_PortalConfig(t *testing.T) {
	var (
		c         = context.TODO()
		err       error
		res       []*app.Portal
		artAuthor = 1
		build     = 524000
		ty        = 0
		plat      = "android"
		mid       = int64(2089809)
	)
	Convey("Portals", t, WithService(func(s *Service) {
		res, err = s.Portals(c, mid, artAuthor, build, ty, plat, 2)
		So(err, ShouldBeNil)
		So(len(res), ShouldBeGreaterThan, 0)
	}))
}
