package web

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestWebUgcIncre(t *testing.T) {
	convey.Convey("UgcIncre", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			pn    = int(1)
			ps    = int(10)
			start = int64(1505876448)
			end   = int64(1505876450)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.UgcIncre(c, pn, ps, start, end)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}
