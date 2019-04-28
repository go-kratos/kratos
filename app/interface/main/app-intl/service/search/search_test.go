package search

import (
	"context"
	"flag"
	"path/filepath"
	"testing"
	"time"

	"go-common/app/interface/main/app-intl/conf"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	s *Service
)

func init() {
	dir, _ := filepath.Abs("../cmd/app-intl-test.toml")
	flag.Set("conf", dir)
	conf.Init()
	s = New(conf.Conf)
	time.Sleep(time.Second)
}

func WithService(f func(s *Service)) func() {
	return func() {
		f(s)
	}
}

func Test_Search(t *testing.T) {
	Convey("get app banner", t, WithService(func(s *Service) {
		_, err := s.Search(context.Background(), 1, 2, "", "", "", "", "", "", "", "", "", "", "", 3, 4, 5, 6, 7, 8, time.Now())
		So(err, ShouldBeNil)
	}))
}

func Test_SearchByType(t *testing.T) {
	Convey("get app banner", t, WithService(func(s *Service) {
		_, err := s.SearchByType(context.Background(), 1, 2, "", "", "", "", "", "", "", "", 3, 4, 5, 6, 7, 8, 8, 8, time.Now())
		So(err, ShouldBeNil)
	}))
}
