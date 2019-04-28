package dao

import (
	"context"
	"go-common/app/admin/main/tag/model"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoReport(t *testing.T) {
	var (
		c  = context.TODO()
		id = int64(0)
	)
	convey.Convey("Report", t, func(ctx convey.C) {
		res, err := d.Report(c, id)
		ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(res, convey.ShouldBeNil)
		})
	})
}

func TestDaoReports(t *testing.T) {
	var (
		c   = context.TODO()
		ids = []int64{1, 2, 3}
	)
	convey.Convey("Reports", t, func(ctx convey.C) {
		rpt, rptMap, rptIDs, err := d.Reports(c, ids)
		ctx.Convey("Then err should be nil.rpt,rptMap,rptIDs should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(rptIDs, convey.ShouldHaveLength, 0)
			ctx.So(rptMap, convey.ShouldHaveLength, 0)
			ctx.So(rpt, convey.ShouldHaveLength, 0)
		})
	})
}

func TestDaoReportUser(t *testing.T) {
	var (
		c  = context.TODO()
		id = int64(0)
	)
	convey.Convey("ReportUser", t, func(ctx convey.C) {
		res, err := d.ReportUser(c, id)
		ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(res, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoReportUsers(t *testing.T) {
	var (
		c   = context.TODO()
		ids = []int64{1, 2, 3}
	)
	convey.Convey("ReportUsers", t, func(ctx convey.C) {
		users, userMap, err := d.ReportUsers(c, ids)
		ctx.Convey("Then err should be nil.users,userMap should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(userMap, convey.ShouldHaveLength, 0)
			ctx.So(users, convey.ShouldHaveLength, 0)
		})
	})
}

func TestDaoUpReportState(t *testing.T) {
	var (
		c     = context.TODO()
		id    = int64(0)
		state = int32(0)
	)
	convey.Convey("UpReportState", t, func(ctx convey.C) {
		affect, err := d.UpReportState(c, id, state)
		ctx.Convey("Then err should be nil.affect should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(affect, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoUpReportsState(t *testing.T) {
	var (
		c     = context.TODO()
		ids   = []int64{0}
		state = int32(0)
	)
	convey.Convey("UpReportsState", t, func(ctx convey.C) {
		affect, err := d.UpReportsState(c, ids, state)
		ctx.Convey("Then err should be nil.affect should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(affect, convey.ShouldBeGreaterThanOrEqualTo, 0)
		})
	})
}

func TestDaoAddReportLog(t *testing.T) {
	var (
		c = context.TODO()
		l = &model.ReportLog{}
	)
	convey.Convey("AddReportLog", t, func(ctx convey.C) {
		id, err := d.AddReportLog(c, l)
		ctx.Convey("Then err should be nil.id should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(id, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoAddReportLogs(t *testing.T) {
	var (
		c    = context.TODO()
		sqls = []string{}
	)
	convey.Convey("AddReportLogs", t, func(ctx convey.C) {
		id, err := d.AddReportLogs(c, sqls)
		ctx.Convey("Then err should be nil.id should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldNotBeNil)
			ctx.So(id, convey.ShouldBeGreaterThanOrEqualTo, 0)
		})
	})
}

func TestDaoUpdateMorals(t *testing.T) {
	var (
		c     = context.TODO()
		moral = int32(0)
		ids   = []int64{0}
	)
	convey.Convey("UpdateMorals", t, func(ctx convey.C) {
		affect, err := d.UpdateMorals(c, moral, ids)
		ctx.Convey("Then err should be nil.affect should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(affect, convey.ShouldBeGreaterThanOrEqualTo, 0)
		})
	})
}

func TestDaoReportByOidMid(t *testing.T) {
	var (
		c   = context.TODO()
		mid = int64(0)
		oid = int64(0)
		tp  = int32(0)
	)
	convey.Convey("ReportByOidMid", t, func(ctx convey.C) {
		rpt, rptIDs, tids, err := d.ReportByOidMid(c, mid, oid, tp)
		ctx.Convey("Then err should be nil.rpt,rptIDs,tids should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(tids, convey.ShouldHaveLength, 0)
			ctx.So(rptIDs, convey.ShouldHaveLength, 0)
			ctx.So(rpt, convey.ShouldHaveLength, 0)
		})
	})
}

func TestDaoReportLog(t *testing.T) {
	var (
		c  = context.TODO()
		id = int64(0)
	)
	convey.Convey("ReportLog", t, func(ctx convey.C) {
		res, tids, err := d.ReportLog(c, id)
		ctx.Convey("Then err should be nil.res,tids should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(tids, convey.ShouldNotBeNil)
			ctx.So(res, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoReportLogCount(t *testing.T) {
	var (
		c   = context.TODO()
		sql = ""
	)
	convey.Convey("ReportLogCount", t, func(ctx convey.C) {
		count, err := d.ReportLogCount(c, sql)
		ctx.Convey("Then err should be nil.count should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(count, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoReportLogList(t *testing.T) {
	var (
		c     = context.TODO()
		sql   = ""
		start = int32(0)
		end   = int32(0)
	)
	convey.Convey("ReportLogList", t, func(ctx convey.C) {
		res, tids, err := d.ReportLogList(c, sql, start, end)
		ctx.Convey("Then err should be nil.res,tids should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(tids, convey.ShouldHaveLength, 0)
			ctx.So(res, convey.ShouldHaveLength, 0)
		})
	})
}

func TestDaoReportLogByRptID(t *testing.T) {
	var (
		c      = context.TODO()
		rptIDs = []int64{1, 2, 3}
	)
	convey.Convey("ReportLogByRptID", t, func(ctx convey.C) {
		logs, logMap, err := d.ReportLogByRptID(c, rptIDs)
		ctx.Convey("Then err should be nil.logs,logMap should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(logMap, convey.ShouldHaveLength, 0)
			ctx.So(logs, convey.ShouldHaveLength, 0)
		})
	})
}

func TestDaoReportCount(t *testing.T) {
	var (
		c      = context.TODO()
		sqlStr = []string{}
		order  = "rpt.id"
		state  = int32(0)
	)
	convey.Convey("ReportCount", t, func(ctx convey.C) {
		count, err := d.ReportCount(c, sqlStr, order, state)
		ctx.Convey("Then err should be nil.count should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(count, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoReportInfoList(t *testing.T) {
	var (
		c      = context.TODO()
		sqlStr = []string{}
		order  = "rpt.ctime"
		state  = int32(0)
		start  = int32(0)
		end    = int32(0)
	)
	convey.Convey("ReportInfoList", t, func(ctx convey.C) {
		res, tagIDs, rptIDs, oids, err := d.ReportInfoList(c, sqlStr, order, state, start, end)
		ctx.Convey("Then err should be nil.res,tagIDs,rptIDs,oids should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(oids, convey.ShouldHaveLength, 0)
			ctx.So(rptIDs, convey.ShouldHaveLength, 0)
			ctx.So(tagIDs, convey.ShouldHaveLength, 0)
			ctx.So(res, convey.ShouldHaveLength, 0)
		})
	})
}
