package tag

import (
	"context"
	"flag"
	"path/filepath"
	"testing"
	"time"

	"go-common/app/interface/main/app-tag/conf"
	"go-common/app/interface/main/app-tag/model"

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
	dir, _ := filepath.Abs("../../cmd/app-tag-test.toml")
	flag.Set("conf", dir)
	conf.Init()
	s = New(conf.Conf)
	time.Sleep(time.Second)
}

func TestTagDynamic(t *testing.T) {
	Convey("get TagDynamic data", t, WithService(func(s *Service) {
		res := s.TagDynamic(context.TODO(), model.PlatIPhone, 0, 0, 0, 1217733, 0, "iphone", "", "", time.Now(), false)
		So(res, ShouldNotBeEmpty)
	}))
}

func TestTagDetail(t *testing.T) {
	Convey("get TagDetail data", t, WithService(func(s *Service) {
		res, _, err := s.TagDetail(context.TODO(), model.PlatIPhone, 0, 1217733, time.Now())
		So(res, ShouldNotBeEmpty)
		So(err, ShouldBeNil)
	}))
}
