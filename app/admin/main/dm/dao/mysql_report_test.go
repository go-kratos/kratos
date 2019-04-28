package dao

import (
	"context"
	"go-common/app/admin/main/dm/model"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoRptTable(t *testing.T) {
	convey.Convey("RptTable", t, func(ctx convey.C) {
		var (
			cid = int64(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			p1 := RptTable(cid)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoUserTable(t *testing.T) {
	convey.Convey("UserTable", t, func(ctx convey.C) {
		var (
			dmid = int64(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			p1 := UserTable(dmid)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoLogTable(t *testing.T) {
	convey.Convey("LogTable", t, func(ctx convey.C) {
		var (
			dmid = int64(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			p1 := LogTable(dmid)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoChangeReportStat(t *testing.T) {
	convey.Convey("ChangeReportStat", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			cid   = int64(0)
			dmids = []int64{}
			state = int8(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			testDao.ChangeReportStat(c, cid, dmids, state)
		})
	})
}

func TestDaoIgnoreReport(t *testing.T) {
	convey.Convey("IgnoreReport", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			cid   = int64(0)
			dmids = []int64{}
			state = int8(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			testDao.IgnoreReport(c, cid, dmids, state)
		})
	})
}

func TestDaoReports(t *testing.T) {
	convey.Convey("Reports", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			cid   = int64(0)
			dmids = []int64{}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			testDao.Reports(c, cid, dmids)
		})
	})
}

func TestDaoReportUsers(t *testing.T) {
	convey.Convey("ReportUsers", t, func(ctx convey.C) {
		var (
			c       = context.Background()
			tableID = int64(0)
			dmids   = []int64{}
			state   = int8(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			testDao.ReportUsers(c, tableID, dmids, state)
		})
	})
}

func TestDaoUpReportUserState(t *testing.T) {
	convey.Convey("UpReportUserState", t, func(ctx convey.C) {
		var (
			c       = context.Background()
			tableID = int64(0)
			dmids   = []int64{}
			state   = int8(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			testDao.UpReportUserState(c, tableID, dmids, state)
		})
	})
}

func TestDaoAddReportLog(t *testing.T) {
	convey.Convey("AddReportLog", t, func(ctx convey.C) {
		var (
			c       = context.Background()
			tableID = int64(0)
			lg      = []*model.ReportLog{}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			testDao.AddReportLog(c, tableID, lg)
		})
	})
}

func TestDaoReportLog(t *testing.T) {
	convey.Convey("ReportLog", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			dmid = int64(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			res, err := testDao.ReportLog(c, dmid)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}
