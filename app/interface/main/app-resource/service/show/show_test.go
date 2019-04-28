package show

import (
	"context"
	"flag"
	"path/filepath"
	"testing"
	"time"

	"go-common/app/interface/main/app-resource/conf"

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

func TestTabs(t *testing.T) {
	Convey("get Tabs data", t, WithService(func(s *Service) {
		res, _, _, err := s.Tabs(context.TODO(), 1, 1, "xxxx", "xxx", "iphone", "hans", 111)
		So(res, ShouldNotBeEmpty)
		So(err, ShouldBeNil)
	}))
}
