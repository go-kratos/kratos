package geetest

import (
	"context"
	"flag"
	"go-common/app/interface/main/creative/conf"
	m "go-common/app/interface/main/creative/model/geetest"
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

func Test_PreProcess(t *testing.T) {
	var (
		c          = context.TODO()
		err        error
		res        *m.ProcessRes
		mid        = int64(2089809)
		ip         = "127.0.0.1"
		clientType = "clientType"
		newCaptcha = 1
	)
	Convey("PreProcess", t, WithService(func(s *Service) {
		res, err = s.PreProcess(c, mid, ip, clientType, newCaptcha)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeNil)
	}))
}
