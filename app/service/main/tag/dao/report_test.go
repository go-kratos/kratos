package dao

import (
	"context"
	"testing"

	"go-common/app/service/main/tag/model"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoReport(t *testing.T) {
	convey.Convey("Report", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			oid = int64(28843596)
			typ = int32(3)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.Report(c, oid, typ)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoTxAddReport(t *testing.T) {
	rpt := &model.Report{
		Oid:    28843596,
		Type:   3,
		Tid:    1833,
		Mid:    35152246,
		TypeID: 1,
		Action: 0,
	}
	tx, err := d.BeginTran(context.Background())
	if err != nil {
		return
	}
	convey.Convey("TxAddReport", t, func(ctx convey.C) {
		rptID, err := d.TxAddReport(tx, rpt)
		ctx.Convey("TxAddReport,Then err should be nil.rptID should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(rptID, convey.ShouldBeGreaterThanOrEqualTo, 0)
		})
		r := &model.ReportUser{
			RptID: rptID,
			Mid:   35152246,
		}
		rptID, err = d.TxAddUserReport(tx, r)
		ctx.Convey("TxAddUserReport,Then err should be nil.rptID should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(rptID, convey.ShouldBeGreaterThanOrEqualTo, 0)
		})
	})
	tx.Commit()
}

func TestDaoAddUserReport(t *testing.T) {
	convey.Convey("AddUserReport", t, func(ctx convey.C) {
		var (
			c = context.Background()
			r = &model.ReportUser{
				RptID: 2233,
				Mid:   35152246,
			}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1, err := d.AddUserReport(c, r)
			ctx.Convey("Then err should be nil.p1 should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoReportAndUser(t *testing.T) {
	convey.Convey("ReportAndUser", t, func(ctx convey.C) {
		var (
			c      = context.Background()
			oid    = int64(28843596)
			mid    = int64(35152246)
			tid    = int64(1833)
			typ    = int32(3)
			action = int32(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			d.ReportAndUser(c, oid, mid, tid, typ, action)
		})
	})
}

func TestDaoReportUser(t *testing.T) {
	convey.Convey("ReportUser", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			lid = int64(12345)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.ReportUser(c, lid)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}
