package income

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestIncomeColumnIncomeStat(t *testing.T) {
	convey.Convey("ColumnIncomeStat", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			id    = int64(0)
			limit = int64(100)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			m, last, err := d.ColumnIncomeStat(c, id, limit)
			ctx.Convey("Then err should be nil.m,last should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(last, convey.ShouldNotBeNil)
				ctx.So(m, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestIncomeInsertColumnIncomeStat(t *testing.T) {
	convey.Convey("InsertColumnIncomeStat", t, func(ctx convey.C) {
		var (
			c      = context.Background()
			values = "(1,'test',12,2,'2018-06-24',100)"
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			rows, err := d.InsertColumnIncomeStat(c, values)
			ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rows, convey.ShouldNotBeNil)
			})
		})
	})
}
