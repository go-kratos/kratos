package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoFansCount(t *testing.T) {
	convey.Convey("FansCount", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			fans, err := d.FansCount(c, mid)
			ctx.Convey("Then err should be nil.fans should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(fans, convey.ShouldNotBeNil)
			})
		})
	})
}
