package roomAdmin

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestRoomAdminAddCacheNoneUser(t *testing.T) {
	convey.Convey("AddCacheNoneUser", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			uid = int64(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := d.AddCacheNoneUser(c, uid)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestRoomAdminAddCacheNoneRoom(t *testing.T) {
	convey.Convey("AddCacheNoneRoom", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			uid = int64(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := d.AddCacheNoneRoom(c, uid)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}
