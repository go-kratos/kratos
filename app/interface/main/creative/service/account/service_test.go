package account

import (
	"context"
	"flag"
	"go-common/app/interface/main/creative/conf"
	accmdl "go-common/app/interface/main/creative/model/account"
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

func Test_Account(t *testing.T) {
	var (
		c         = context.TODO()
		err       error
		MID       = int64(2089809)
		localHost = "127.0.0.1"
		v         *accmdl.UpInfo
		m         *accmdl.MyInfo
		ip        = ""
	)
	Convey("MyInfo", t, WithService(func(s *Service) {
		m, err = s.MyInfo(c, MID, ip, time.Now())
		So(err, ShouldBeNil)
		So(m, ShouldNotBeNil)
	}))
	Convey("UpInfo", t, WithService(func(s *Service) {
		v, err = s.UpInfo(c, MID, localHost)
		So(err, ShouldBeNil)
		So(v, ShouldNotBeNil)
	}))
}
