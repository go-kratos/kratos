package income

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestIncomeBgmIncomeStat(t *testing.T) {
	convey.Convey("BgmIncomeStat", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			id    = int64(0)
			limit = int64(100)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			m, last, err := d.BgmIncomeStat(c, id, limit)
			ctx.Convey("Then err should be nil.m,last should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(last, convey.ShouldNotBeNil)
				ctx.So(m, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestIncomeInsertBgmIncomeStat(t *testing.T) {
	convey.Convey("InsertBgmIncomeStat", t, func(ctx convey.C) {
		var (
			c      = context.Background()
			values = "(100,200)"
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			rows, err := d.InsertBgmIncomeStat(c, values)
			ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rows, convey.ShouldNotBeNil)
			})
		})
	})
}
