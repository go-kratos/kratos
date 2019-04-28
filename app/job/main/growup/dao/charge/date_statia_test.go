package charge

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestChargeInsertStatisTable(t *testing.T) {
	convey.Convey("InsertStatisTable", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			table = "av_charge_daily_statis"
			vals  = "(100,1,'100-200',100,1,'2018-06-24')"
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			rows, err := d.InsertStatisTable(c, table, vals)
			ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rows, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestChargeDelStatisTable(t *testing.T) {
	convey.Convey("DelStatisTable", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			table = "av_charge_daily_statis"
			date  = "2018-06-24"
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			rows, err := d.DelStatisTable(c, table, date)
			ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rows, convey.ShouldNotBeNil)
			})
		})
	})
}
