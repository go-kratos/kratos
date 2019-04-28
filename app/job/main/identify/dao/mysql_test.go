package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoCookieDeleted(t *testing.T) {
	convey.Convey("CookieDeleted", t, func(ctx convey.C) {
		var (
			c      = context.Background()
			suffix = "201810"
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			res, err := d.CookieDeleted(c, 0, 100, suffix)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoTokenDeleted(t *testing.T) {
	convey.Convey("TokenDeleted", t, func(ctx convey.C) {
		var (
			c      = context.Background()
			suffix = "201810"
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			res, err := d.TokenDeleted(c, 0, 100, suffix)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDao_encodeSession(t *testing.T) {
	convey.Convey("encodeSession", t, func(ctx convey.C) {
		res := encodeSession([]byte{1})
		ctx.Convey("Then err should not be nil.", func(ctx convey.C) {
			ctx.So(res, convey.ShouldNotBeNil)
		})
	})
}
