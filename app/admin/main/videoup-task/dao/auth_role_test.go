package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoGetUserRole(t *testing.T) {
	convey.Convey("GetUserRole", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			uid = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			role, err := d.GetUserRole(c, uid)
			ctx.Convey("Then err should be nil.role should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(role, convey.ShouldNotBeNil)
			})
		})
	})
}
