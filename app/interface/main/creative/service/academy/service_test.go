package academy

import (
	"context"
	"flag"
	"path/filepath"
	"testing"
	"time"

	"go-common/app/interface/main/creative/conf"
	"go-common/app/interface/main/creative/model/academy"
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

func Test_TagList(t *testing.T) {
	var (
		c = context.TODO()
	)
	Convey("TagList", t, WithService(func(s *Service) {
		res, err := s.TagList(c)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeNil)
	}))
}

func Test_ArchivesWithES(t *testing.T) {
	var (
		c   = context.TODO()
		aca = &academy.EsParam{
			Tid:      []int64{},
			Business: 1,
			Pn:       1,
			Ps:       10,
			Keyword:  "",
			Order:    "",
			IP:       "127.0.0.1",
		}
	)
	Convey("Archives", t, WithService(func(s *Service) {
		res, err := s.ArchivesWithES(c, aca)
		//spew.Dump(res, err)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeNil)
	}))
}

func Test_AddFeedBack(t *testing.T) {
	var (
		c        = context.TODO()
		category = "视频"
		course   = "图像处理"
		suggest  = "画质太差"
		mid      = int64(123)
	)
	Convey("AddFeedBack", t, WithService(func(s *Service) {
		id, err := s.AddFeedBack(c, category, course, suggest, mid)
		So(err, ShouldBeNil)
		So(id, ShouldBeGreaterThan, 0)
	}))
}

func Test_RecommendV2(t *testing.T) {
	var (
		c   = context.TODO()
		mid = int64(123)
	)
	Convey("RecommendV2", t, WithService(func(s *Service) {
		_, err := s.RecommendV2(c, mid)
		So(err, ShouldBeNil)
	}))
}
