package service

import (
	"context"
	"testing"

	"go-common/app/admin/main/reply/model"

	. "github.com/smartystreets/goconvey/convey"
)

func TestEvent(t *testing.T) {
	var (
		sub = &model.Subject{}
		rp  = &model.Reply{}
		rpt = &model.Report{}
		c   = context.Background()
	)
	Convey("test pub a event for reply", t, WithService(func(s *Service) {
		err := s.pubEvent(c, "test", 0, sub, rp, rpt)
		So(err, ShouldBeNil)
	}))
	Convey("test pub a event for search index", t, WithService(func(s *Service) {
		rps := map[int64]*model.Reply{}
		rps[rp.ID] = rp
		err := s.pubSearchReply(c, rps, 0)
		So(err, ShouldBeNil)
	}))
	Convey("test pub a event for search index", t, WithService(func(s *Service) {
		rpts := map[int64]*model.Report{}
		rpts[rpt.RpID] = rpt
		err := s.pubSearchReport(c, rpts, nil)
		So(err, ShouldBeNil)
	}))
}
