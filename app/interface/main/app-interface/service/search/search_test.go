package search

import (
	"context"
	"flag"
	"path/filepath"
	"testing"
	"time"

	"go-common/app/interface/main/app-interface/conf"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	s *Service
)

func init() {
	dir, _ := filepath.Abs("../cmd/app-interface-test.toml")
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
		_, err := s.Search(context.Background(), 1, 2, "", "", "", "", "", "", "", "", "", "", "", "", 3, 4, 5, 6, 7, 8, 1, false, time.Now())
		So(err, ShouldBeNil)
	}))
}

func Test_SearchByType(t *testing.T) {
	Convey("get app banner", t, WithService(func(s *Service) {
		_, err := s.SearchByType(context.Background(), 1, 2, "", "", "", "", "", "", "", "", 3, 4, 5, 6, 7, 8, 8, 8, false, time.Now())
		So(err, ShouldBeNil)
	}))
}

func Test_SearchLive(t *testing.T) {
	Convey("get app SearchLive", t, WithService(func(s *Service) {
		_, err := s.SearchLive(context.Background(), 1, "", "", "", "", "", "", "", 3, 4, 5)
		So(err, ShouldBeNil)
	}))
}

func Test_upper(t *testing.T) {
	Convey("get app upper", t, WithService(func(s *Service) {
		_, err := s.upper(context.Background(), 1, "", "", "", "", "", "", "", 3, 4, 5, 3, 4, 5, 7, false, time.Now())
		So(err, ShouldBeNil)
	}))
}
