package banner

import (
	"context"
	"flag"
	"path/filepath"
	"testing"
	"time"

	"go-common/app/interface/main/app-show/conf"
	"go-common/app/interface/main/app-show/model"

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
	dir, _ := filepath.Abs("../../cmd/app-show-test.toml")
	flag.Set("conf", dir)
	conf.Init()
	s = New(conf.Conf)
	time.Sleep(time.Second)
}

func TestDisplay(t *testing.T) {
	Convey("get Display data", t, WithService(func(s *Service) {
		res := s.Display(context.TODO(), model.PlatIPhone, 0, "", "", "", "iphone")
		So(res, ShouldNotBeEmpty)
	}))
}
