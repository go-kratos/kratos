package dao

import (
	"context"
	"go-common/app/job/main/growup/model"
	"testing"
	"time"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoGetTotalExpenseByDate(t *testing.T) {
	convey.Convey("GetTotalExpenseByDate", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			date = "2018-06-23"
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			expense, err := d.GetTotalExpenseByDate(c, date)
			ctx.Convey("Then err should be nil.expense should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(expense, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoInsertDailyExpense(t *testing.T) {
	convey.Convey("InsertDailyExpense", t, func(ctx convey.C) {
		var (
			c = context.Background()
			e = &model.BudgetExpense{
				Expense:      100,
				UpCount:      100,
				AvCount:      100,
				UpAvgExpense: 100,
				AvAvgExpense: 100,
				Date:         time.Now(),
				TotalExpense: 100,
				CType:        1,
			}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			rows, err := d.InsertDailyExpense(c, e)
			ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rows, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoInsertMonthlyExpense(t *testing.T) {
	convey.Convey("InsertMonthlyExpense", t, func(ctx convey.C) {
		var (
			c = context.Background()
			e = &model.BudgetExpense{
				Expense:      100,
				UpCount:      100,
				AvCount:      100,
				UpAvgExpense: 100,
				AvAvgExpense: 100,
				Date:         time.Now(),
				TotalExpense: 100,
				CType:        1,
			}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			rows, err := d.InsertMonthlyExpense(c, e)
			ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rows, convey.ShouldNotBeNil)
			})
		})
	})
}
