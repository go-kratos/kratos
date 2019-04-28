package rank

import (
	"context"
	"flag"
	"go-common/app/interface/main/app-show/conf"
	"path/filepath"
	"testing"
	"time"

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

func TestRankShow(t *testing.T) {
	Convey("get RankShow data", t, WithService(func(s *Service) {
		res := s.RankShow(context.TODO(), 0, 1, 1, 20, 0, "all")
		So(res, ShouldNotBeEmpty)
	}))
}
