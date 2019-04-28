package medal

import (
	"context"
	"flag"
	"go-common/app/interface/main/creative/conf"
	mdMdl "go-common/app/interface/main/creative/model/medal"
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

func Test_Medal(t *testing.T) {
	var (
		mid  int64 = 1
		name       = "name"
		c          = context.TODO()
		err  error
		res  *mdMdl.Medal
	)
	Convey("Medal", t, WithService(func(s *Service) {
		res, err = s.Medal(c, mid)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeNil)
	}))
	Convey("OpenMedal", t, WithService(func(s *Service) {
		err = s.OpenMedal(c, mid, name)
		So(err, ShouldBeNil)
	}))
}
