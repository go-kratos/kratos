package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoGetAvBreach(t *testing.T) {
	convey.Convey("GetAvBreach", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			start = "2018-06-20"
			end   = "2018-10-20"
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			avs, err := d.GetAvBreach(c, start, end)
			ctx.Convey("Then err should be nil.avs should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(avs, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoGetAvBreachPre(t *testing.T) {
	convey.Convey("GetAvBreachPre", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			ctype = int(1)
			state = int(3)
			cdate = "2018-10-23"
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			avs, err := d.GetAvBreachPre(c, ctype, state, cdate)
			ctx.Convey("Then err should be nil.avs should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(avs, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoInsertAvBreachPre(t *testing.T) {
	convey.Convey("InsertAvBreachPre", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			val = "100, 29, '2018-06-23', 1, 3"
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			Exec(c, "DELETE FROM av_breach_pre WHERE aid=100")
			rows, err := d.InsertAvBreachPre(c, val)
			ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rows, convey.ShouldNotBeNil)
			})
		})
	})
}
