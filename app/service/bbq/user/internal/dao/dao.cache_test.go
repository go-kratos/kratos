package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoUserBase(t *testing.T) {
	convey.Convey("UserBase", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			keys = []int64{88895104, 88895105}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.UserBase(c, keys)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}
