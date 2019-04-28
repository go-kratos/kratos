package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoAccessCookie(t *testing.T) {
	var (
		c      = context.Background()
		cookie = "cookie=cookie"
	)
	convey.Convey("AccessCookie", t, func(ctx convey.C) {
		res, err := d.AccessCookie(c, cookie)
		ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldNotBeNil)
			ctx.So(res, convey.ShouldBeNil)
		})
	})
}

func TestDaoAccessToken(t *testing.T) {
	var (
		c         = context.Background()
		accesskey = "123456"
	)
	convey.Convey("AccessToken", t, func(ctx convey.C) {
		res, err := d.AccessToken(c, accesskey)
		ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldNotBeNil)
			ctx.So(res, convey.ShouldBeNil)
		})
	})
}
