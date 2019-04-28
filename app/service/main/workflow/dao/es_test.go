package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoUpdateUserTag(t *testing.T) {
	convey.Convey("UpdateUserTag", t, func(ctx convey.C) {
		var (
			c       = context.Background()
			gid     = int64(1)
			userTid = int32(1)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := d.UpdateUserTag(c, gid, userTid)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}
