package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoAbServer(t *testing.T) {
	convey.Convey("AbServer", t, func(ctx convey.C) {
		var (
			c        = context.Background()
			mid      = int64(0)
			platform = int(0)
			channel  = "test"
			buvid    = "test"
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.AbServer(c, mid, platform, channel, buvid)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}
