package service

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestServiceReport(t *testing.T) {
	var (
		tp  int32 = 3
		pn  int32 = 1
		ps  int32 = 10
		oid int64 = 22786099
		mid int64 = 35152246
		// rptMid int64 = 35152245
		uname        = "shunza"
		stime        = "2018-05-02 00:00:00"
		etime        = "2018-05-03 00:00:00"
		state  int32 = 1
		reason int32 = 2
		// rid          = []int64{1, 2, 3}
		ids          = []int64{1234, 2345, 34567}
		id     int64 = 123456
		uid    int64 = 35152246
		action int32 = 1
		audit  int32 = 4
	)
	Convey("ReportList", func() {
		testSvc.ReportList(context.TODO(), nil)
	})
	Convey("ReportInfo", func() {
		testSvc.ReportInfo(context.TODO(), id)
	})
	Convey("ReportHandle", func() {
		testSvc.ReportHandle(context.TODO(), uname, uid, id, audit, action)
	})
	Convey("ReportState", func() {
		testSvc.ReportState(context.TODO(), id, state)
	})
	Convey("ReportIgnore", func() {
		testSvc.ReportIgnore(context.TODO(), uname, audit, ids)
	})
	Convey("ReportDelete", func() {
		testSvc.ReportDelete(context.TODO(), uname, uid, ids, audit)
	})
	Convey("ReportPunish", func() {
		testSvc.ReportPunish(context.TODO(), uname, "", "", ids, reason, -10, 1, 2, 0)
	})
	Convey("ReportLogList", func() {
		testSvc.ReportLogList(context.TODO(), oid, 233, mid, 123, tp, pn, ps, []int64{}, stime, etime, uname)
	})
	Convey("ReportLogInfo", func() {
		testSvc.ReportLogInfo(context.TODO(), id)
	})
}
