package dao

import (
	"context"
	"testing"

	"go-common/app/interface/main/dm/model"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaorptTable(t *testing.T) {
	convey.Convey("rptTable", t, func(convCtx convey.C) {
		var (
			cid = int64(0)
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			p1 := rptTable(cid)
			convCtx.Convey("Then p1 should not be nil.", func(convCtx convey.C) {
				convCtx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaorptUserTable(t *testing.T) {
	convey.Convey("rptUserTable", t, func(convCtx convey.C) {
		var (
			dmid = int64(0)
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			p1 := rptUserTable(dmid)
			convCtx.Convey("Then p1 should not be nil.", func(convCtx convey.C) {
				convCtx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoRptLogTable(t *testing.T) {
	convey.Convey("RptLogTable", t, func(convCtx convey.C) {
		var (
			dmid = int64(0)
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			p1 := RptLogTable(dmid)
			convCtx.Convey("Then p1 should not be nil.", func(convCtx convey.C) {
				convCtx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoAddReport(t *testing.T) {
	convey.Convey("AddReport", t, func(convCtx convey.C) {
		var (
			c   = context.Background()
			rpt = &model.Report{}
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			id, err := testDao.AddReport(c, rpt)
			convCtx.Convey("Then err should be nil.id should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(id, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoAddReportUser(t *testing.T) {
	convey.Convey("AddReportUser", t, func(convCtx convey.C) {
		var (
			c = context.Background()
			u = &model.User{}
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			id, err := testDao.AddReportUser(c, u)
			convCtx.Convey("Then err should be nil.id should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(id, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoReportLog(t *testing.T) {
	convey.Convey("ReportLog", t, func(convCtx convey.C) {
		var (
			c    = context.Background()
			dmid = int64(0)
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			res, err := testDao.ReportLog(c, dmid)
			convCtx.Convey("Then err should be nil.res should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoReport(t *testing.T) {
	convey.Convey("Report", t, func(convCtx convey.C) {
		var (
			c    = context.Background()
			cid  = int64(0)
			dmid = int64(0)
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			rpt, err := testDao.Report(c, cid, dmid)
			convCtx.Convey("Then err should be nil.rpt should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(rpt, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoUpdateReportStat(t *testing.T) {
	convey.Convey("UpdateReportStat", t, func(convCtx convey.C) {
		var (
			c     = context.Background()
			cid   = int64(0)
			dmid  = int64(0)
			state = int8(0)
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			affect, err := testDao.UpdateReportStat(c, cid, dmid, state)
			convCtx.Convey("Then err should be nil.affect should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(affect, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoUpdateReportUPOp(t *testing.T) {
	convey.Convey("UpdateReportUPOp", t, func(convCtx convey.C) {
		var (
			c    = context.Background()
			cid  = int64(0)
			dmid = int64(0)
			op   = int8(0)
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			affect, err := testDao.UpdateReportUPOp(c, cid, dmid, op)
			convCtx.Convey("Then err should be nil.affect should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(affect, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoAddReportLog(t *testing.T) {
	convey.Convey("AddReportLog", t, func(convCtx convey.C) {
		var (
			c  = context.Background()
			lg = &model.RptLog{}
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			err := testDao.AddReportLog(c, lg)
			convCtx.Convey("Then err should be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoReportUser(t *testing.T) {
	convey.Convey("ReportUser", t, func(convCtx convey.C) {
		var (
			c    = context.Background()
			dmid = int64(123)
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			_, err := testDao.ReportUser(c, dmid)
			convCtx.Convey("Then err should be nil.users should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoSetReportUserFinished(t *testing.T) {
	convey.Convey("SetReportUserFinished", t, func(convCtx convey.C) {
		var (
			c    = context.Background()
			dmid = int64(0)
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			err := testDao.SetReportUserFinished(c, dmid)
			convCtx.Convey("Then err should be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
			})
		})
	})
}
