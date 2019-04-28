package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoCookieCache(t *testing.T) {
	convey.Convey("CookieCache", t, func(ctx convey.C) {
		var (
			c       = context.Background()
			session = ""
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			res, err := d.CookieCache(c, session)
			ctx.Convey("Then res should be nil.", func(ctx convey.C) {
				ctx.So(res, convey.ShouldBeNil)
			})
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoDelCookieCache(t *testing.T) {
	convey.Convey("DelCookieCache", t, func(ctx convey.C) {
		var (
			c       = context.Background()
			session = ""
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := d.DelCookieCache(c, session)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoTokenCache(t *testing.T) {
	convey.Convey("TokenCache", t, func(ctx convey.C) {
		var (
			c       = context.Background()
			session = ""
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			res, err := d.TokenCache(c, session)
			ctx.Convey("Then res should be nil.", func(ctx convey.C) {
				ctx.So(res, convey.ShouldBeNil)
			})
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoDelTokenCache(t *testing.T) {
	convey.Convey("DelTokenCache", t, func(ctx convey.C) {
		var (
			c       = context.Background()
			session = ""
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := d.DelTokenCache(c, session)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}
