package dao

import (
	"context"
	"go-common/app/service/main/passport-auth/model"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaockKey(t *testing.T) {
	var (
		session = "9f1c9145,1536117849,c9fb62a9"
	)
	convey.Convey("ckKey", t, func(ctx convey.C) {
		p1 := ckKey(session)
		ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
			ctx.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoakKey(t *testing.T) {
	var (
		token = "b8b544c602557c27d454911d7ecc006c"
	)
	convey.Convey("akKey", t, func(ctx convey.C) {
		p1 := akKey(token)
		ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
			ctx.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func TestDaorkKey(t *testing.T) {
	var (
		refresh = "5f263d1297aa40ea0252c0963e29c6eb"
	)
	convey.Convey("rkKey", t, func(ctx convey.C) {
		p1 := rkKey(refresh)
		ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
			ctx.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoSetCookieCache(t *testing.T) {
	var (
		c       = context.TODO()
		session = "9f1c9145,1536117849,c9fb62a9"
		res     = &model.Cookie{}
	)
	convey.Convey("SetCookieCache", t, func(ctx convey.C) {
		err := d.SetCookieCache(c, session, res)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoCookieCache(t *testing.T) {
	var (
		c       = context.TODO()
		session = "9f1c9145,1536117849,c9fb62a9"
	)
	convey.Convey("CookieCache", t, func(ctx convey.C) {
		res, err := d.CookieCache(c, session)
		ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(res, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoDelCookieCache(t *testing.T) {
	var (
		c       = context.TODO()
		session = "9f1c9145,1536117849,c9fb62a9"
	)
	convey.Convey("DelCookieCache", t, func(ctx convey.C) {
		err := d.DelCookieCache(c, session)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoSetTokenCache(t *testing.T) {
	var (
		c   = context.TODO()
		k   = "5f263d1297aa40ea0252c0963e29c6eb"
		res = &model.Token{}
	)
	convey.Convey("SetTokenCache", t, func(ctx convey.C) {
		err := d.SetTokenCache(c, k, res)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoTokenCache(t *testing.T) {
	var (
		c  = context.TODO()
		sd = "5f263d1297aa40ea0252c0963e29c6eb"
	)
	convey.Convey("TokenCache", t, func(ctx convey.C) {
		res, err := d.TokenCache(c, sd)
		ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(res, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoDelTokenCache(t *testing.T) {
	var (
		c     = context.TODO()
		token = "5f263d1297aa40ea0252c0963e29c6eb"
	)
	convey.Convey("DelTokenCache", t, func(ctx convey.C) {
		err := d.DelTokenCache(c, token)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoSetRefreshCache(t *testing.T) {
	var (
		c       = context.TODO()
		refresh = &model.Refresh{
			Mid:     123,
			AppID:   430,
			Token:   "b8b544c602557c27d454911d7ecc006c",
			Refresh: "5f263d1297aa40ea0252c0963e29c6e1",
			Expires: 1850953187,
		}
	)
	convey.Convey("SetRefreshCache", t, func(ctx convey.C) {
		err := d.SetRefreshCache(c, refresh)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoRefreshCache(t *testing.T) {
	var (
		c       = context.TODO()
		refresh = "5f263d1297aa40ea0252c0963e29c6e1"
	)
	convey.Convey("RefreshCache", t, func(ctx convey.C) {
		res, err := d.RefreshCache(c, refresh)
		ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(res, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoDelRefreshCache(t *testing.T) {
	var (
		c       = context.TODO()
		refresh = "5f263d1297aa40ea0252c0963e29c6e0"
	)
	convey.Convey("DelRefreshCache", t, func(ctx convey.C) {
		err := d.DelRefreshCache(c, refresh)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaopingMC(t *testing.T) {
	var (
		c = context.TODO()
	)
	convey.Convey("pingMC", t, func(ctx convey.C) {
		err := d.pingMC(c)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}
