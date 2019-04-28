package income

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestIncomeInsertColumnIncome(t *testing.T) {
	convey.Convey("InsertColumnIncome", t, func(ctx convey.C) {
		var (
			c      = context.Background()
			values = "(1,2,3,'2018-06-24',100,10,100,10,'2018-06-24',100)"
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			rows, err := d.InsertColumnIncome(c, values)
			ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rows, convey.ShouldNotBeNil)
			})
		})
	})
}
