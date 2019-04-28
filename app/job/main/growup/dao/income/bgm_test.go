package income

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestIncomeInsertBGM(t *testing.T) {
	convey.Convey("InsertBGM", t, func(ctx convey.C) {
		var (
			c      = context.Background()
			values = "(1,2,3,4,'2018-06-24','test')"
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			rows, err := d.InsertBGM(c, values)
			ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rows, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestIncomeGetBGM(t *testing.T) {
	convey.Convey("GetBGM", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			id    = int64(0)
			limit = int64(100)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			bs, last, err := d.GetBGM(c, id, limit)
			ctx.Convey("Then err should be nil.bs,last should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(last, convey.ShouldNotBeNil)
				ctx.So(bs, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestIncomeDelBGM(t *testing.T) {
	convey.Convey("DelBGM", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			limit = int64(100)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			rows, err := d.DelBGM(c, limit)
			ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rows, convey.ShouldNotBeNil)
			})
		})
	})
}
