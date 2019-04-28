package charge

import (
	"context"
	"testing"
	"time"

	"github.com/smartystreets/goconvey/convey"
)

func TestChargeColumnCharge(t *testing.T) {
	convey.Convey("ColumnCharge", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			date  = time.Date(2018, 6, 24, 0, 0, 0, 0, time.Local)
			id    = int64(0)
			limit = int(100)
			table = "column_weekly_charge"
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			d.db.Exec(c, "INSERT INTO column_weekly_charge(aid,mid,date) VALUES(1,2,'2018-06-24')")
			columns, err := d.ColumnCharge(c, date, id, limit, table)
			ctx.Convey("Then err should be nil.columns should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(columns, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestChargeCmStatis(t *testing.T) {
	convey.Convey("CmStatis", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			id    = int64(0)
			limit = int(100)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			d.db.Exec(c, "INSERT INTO column_charge_statis(aid,mid) VALUES(1,2)")
			columns, err := d.CmStatis(c, id, limit)
			ctx.Convey("Then err should be nil.columns should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(columns, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestChargeInsertCmChargeTable(t *testing.T) {
	convey.Convey("InsertCmChargeTable", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			vals  = "(1,2,3,100,'2018-06-24','2018-06-24')"
			table = "column_monthly_charge"
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			rows, err := d.InsertCmChargeTable(c, vals, table)
			ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rows, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestChargeInsertCmStatisBatch(t *testing.T) {
	convey.Convey("InsertCmStatisBatch", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			vals = "(1,2,3,4,'2018-06-24')"
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			rows, err := d.InsertCmStatisBatch(c, vals)
			ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rows, convey.ShouldNotBeNil)
			})
		})
	})
}
