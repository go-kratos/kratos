package income

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestIncomeUpIncomeCount(t *testing.T) {
	convey.Convey("UpIncomeCount", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			table = "up_income"
			query = "id > 0"
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			Exec(c, "INSERT INTO up_income(mid, income, date) VALUS(1993, 10, '2018-01-01')")
			count, err := d.UpIncomeCount(c, table, query)
			ctx.Convey("Then err should be nil.count should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(count, convey.ShouldNotBeNil)
			})
		})
	})

	convey.Convey("UpIncomeCount table error", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			table = ""
			query = ""
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			_, err := d.UpIncomeCount(c, table, query)
			ctx.Convey("Then err should be nil.count should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestIncomeGetUpIncome(t *testing.T) {
	convey.Convey("GetUpIncome", t, func(ctx convey.C) {
		var (
			c          = context.Background()
			table      = "up_income"
			incomeType = "av_income"
			query      = "is_deleted = 0"
			id         = int64(0)
			limit      = int(10)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			upIncome, err := d.GetUpIncome(c, table, incomeType, query, id, limit)
			ctx.Convey("Then err should be nil.upIncome should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(upIncome, convey.ShouldNotBeNil)
			})
		})
	})

	convey.Convey("GetUpIncome query == nil", t, func(ctx convey.C) {
		var (
			c          = context.Background()
			table      = "up_income"
			incomeType = "av_income"
			query      = ""
			id         = int64(0)
			limit      = int(10)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			_, err := d.GetUpIncome(c, table, incomeType, query, id, limit)
			ctx.Convey("Then err should be nil.upIncome should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestIncomeGetUpIncomeBySort(t *testing.T) {
	convey.Convey("GetUpIncomeBySort", t, func(ctx convey.C) {
		var (
			c         = context.Background()
			table     = "up_income"
			typeField = "av_income,av_tax,av_base_income,av_total_income"
			sort      = "id"
			query     = "id > 0"
			from      = int(0)
			limit     = int(10)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			upIncome, err := d.GetUpIncomeBySort(c, table, typeField, sort, query, from, limit)
			ctx.Convey("Then err should be nil.upIncome should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(upIncome, convey.ShouldNotBeNil)
			})
		})
	})

	convey.Convey("GetUpIncomeBySort table == nil", t, func(ctx convey.C) {
		var (
			c         = context.Background()
			table     = ""
			typeField = "av_income"
			sort      = "id"
			query     = "id > 0"
			from      = int(0)
			limit     = int(10)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			_, err := d.GetUpIncomeBySort(c, table, typeField, sort, query, from, limit)
			ctx.Convey("Then err should be nil.upIncome should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestIncomeGetUpDailyStatis(t *testing.T) {
	convey.Convey("GetUpDailyStatis", t, func(ctx convey.C) {
		var (
			c        = context.Background()
			table    = "up_income_daily_statis"
			fromTime = "2018-01-01"
			toTime   = "2018-01-10"
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			Exec(c, "INSERT INTO up_income_daily_statis(ups, cdate) VALUES(10, '2018-01-02')")
			s, err := d.GetUpDailyStatis(c, table, fromTime, toTime)
			ctx.Convey("Then err should be nil.s should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(s, convey.ShouldNotBeNil)
			})
		})
	})
}
