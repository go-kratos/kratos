package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoListColumnIncome(t *testing.T) {
	convey.Convey("ListColumnIncome", t, func(ctx convey.C) {
		var (
			c         = context.Background()
			mid       = int64(1001)
			startTime = "2018-01-01"
			endTime   = "2019-01-01"
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			Exec(c, "INSERT INTO column_income(aid, mid, total_income, date) VALUES(1000, 1001, 100, '2018-06-01')")
			columns, err := d.ListColumnIncome(c, mid, startTime, endTime)
			ctx.Convey("Then err should be nil.columns should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(columns, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoListColumnIncomeByID(t *testing.T) {
	convey.Convey("ListColumnIncomeByID", t, func(ctx convey.C) {
		var (
			c       = context.Background()
			id      = int64(1000)
			endTime = "2019-01-01"
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			Exec(c, "INSERT INTO column_income(aid, mid, total_income, date) VALUES(1000, 1001, 100, '2018-06-01')")
			columns, err := d.ListColumnIncomeByID(c, id, endTime)
			ctx.Convey("Then err should be nil.columns should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(columns, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoGetColumnTitle(t *testing.T) {
	convey.Convey("GetColumnTitle", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			ids = []int64{1000}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			Exec(c, "INSERT INTO column_income_statis(aid, title) VALUES(1000, 'test')")
			titles, err := d.GetColumnTitle(c, ids)
			ctx.Convey("Then err should be nil.titles should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(titles, convey.ShouldNotBeNil)
			})
		})
	})
}
