package service

import (
	"context"
	"time"

	"go-common/app/admin/main/reply/model"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestReport(t *testing.T) {
	c := context.Background()
	Convey("report del and recover", t, WithService(func(s *Service) {
		sub, err := s.subject(c, 5464686, 1)
		So(err, ShouldBeNil)
		err = s.reportDel(c, []int64{5464686}, []int64{894726080}, 23, 0, 1, 2, 0, 0, 0, false, "testadmin", "test", "content")
		So(err, ShouldBeNil)
		rpt, err := s.report(c, 5464686, 894726080)
		So(err, ShouldBeNil)
		So(rpt.State, ShouldEqual, model.ReportStateDelete2)
		rp, err := s.reply(c, 5464686, 894726080)
		So(err, ShouldBeNil)
		So(rp.State, ShouldEqual, model.StateDelAdmin)
		sub2, err := s.subject(c, 5464686, 1)
		So(err, ShouldBeNil)
		So(sub2.RCount-sub.RCount, ShouldEqual, -1)
		So(sub2.ACount-sub.ACount, ShouldEqual, -1)
		alog, err := s.dao.AdminLog(c, 894726080)
		So(err, ShouldBeNil)
		So(alog.State, ShouldEqual, model.AdminOperRptDel2)
		err = s.reportRecover(c, []int64{5464686}, []int64{894726080}, 23, 1, 2, "test")
		So(err, ShouldBeNil)
		rpt, err = s.report(c, 5464686, 894726080)
		So(err, ShouldBeNil)
		So(rpt.State, ShouldEqual, model.ReportStateIgnore2)
		rp, err = s.reply(c, 5464686, 894726080)
		So(err, ShouldBeNil)
		So(rp.State, ShouldEqual, model.StateNormal)
		sub2, err = s.subject(c, 5464686, 1)
		So(err, ShouldBeNil)
		So(sub2.RCount-sub.RCount, ShouldEqual, 0)
		So(sub2.ACount-sub.ACount, ShouldEqual, 0)
		alog, err = s.dao.AdminLog(c, 894726080)
		So(err, ShouldBeNil)
		So(alog.State, ShouldEqual, model.AdminOperRptRecover2)
	}))

	Convey("reply transfer", t, WithService(func(s *Service) {
		err := s.ReportTransfer(c, []int64{5464686}, []int64{894726080}, 23, "asd", 1, 2, "test")
		So(err, ShouldBeNil)
		rpt, err := s.report(c, 5464686, 894726080)
		So(err, ShouldBeNil)
		So(rpt.State, ShouldEqual, model.ReportStateNew2)
		So(rpt.AttrVal(model.ReportAttrTransferred), ShouldEqual, model.AttrYes)
		s.dao.UpReportsAttrBit(c, []int64{5464686}, []int64{894726080}, model.ReportAttrTransferred, model.AttrNo, time.Now())
	}))

	Convey("reply ignore", t, WithService(func(s *Service) {
		err := s.ReportIgnore(c, []int64{5464686}, []int64{894726080}, 23, "asd", 1, 2, "test", false)
		rpt, err := s.report(c, 5464686, 894726080)
		So(err, ShouldBeNil)
		So(rpt.State, ShouldEqual, model.ReportStateIgnore2)
	}))
}
