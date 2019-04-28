package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoUserInfoByMid(t *testing.T) {
	convey.Convey("UserInfoByMid", t, func(ctx convey.C) {
		var (
			c  = context.Background()
			id = int64(27515308)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			res, err := d.UserInfoByMid(c, id)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}
