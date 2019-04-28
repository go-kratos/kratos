package charge

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestChargeGetUpCharges(t *testing.T) {
	convey.Convey("GetUpCharges", t, func(ctx convey.C) {
		var (
			c      = context.Background()
			table  = "up_daily_charge"
			date   = "2018-06-24"
			offset = int64(0)
			limit  = int64(100)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			d.db.Exec(c, "INSERT INTO up_daily_charge(mid,date) VALUES(1, '2018-06-24')")
			last, charges, err := d.GetUpCharges(c, table, date, offset, limit)
			ctx.Convey("Then err should be nil.last,charges should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(charges, convey.ShouldNotBeNil)
				ctx.So(last, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestChargeInsertUpCharge(t *testing.T) {
	convey.Convey("InsertUpCharge", t, func(ctx convey.C) {
		var (
			c      = context.Background()
			table  = "up_daily_charge"
			values = "(1,2,3,'2018-06-24')"
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			rows, err := d.InsertUpCharge(c, table, values)
			ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rows, convey.ShouldNotBeNil)
			})
		})
	})
}
