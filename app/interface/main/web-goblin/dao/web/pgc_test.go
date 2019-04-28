package web

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestWebPgcFull(t *testing.T) {
	convey.Convey("PgcFull", t, func(ctx convey.C) {
		var (
			c      = context.Background()
			tp     = int(1)
			pn     = int64(1)
			ps     = int64(10)
			source = "youku"
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.PgcFull(c, tp, pn, ps, source)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestWebPgcIncre(t *testing.T) {
	convey.Convey("PgcIncre", t, func(ctx convey.C) {
		var (
			c      = context.Background()
			tp     = int(2)
			pn     = int64(1)
			ps     = int64(10)
			start  = int64(1505876448)
			end    = int64(1505876450)
			source = "youku"
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.PgcIncre(c, tp, pn, ps, start, end, source)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}
