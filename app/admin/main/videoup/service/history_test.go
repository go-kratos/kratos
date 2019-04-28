package service

import (
	"context"
	"github.com/smartystreets/goconvey/convey"
	"testing"
)

func Test_trackArchiveInfo(t *testing.T) {
	convey.Convey("稿件编辑历史和track信息合并", t, WithService(func(s *Service) {
		c := context.TODO()
		aid := int64(1)
		inf, err := s.TrackArchiveInfo(c, aid)
		convey.So(err, convey.ShouldBeNil)
		convey.So(len(inf.Relation), convey.ShouldBeGreaterThanOrEqualTo, len(inf.EditHistory))
		for k, it := range inf.Relation {
			t.Logf("relation k(%d) it(%v)", k, it)
		}
		t.Logf("tr len(%d)", len(inf.Track))
	}))
}

func Test_editHistory(t *testing.T) {
	convey.Convey("hid获取稿件+分P编辑历史", t, WithService(func(s *Service) {
		c := context.TODO()
		hid := int64(1)
		_, err := s.EditHistory(c, hid)
		convey.So(err, convey.ShouldBeNil)
	}))
}

func Test_allEditHistory(t *testing.T) {
	convey.Convey("稿件编辑历史和track信息合并", t, WithService(func(s *Service) {
		c := context.TODO()
		aid := int64(10107879)
		h, err := s.AllEditHistory(c, aid)
		for _, o := range h {
			t.Logf("ah(%+v) vh(%+v)", o.ArcHistory, o.VHistory)
		}

		convey.So(err, convey.ShouldBeNil)
	}))
}
