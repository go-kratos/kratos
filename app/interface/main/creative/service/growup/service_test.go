package growup

import (
	"context"
	"flag"
	"go-common/app/interface/main/creative/conf"
	"go-common/app/interface/main/creative/model/growup"
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

func Test_UpInfo(t *testing.T) {
	var (
		c         = context.TODO()
		err       error
		res       *growup.UpInfo
		MID       = int64(27515256)
		localHost = "127.0.0.1"
	)
	Convey("UpInfo", t, WithService(func(s *Service) {
		res, err = s.UpInfo(c, MID, localHost)
		So(err, ShouldNotBeNil)
		So(res, ShouldBeNil)
	}))
}
