package service

import (
	"context"
	. "github.com/smartystreets/goconvey/convey"
	"testing"

	"go-common/app/admin/main/videoup/model/archive"
)

func TestService_Submit(t *testing.T) {
	var (
		c  = context.TODO()
		ap = &archive.ArcParam{}
	)
	Convey("Submit", t, WithService(func(s *Service) {
		err := svr.Submit(c, ap)
		So(err, ShouldBeNil)
	}))
}

func TestService_UpAccess(t *testing.T) {
	var (
		c  = context.TODO()
		ap = &archive.ArcParam{}
	)
	Convey("UpAccess", t, WithService(func(s *Service) {
		err := svr.UpAccess(c, ap)
		So(err, ShouldBeNil)
	}))
}

func TestService_UpArcDtime(t *testing.T) {
	var (
		c = context.TODO()
	)
	Convey("UpArcDtime", t, WithService(func(s *Service) {
		err := svr.UpArcDtime(c, 1, 12345)
		So(err, ShouldNotBeNil)
	}))
}

func TestService_UpAuther(t *testing.T) {
	var (
		c  = context.TODO()
		ap = &archive.ArcParam{}
	)
	Convey("UpAuther", t, WithService(func(s *Service) {
		err := svr.UpAuther(c, ap)
		So(err, ShouldBeNil)
	}))
}

func TestService_UpArchiveAttr(t *testing.T) {
	var (
		c = context.TODO()
	)
	attrs := make(map[uint]int32, 6)
	attrs[archive.AttrBitNoRank] = 0
	attrs[archive.AttrBitNoDynamic] = 0
	attrs[archive.AttrBitNoRecommend] = 0
	// forbid
	forbidAttrs := make(map[string]map[uint]int32, 3)
	forbidAttrs[archive.ForbidRank] = map[uint]int32{
		archive.ForbidRankMain:      0,
		archive.ForbidRankRecentArc: 0,
		archive.ForbidRankAllArc:    0,
	}
	forbidAttrs[archive.ForbidDynamic] = map[uint]int32{
		archive.ForbidDynamicMain: 0,
	}
	forbidAttrs[archive.ForbidRecommend] = map[uint]int32{
		archive.ForbidRecommendMain: 0,
	}
	Convey("UpArchiveAttr", t, WithService(func(s *Service) {
		err := svr.UpArchiveAttr(c, 1, 2, attrs, forbidAttrs, "")
		So(err, ShouldBeNil)
	}))
}

func TestService_Next(t *testing.T) {
	var (
		c = context.TODO()
	)
	Convey("Next", t, WithService(func(s *Service) {
		task, err := svr.Next(c, 6)
		So(task, ShouldNotBeNil)
		So(err, ShouldBeNil)
	}))
}

func TestService_UpArcTag(t *testing.T) {
	Convey("UpArcTag", t, WithService(func(s *Service) {
		//a.频道回查列表进入并提交的  b.tag改变
		c := context.TODO()
		//pm1(~a && b) -- archive_oper新增记录
		pm1 := &archive.TagParam{AID: 6, Tags: "haha1,haha2", FromChannelReview: ""}
		//pm2 (~a && ~b) -- 啥都没有
		pm2 := &archive.TagParam{AID: 6, Tags: "haha1,haha2", FromChannelReview: ""}
		//pm1(a && ~b) -- 新增flow_design
		pm3 := &archive.TagParam{AID: 6, Tags: "haha1,haha2", FromChannelReview: "1"}
		//pm2 (a && b) -- archive_oper新增
		pm4 := &archive.TagParam{AID: 6, Tags: "haha", FromChannelReview: "1"}
		err := svr.UpArcTag(c, 421, pm1)
		So(err, ShouldBeNil)

		err = svr.UpArcTag(c, 421, pm2)
		So(err, ShouldBeNil)

		err = svr.UpArcTag(c, 421, pm3)
		So(err, ShouldNotBeNil)

		err = svr.UpArcTag(c, 421, pm4)
		So(err, ShouldNotBeNil)
	}))
}

func TestService_GetChannelInfo(t *testing.T) {
	Convey("GetChannelInfo", t, WithService(func(s *Service) {
		info, err := s.GetChannelInfo(context.TODO(), []int64{10110255, 10110250})
		for aid, in := range info {
			channes := []*archive.Channel{}
			if in != nil {
				channes = in.Channels
			}

			t.Logf("aid=%d, in=%+v list the channels\r\n", aid, in)
			for _, ch := range channes {
				t.Logf("channel(%+v)\r\n", ch)
			}
		}

		So(err, ShouldBeNil)
	}))
}
