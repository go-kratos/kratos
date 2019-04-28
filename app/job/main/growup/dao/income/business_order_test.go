package income

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestIncomeBusinessOrders(t *testing.T) {
	convey.Convey("BusinessOrders", t, func(ctx convey.C) {
		var (
			c      = context.Background()
			offset = int64(0)
			limit  = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			last, m, err := d.BusinessOrders(c, offset, limit)
			ctx.Convey("Then err should not be nil.last,m should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(m, convey.ShouldNotBeNil)
				ctx.So(last, convey.ShouldNotBeNil)
			})
		})
	})
}
