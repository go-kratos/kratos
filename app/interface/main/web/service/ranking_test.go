package service

import (
	"context"
	"flag"
	"path/filepath"
	"testing"
	"time"

	"go-common/app/interface/main/web/conf"

	. "github.com/smartystreets/goconvey/convey"
)

var svf *Service

func init() {
	dir, _ := filepath.Abs("../cmd/web-interface-test.toml")
	flag.Set("conf", dir)
	if err := conf.Init(); err != nil {
		panic(err)
	}
	if svf == nil {
		svf = New(conf.Conf)
	}
	time.Sleep(time.Second)
}

func WithService(f func(s *Service)) func() {
	return func() {
		f(svf)
	}
}

func TestService_RankingRegion1(t *testing.T) {
	Convey("should return without err", t, WithService(func(svf *Service) {
		res, err := svf.RankingRegion(context.Background(), 129, 3, 0)
		So(err, ShouldBeNil)
		So(len(res), ShouldBeGreaterThan, 0)
	}))
}

func TestService_RankingRegion2(t *testing.T) {
	Convey("should return without err", t, WithService(func(svf *Service) {
		res, err := svf.RankingRegion(context.Background(), 129, 3, 1)
		So(err, ShouldBeNil)
		So(len(res), ShouldBeGreaterThan, 0)
	}))
}

func TestService_RankingTag(t *testing.T) {
	rid := int16(24)
	tagID := int64(358)
	Convey("should return without err", t, WithService(func(svf *Service) {
		res, err := svf.RankingTag(context.Background(), rid, tagID)
		So(err, ShouldBeNil)
		So(len(res), ShouldBeGreaterThan, 0)
	}))
}

func TestService_Ranking(t *testing.T) {
	var (
		rid      int16 = 23
		rankType       = 1
		day            = 1
		arcType        = 1
	)
	Convey("test ranking Ranking", t, WithService(func(s *Service) {
		res, err := s.Ranking(context.Background(), rid, rankType, day, arcType)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeNil)
	}))
}

func TestService_RankingIndex(t *testing.T) {
	Convey("test ranking RankingIndex", t, WithService(func(s *Service) {
		day := 1
		res, err := s.RankingIndex(context.Background(), day)
		So(err, ShouldBeNil)
		So(len(res), ShouldBeGreaterThan, 0)
	}))
}

func TestService_RankingRecommend(t *testing.T) {
	Convey("test ranking RankingRecommend", t, WithService(func(s *Service) {
		rid := 1
		res, err := s.RankingRecommend(context.Background(), int16(rid))
		So(err, ShouldBeNil)
		So(len(res), ShouldBeGreaterThan, 0)
	}))

}

func TestService_RegionCustom(t *testing.T) {
	Convey("test ranking ReginCustom", t, WithService(func(s *Service) {
		res, err := s.RegionCustom(context.Background())
		So(err, ShouldBeNil)
		So(res, ShouldNotBeNil)
	}))
}
