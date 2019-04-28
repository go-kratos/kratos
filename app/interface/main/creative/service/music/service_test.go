package music

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

func Test_loadPreValues(t *testing.T) {
	Convey("loadPreValues", t, WithService(func(s *Service) {
		s.loadPreValues()
		vals := s.BgmList(context.Background(), 1)
		So(vals, ShouldNotBeNil)
	}))
}
