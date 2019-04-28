package param

import (
	"flag"
	"path/filepath"
	"testing"
	"time"

	"go-common/app/interface/main/app-resource/conf"
	"go-common/app/interface/main/app-resource/model"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	s *Service
)

func WithService(f func(s *Service)) func() {
	return func() {
		f(s)
	}
}

func init() {
	dir, _ := filepath.Abs("../../cmd/app-resource-test.toml")
	flag.Set("conf", dir)
	conf.Init()
	s = New(conf.Conf)
	time.Sleep(time.Second)
}

func TestParam(t *testing.T) {
	Convey("get Param data", t, WithService(func(s *Service) {
		res, _, err := s.Param(model.PlatIPhone, 1, "")
		So(res, ShouldNotBeEmpty)
		So(err, ShouldBeNil)
	}))
}
