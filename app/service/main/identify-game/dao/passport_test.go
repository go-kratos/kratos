package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoAccessToken(t *testing.T) {
	var (
		c         = context.Background()
		accesskey = "accessKey"
		target    = "origin"
	)
	convey.Convey("AccessToken", t, func(ctx convey.C) {
		token, err := d.AccessToken(c, accesskey, target)
		ctx.Convey("Then err should be nil.token should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldEqual, -101)
			ctx.So(token, convey.ShouldBeNil)
		})
	})
}

func TestDaoRenewToken(t *testing.T) {
	var (
		c         = context.Background()
		accesskey = "accessKey"
		target    = "origin"
	)
	convey.Convey("RenewToken", t, func(ctx convey.C) {
		renewToken, err := d.RenewToken(c, accesskey, target)
		ctx.Convey("Then err should be nil.renewToken should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldEqual, -101)
			ctx.So(renewToken, convey.ShouldBeNil)
		})
	})
}

func TestGetCookieByMid(t *testing.T) {
	var (
		c   = context.Background()
		mid = int64(234112334566723432)
	)
	convey.Convey("GetCookieByMid", t, func(ctx convey.C) {
		cookies, err := d.GetCookieByMid(c, mid)
		convey.So(err, convey.ShouldNotBeNil)
		convey.So(cookies, convey.ShouldBeNil)
	})
}
