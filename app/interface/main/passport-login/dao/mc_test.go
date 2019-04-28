package dao

import (
	"context"
	"testing"
	"time"

	"go-common/app/interface/main/passport-login/model"

	"github.com/smartystreets/goconvey/convey"
)

func TestDao_SetCookieCache(t *testing.T) {
	var (
		now    = time.Now()
		c      = context.Background()
		cookie = &model.CookieProto{
			Mid:     1,
			Session: "aac24b1cccebc4c85fcf5c3b65116ba8",
			CSRF:    "aac24b1cccebc4c85fcf5c3b65116ba8",
			Type:    0,
			Expires: now.Unix() + _expireSeconds,
		}
	)
	convey.Convey("SetCookieCache", t, func(ctx convey.C) {
		err := d.SetCookieCache(c, cookie)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDao_CookieCache(t *testing.T) {
	var (
		c       = context.Background()
		session = "aac24b1cccebc4c85fcf5c3b65116ba8"
	)
	convey.Convey("CookieCache", t, func(ctx convey.C) {
		res, err := d.CookieCache(c, session)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(res, convey.ShouldNotBeNil)
		})
	})
}

func TestDao_DelCookieCache(t *testing.T) {
	var (
		c       = context.Background()
		session = "aac24b1cccebc4c85fcf5c3b65116ba8"
	)
	convey.Convey("DelCookieCache", t, func(ctx convey.C) {
		err := d.DelCookieCache(c, session)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDao_SetTokenCache(t *testing.T) {
	var (
		now   = time.Now()
		c     = context.Background()
		token = &model.TokenProto{
			Mid:     1,
			Token:   "aac24b1cccebc4c85fcf5c3b65116ba8",
			Type:    0,
			Expires: now.Unix() + _expireSeconds,
		}
	)
	convey.Convey("SetTokenCache", t, func(ctx convey.C) {
		err := d.SetTokenCache(c, token)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDao_TokenCache(t *testing.T) {
	var (
		c     = context.Background()
		token = "aac24b1cccebc4c85fcf5c3b65116ba8"
	)
	convey.Convey("TokenCache", t, func(ctx convey.C) {
		res, err := d.TokenCache(c, token)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(res, convey.ShouldNotBeNil)
		})
	})
}

func TestDao_SetRefreshCache(t *testing.T) {
	var (
		now     = time.Now()
		c       = context.Background()
		refresh = &model.RefreshProto{
			Mid:     1,
			Token:   "aac24b1cccebc4c85fcf5c3b65116ba8",
			Refresh: "aac24b1cccebc4c85fcf5c3b65116ba8",
			Expires: now.Unix() + _expireSeconds,
		}
	)
	convey.Convey("SetRefreshCache", t, func(ctx convey.C) {
		err := d.SetRefreshCache(c, refresh)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDao_DelTokenCache(t *testing.T) {
	var (
		c     = context.Background()
		token = "aac24b1cccebc4c85fcf5c3b65116ba8"
	)
	convey.Convey("DelTokenCache", t, func(ctx convey.C) {
		err := d.DelTokenCache(c, token)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}
