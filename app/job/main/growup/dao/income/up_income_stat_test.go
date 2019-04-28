package income

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestIncomeUpIncomeStat(t *testing.T) {
	convey.Convey("UpIncomeStat", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			id    = int64(0)
			limit = int64(100)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			d.db.Exec(c, "INSERT INTO up_income(mid,date) VALUES(1,'2018-06-24') ON DUPLICATE KEY UPDATE date=VALUES(date)")
			m, last, err := d.UpIncomeStat(c, id, limit)
			ctx.Convey("Then err should be nil.m,last should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(last, convey.ShouldNotBeNil)
				ctx.So(m, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestIncomeInsertUpIncomeStat(t *testing.T) {
	convey.Convey("InsertUpIncomeStat", t, func(ctx convey.C) {
		var (
			c      = context.Background()
			values = "(1,2,3,4,5)"
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			rows, err := d.InsertUpIncomeStat(c, values)
			ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rows, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestIncomeFixInsertUpIncomeStat(t *testing.T) {
	convey.Convey("FixInsertUpIncomeStat", t, func(ctx convey.C) {
		var (
			c      = context.Background()
			values = "(1,2,3,4)"
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			rows, err := d.FixInsertUpIncomeStat(c, values)
			ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rows, convey.ShouldNotBeNil)
			})
		})
	})
}
