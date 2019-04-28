package goblin

import (
	"flag"
	"fmt"
	"path/filepath"
	"testing"
	"time"

	"go-common/app/interface/main/tv/conf"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	srv *Service
)

func init() {
	dir, _ := filepath.Abs("../../cmd/tv-interface.toml")
	flag.Set("conf", dir)
	conf.Init()
	srv = New(conf.Conf)
	time.Sleep(time.Second)
}

func WithService(f func(s *Service)) func() {
	return func() {
		Reset(func() {})
		f(srv)
	}
}

func TestService_Labels(t *testing.T) {
	Convey("TestService_Labels", t, WithService(func(s *Service) {
		results, err := s.Labels(ctx, 1, 2)
		So(err, ShouldBeNil)
		So(results, ShouldNotBeNil)
		for _, v := range results {
			fmt.Println(v.ParamName)
		}
	}))
}
