package service

import (
	"flag"
	"path/filepath"
	"testing"
	"time"

	"go-common/app/job/main/growup/conf"
	"go-common/app/job/main/growup/service/ctrl"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_Temp(t *testing.T) {
	Convey("temp test", t, func() {})
}

var (
	srv *Service
)

func init() {
	dir, _ := filepath.Abs("../cmd/growup-job.toml")
	flag.Set("conf", dir)
	conf.Init()
	srv = New(conf.Conf, ctrl.NewUnboundedExecutor())
	time.Sleep(time.Second)
}

func WithService(f func(s *Service)) func() {
	return func() {
		// Reset(func() { CleanCache() })
		f(srv)
	}
}
