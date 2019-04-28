package income

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestIncomeInsertBgmIncome(t *testing.T) {
	convey.Convey("InsertBgmIncome", t, func(ctx convey.C) {
		var (
			c      = context.Background()
			values = "(1,2,3,4,5,6,7,'2018-06-24',100,100)"
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			rows, err := d.InsertBgmIncome(c, values)
			ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rows, convey.ShouldNotBeNil)
			})
		})
	})
}
