package block

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestBlockuserKey(t *testing.T) {
	convey.Convey("userKey", t, func(ctx convey.C) {
		var (
			mid = int64(46333)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			key := userKey(mid)
			ctx.Convey("Then key should not be nil.", func(ctx convey.C) {
				ctx.So(key, convey.ShouldEqual, "u_46333")
			})
		})
	})
}

func TestBlockDeleteUserCache(t *testing.T) {
	convey.Convey("DeleteUserCache", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(1)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.DeleteUserCache(c, mid)
			ctx.Convey("test DeleteUserCache", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}
