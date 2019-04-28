package service

import (
	"context"
	"flag"
	"path/filepath"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
	"go-common/app/service/main/videoup/conf"
	"go-common/library/queue/databus/report"
)

var (
	svr *Service
)

func init() {
	dir, _ := filepath.Abs("../cmd/videoup-service.toml")
	flag.Set("conf", dir)
	conf.Init()
	svr = New(conf.Conf)
	time.Sleep(time.Second)
	report.InitUser(nil)
}

func WithService(f func(s *Service)) func() {
	return func() {
		f(svr)
	}
}

func TestService_Types(t *testing.T) {
	var (
		c = context.TODO()
	)
	Convey("Types", t, WithService(func(s *Service) {
		data := svr.Types(c)
		So(data, ShouldNotBeNil)
	}))
}
