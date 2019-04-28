package dao

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoGetDayExpenseCount(t *testing.T) {
	convey.Convey("GetDayExpenseCount", t, func(ctx convey.C) {
		var (
			c         = context.Background()
			beginDate = time.Now()
			ctype     = int(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			total, err := d.GetDayExpenseCount(c, beginDate, ctype)
			ctx.Convey("Then err should be nil.total should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(total, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoGetAllDayExpenseInfo(t *testing.T) {
	convey.Convey("GetAllDayExpenseInfo", t, func(ctx convey.C) {
		var (
			c            = context.Background()
			beginDate, _ = time.ParseInLocation("2018-06-01", "2018-01-01", time.Local)
			ctype        = int(0)
			from         = int(0)
			limit        = int(10)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			Exec(c, "INSERT INTO expense_daily_info(up_acount, date) VALUES(100, '2018-01-01')")
			_, err := d.GetAllDayExpenseInfo(c, beginDate, ctype, from, limit)
			ctx.Convey("Then err should be nil.infos should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoGetDayTotalExpenseInfo(t *testing.T) {
	convey.Convey("GetDayTotalExpenseInfo", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			date  = time.Date(2018, 9, 1, 0, 0, 0, 0, time.Local)
			ctype = int(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			Exec(c, fmt.Sprintf("INSERT INTO expense_daily_info(date, total_expense) VALUES('%s', 100)", date.Format("2006-01-02")))
			_, err := d.GetDayTotalExpenseInfo(c, date, ctype)
			ctx.Convey("Then err should be nil.totalExpense should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoGetMonthExpenseCount(t *testing.T) {
	convey.Convey("GetMonthExpenseCount", t, func(ctx convey.C) {
		var (
			c          = context.Background()
			month      = "2019-01-01"
			beginMonth = "2018-01-01"
			ctype      = int(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			total, err := d.GetMonthExpenseCount(c, month, beginMonth, ctype)
			ctx.Convey("Then err should be nil.total should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(total, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoGetAllMonthExpenseInfo(t *testing.T) {
	convey.Convey("GetAllMonthExpenseInfo", t, func(ctx convey.C) {
		var (
			c          = context.Background()
			month      = "2019-01"
			beginMonth = "2018-01"
			ctype      = int(0)
			from       = int(0)
			limit      = int(10)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			Exec(c, fmt.Sprintf("INSERT INTO expense_monthly_info(date, total_expense) VALUES('2018-02', 100)"))
			_, err := d.GetAllMonthExpenseInfo(c, month, beginMonth, ctype, from, limit)
			ctx.Convey("Then err should be nil.infos should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoGetLatelyExpenseDate(t *testing.T) {
	convey.Convey("GetLatelyExpenseDate", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			table = "daily"
			ctype = int(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			date, err := d.GetLatelyExpenseDate(c, table, ctype)
			ctx.Convey("Then err should be nil.date should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(date, convey.ShouldNotBeNil)
			})
		})
	})
}
