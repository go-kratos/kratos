package data

import (
	"context"
	"flag"
	"go-common/app/interface/main/creative/conf"
	"path/filepath"
	"testing"
	"time"

	"go-common/app/interface/main/creative/service"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	s *Service
	p *service.Public
)

func init() {
	dir, _ := filepath.Abs("../../cmd/creative.toml")
	flag.Set("conf", dir)
	conf.Init()
	rpcdaos := service.NewRPCDaos(conf.Conf)
	p = service.New(conf.Conf, rpcdaos)
	s = New(conf.Conf, rpcdaos, p)
	time.Sleep(time.Second)
}

func WithService(f func(s *Service)) func() {
	return func() {
		Reset(func() {})
		f(s)
	}
}

var (
	c   = context.TODO()
	MID = int64(27515256)
	ty  = int8(1)
	ip  = "127.0.0.1"
	dt  = "20180301"
)

func Test_AppStat(t *testing.T) {
	Convey("AppStat", t, WithService(func(s *Service) {
		Convey("AppStat", func() {
			res, err := s.AppStat(c, MID)
			So(err, ShouldBeNil)
			So(res, ShouldNotBeNil)
		})
	}))
}

func Test_ViewerBase(t *testing.T) {
	Convey("ViewerBase", t, WithService(func(s *Service) {
		Convey("ViewerBase", func() {
			res, err := s.ViewerBase(c, MID)
			So(err, ShouldBeNil)
			So(res, ShouldNotBeNil)
		})
	}))
}

func Test_ViewerArea(t *testing.T) {
	Convey("ViewerArea", t, WithService(func(s *Service) {
		Convey("ViewerArea", func() {
			res, err := s.ViewerArea(c, MID)
			So(err, ShouldBeNil)
			So(res, ShouldNotBeNil)
		})
	}))
}

func Test_CacheTrend(t *testing.T) {
	Convey("CacheTrend", t, WithService(func(s *Service) {
		Convey("CacheTrend", func() {
			res, err := s.CacheTrend(c, MID)
			So(err, ShouldBeNil)
			So(res, ShouldNotBeNil)
		})
	}))
}

func Test_RelationFansHistory(t *testing.T) {
	Convey("RelationFansHistory", t, WithService(func(s *Service) {
		Convey("RelationFansHistory", func() {
			res, err := s.RelationFansHistory(c, MID, dt)
			So(err, ShouldBeNil)
			So(res, ShouldNotBeNil)
		})
	}))
}

func Test_RelationFansMonth(t *testing.T) {
	Convey("RelationFansMonth", t, WithService(func(s *Service) {
		Convey("RelationFansMonth", func() {
			res, err := s.RelationFansMonth(c, MID)
			So(err, ShouldBeNil)
			So(res, ShouldNotBeNil)
		})
	}))
}

func Test_ViewerActionHour(t *testing.T) {
	Convey("ViewerActionHour", t, WithService(func(s *Service) {
		Convey("ViewerActionHour", func() {
			res, err := s.ViewerActionHour(c, MID)
			So(err, ShouldBeNil)
			So(res, ShouldNotBeNil)
		})
	}))
}

func Test_UpIncr(t *testing.T) {
	Convey("UpIncr", t, WithService(func(s *Service) {
		Convey("UpIncr", func() {
			res, err := s.UpIncr(c, MID, ty, ip)
			So(err, ShouldBeNil)
			So(res, ShouldNotBeNil)
		})
	}))
}

func Test_AppUpIncr(t *testing.T) {
	Convey("AppUpIncr", t, WithService(func(s *Service) {
		Convey("AppUpIncr", func() {
			res, err := s.AppUpIncr(c, MID, ty, ip)
			So(err, ShouldBeNil)
			So(res, ShouldNotBeNil)
		})
	}))
}

func Test_ThirtyDayArchive(t *testing.T) {
	Convey("ThirtyDayArchive", t, WithService(func(s *Service) {
		Convey("ThirtyDayArchive", func() {
			res, err := s.ThirtyDayArchive(c, MID, ty)
			So(err, ShouldBeNil)
			So(res, ShouldNotBeNil)
		})
	}))
}

func Test_ThirtyDayArticle(t *testing.T) {
	Convey("ThirtyDayArticle", t, WithService(func(s *Service) {
		Convey("ThirtyDayArticle", func() {
			res, err := s.ThirtyDayArticle(c, MID, ip)
			So(err, ShouldBeNil)
			So(res, ShouldNotBeNil)
		})
	}))
}
