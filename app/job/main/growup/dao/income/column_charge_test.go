package income

import (
	"context"
	"testing"
	"time"

	"github.com/smartystreets/goconvey/convey"
)

func TestIncomeColumnDailyCharge(t *testing.T) {
	convey.Convey("ColumnDailyCharge", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			date  = time.Now()
			id    = int64(0)
			limit = int(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			columns, err := d.ColumnDailyCharge(c, date, id, limit)
			ctx.Convey("Then err should be nil.columns should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(columns, convey.ShouldNotBeNil)
			})
		})
	})
}
