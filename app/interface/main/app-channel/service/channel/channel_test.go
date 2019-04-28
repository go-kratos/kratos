package channel

import (
	"context"
	"flag"
	"path/filepath"
	"testing"
	"time"

	"go-common/app/interface/main/app-channel/conf"
	"go-common/app/interface/main/app-channel/model"

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
	dir, _ := filepath.Abs("../../cmd/app-channel-test.toml")
	flag.Set("conf", dir)
	conf.Init()
	s = New(conf.Conf)
	time.Sleep(time.Second)
}

func TestTab(t *testing.T) {
	Convey("get Tab data", t, WithService(func(s *Service) {
		res, err := s.Tab(context.TODO(), 1, 1, "", 1)
		So(res, ShouldNotBeEmpty)
		So(err, ShouldBeNil)
	}))
}

func TestSubscribeAdd(t *testing.T) {
	Convey("get SubscribeAdd data", t, WithService(func(s *Service) {
		err := s.SubscribeAdd(context.TODO(), 1, 1, time.Now())
		So(err, ShouldBeNil)
	}))
}

func TestSubscribeCancel(t *testing.T) {
	Convey("get SubscribeCancel data", t, WithService(func(s *Service) {
		err := s.SubscribeCancel(context.TODO(), 1, 1, time.Now())
		So(err, ShouldBeNil)
	}))
}

func TestSubscribeUpdate(t *testing.T) {
	Convey("get SubscribeUpdate data", t, WithService(func(s *Service) {
		err := s.SubscribeUpdate(context.TODO(), 1, "")
		So(err, ShouldBeNil)
	}))
}

func TestList(t *testing.T) {
	Convey("get List data", t, WithService(func(s *Service) {
		res, err := s.List(context.TODO(), 1, model.PlatIPhone, 1, 1, "", "iphone", "phone", "hans")
		So(err, ShouldBeNil)
		So(res, ShouldNotBeEmpty)
	}))
}

func TestRecommend(t *testing.T) {
	Convey("get Recommend data", t, WithService(func(s *Service) {
		res, err := s.Recommend(context.TODO(), 1, 1)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeEmpty)
	}))
}

func TestSubscribe(t *testing.T) {
	Convey("get Subscribe data", t, WithService(func(s *Service) {
		res, err := s.Subscribe(context.TODO(), 1, 1)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeEmpty)
	}))
}

func TestDiscover(t *testing.T) {
	Convey("get Discover data", t, WithService(func(s *Service) {
		res, err := s.Discover(context.TODO(), 1, 1, 1)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeEmpty)
	}))
}

func TestCategory(t *testing.T) {
	Convey("get Category data", t, WithService(func(s *Service) {
		res, err := s.Category(context.TODO(), 1)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeEmpty)
	}))
}

func TestRegionList(t *testing.T) {
	Convey("get RegionList data", t, WithService(func(s *Service) {
		_, _, _, err := s.RegionList(context.TODO(), model.PlatIPhone, 1, "iphone", "device", "hans")
		So(err, ShouldBeNil)
	}))
}

func TestIndex(t *testing.T) {
	Convey("get Index data", t, WithService(func(s *Service) {
		res, err := s.Index(context.TODO(), 1, 1, 1, model.PlatIPhone, "", "", "", "", 1, 1, 1, 0, 0, 0, false, time.Now())
		So(err, ShouldBeNil)
		So(res, ShouldNotBeEmpty)
	}))
}
