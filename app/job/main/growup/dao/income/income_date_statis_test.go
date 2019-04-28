package income

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestIncomeInsertIncomeStatisTable(t *testing.T) {
	convey.Convey("InsertIncomeStatisTable", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			table = "av_income_daily_statis"
			vals  = "(1,2,'100-200',100,1,'2018-06-24')"
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			rows, err := d.InsertIncomeStatisTable(c, table, vals)
			ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rows, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestIncomeInsertUpIncomeDailyStatis(t *testing.T) {
	convey.Convey("InsertUpIncomeDailyStatis", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			table = "up_income_daily_statis"
			vals  = "(10,1,'100-200',1,'2018-06-24')"
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			rows, err := d.InsertUpIncomeDailyStatis(c, table, vals)
			ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rows, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestIncomeDelIncomeStatisTable(t *testing.T) {
	convey.Convey("DelIncomeStatisTable", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			table = "av_income_daily_statis"
			date  = "2018-06-24"
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			rows, err := d.DelIncomeStatisTable(c, table, date)
			ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rows, convey.ShouldNotBeNil)
			})
		})
	})
}
