package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaouserKey(t *testing.T) {
	var (
		mid = int64(0)
	)
	convey.Convey("userKey", t, func(ctx convey.C) {
		key := userKey(mid)
		ctx.Convey("Then key should not be nil.", func(ctx convey.C) {
			ctx.So(key, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoDeleteUserBlockCache(t *testing.T) {
	var (
		c   = context.Background()
		mid = int64(0)
	)
	convey.Convey("DeleteUserBlockCache", t, func(ctx convey.C) {
		err := d.DeleteUserBlockCache(c, mid)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}
