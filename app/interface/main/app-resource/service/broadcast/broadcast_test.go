package broadcast

import (
	"context"
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

func TestServerList(t *testing.T) {
	Convey("get ServerList data", t, WithService(func(s *Service) {
		res, err := s.ServerList(context.Background(), model.PlatAndroid)
		So(res, ShouldNotBeEmpty)
		So(err, ShouldBeNil)
	}))
}
