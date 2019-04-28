package template

import (
	"context"
	"flag"
	"go-common/app/interface/main/creative/conf"
	"go-common/app/interface/main/creative/model/template"
	"path/filepath"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
	"go-common/app/interface/main/creative/service"
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

func Test_Templates(t *testing.T) {
	var (
		mid int64 = 1
		c         = context.TODO()
		err error
		res []*template.Template
	)
	Convey("Templates", t, WithService(func(s *Service) {
		res, err = s.Templates(c, mid)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeNil)
	}))
}
