package job

import (
	"context"
	"flag"
	"path/filepath"
	"testing"
	"time"

	"go-common/app/interface/main/web-show/conf"

	. "github.com/smartystreets/goconvey/convey"
)

var svf *Service

func init() {
	dir, _ := filepath.Abs("../../cmd/web-show-test.toml")
	flag.Set("conf", dir)
	if err := conf.Init(); err != nil {
		panic(err)
	}
	if svf == nil {
		svf = New(conf.Conf)
	}
	time.Sleep(time.Second)
}

func WithService(f func(s *Service)) func() {
	return func() {
		f(svf)
	}
}
func TestService_Jobs(t *testing.T) {
	Convey("should return without err", t, WithService(func(svf *Service) {
		res := svf.Jobs(context.TODO())
		So(res, ShouldNotBeNil)
		Printf("%+v", res)
	}))
}
