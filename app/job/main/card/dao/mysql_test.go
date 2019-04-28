package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoUpdateExpireTime(t *testing.T) {
	convey.Convey("UpdateExpireTime", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			no  = int64(0)
			mid = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.UpdateExpireTime(c, no, mid)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}
