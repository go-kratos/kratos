package show

import (
	"context"
	"flag"
	"go-common/app/interface/main/app-show/conf"
	"go-common/app/interface/main/app-show/model"
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

func TestIndex(t *testing.T) {
	Convey("get Index data", t, WithService(func(s *Service) {
		res := s.Index(context.TODO(), 0, model.PlatIPhone, 0, "", "", "", "", "", "iphone", "phone", _initlanguage, "", false, time.Now())
		So(res, ShouldNotBeEmpty)
	}))
}

func TestChange(t *testing.T) {
	Convey("get Change data", t, WithService(func(s *Service) {
		res := s.Change(context.TODO(), 1, 1, 1, 1, "", "", "", "", "")
		So(res, ShouldNotBeEmpty)
	}))
}

func TestRegionChange(t *testing.T) {
	Convey("get RegionChange data", t, WithService(func(s *Service) {
		res := s.RegionChange(context.TODO(), 1, 1, 1, 1, "")
		So(res, ShouldNotBeEmpty)
	}))
}

func TestBangumiChange(t *testing.T) {
	Convey("get BangumiChange data", t, WithService(func(s *Service) {
		res := s.BangumiChange(context.TODO(), 1, 1)
		So(res, ShouldNotBeEmpty)
	}))
}

func TestFeedIndex(t *testing.T) {
	Convey("get FeedIndex data", t, WithService(func(s *Service) {
		res := s.FeedIndex(context.TODO(), 1, 1, 1, 1, 1, "", "", "", "", time.Now())
		So(res, ShouldNotBeEmpty)
	}))
}
