package dao

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"testing"
	"time"

	"go-common/app/interface/main/passport-login/model"

	"github.com/smartystreets/goconvey/convey"
)

const (
	_expireSeconds = 2592000 // 30 days
)

func TestDaoAddCookie(t *testing.T) {
	var (
		now        = time.Now()
		c          = context.TODO()
		session, _ = hex.DecodeString("25ded96af42eb61677730d0a74eb4ca1")
		csrf, _    = hex.DecodeString("aac24b1cccebc4c85fcf5c3b65116ba8")
		cookie     = &model.Cookie{
			Mid:     1,
			Session: session,
			CSRF:    csrf,
			Type:    0,
			Expires: now.Unix() + _expireSeconds,
		}
	)
	convey.Convey("AddCookie", t, func(ctx convey.C) {
		res, err := d.AddCookie(c, cookie, now)
		ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(res, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoAddToken(t *testing.T) {
	var (
		now   = time.Now()
		c     = context.TODO()
		tk, _ = hex.DecodeString("25ded96af42eb61677730d0a74eb4ca1")
		token = &model.Token{
			Mid:     1,
			AppID:   876,
			Token:   tk,
			Expires: now.Unix() + _expireSeconds,
			Type:    0,
		}
	)
	convey.Convey("AddToken", t, func(ctx convey.C) {
		res, err := d.AddToken(c, token, now)
		ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(res, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoAddRefresh(t *testing.T) {
	var (
		now     = time.Now()
		c       = context.TODO()
		tk, _   = hex.DecodeString("25ded96af42eb61677730d0a74eb4ca1")
		rk, _   = hex.DecodeString("aac24b1cccebc4c85fcf5c3b65116ba8")
		refresh = &model.Refresh{
			Mid:     1,
			AppID:   876,
			Refresh: rk,
			Token:   tk,
			Expires: now.Unix() + _expireSeconds,
		}
	)
	convey.Convey("AddRefresh", t, func(ctx convey.C) {
		res, err := d.AddRefresh(c, refresh, now)
		ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(res, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoAddOldCookie(t *testing.T) {
	var (
		now    = time.Now()
		c      = context.TODO()
		cookie = &model.OldCookie{
			Mid:       1,
			Session:   "396c38bb,1519380539,d73804f2",
			CSRFToken: "aac24b1cccebc4c85fcf5c3b65116ba8",
			Type:      0,
			Expires:   now.Unix() + _expireSeconds,
		}
	)
	convey.Convey("AddOldCookie", t, func(ctx convey.C) {
		res, err := d.AddOldCookie(c, cookie)
		ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(res, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoAddOldToken(t *testing.T) {
	var (
		now   = time.Now()
		c     = context.TODO()
		token = &model.OldToken{
			Mid:          1,
			AppID:        876,
			AccessToken:  "25ded96af42eb61677730d0a74eb4ca1",
			RefreshToken: "aac24b1cccebc4c85fcf5c3b65116ba8",
			Expires:      now.Unix() + _expireSeconds,
		}
	)
	convey.Convey("AddOldToken", t, func(ctx convey.C) {
		res, err := d.AddOldToken(c, token, now)
		ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(res, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoGetCookie(t *testing.T) {
	var (
		now        = time.Now()
		c          = context.TODO()
		session, _ = hex.DecodeString("25ded96af42eb61677730d0a74eb4ca1")
	)
	convey.Convey("GetCookie", t, func(ctx convey.C) {
		res, err := d.GetCookie(c, session, now)
		ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(res, convey.ShouldNotBeNil)
		})

		str, _ := json.Marshal(res)
		t.Logf("res: %s", str)
	})
}

func TestDaoGetToken(t *testing.T) {
	var (
		now      = time.Now()
		c        = context.TODO()
		token, _ = hex.DecodeString("25ded96af42eb61677730d0a74eb4ca1")
	)
	convey.Convey("GetToken", t, func(ctx convey.C) {
		res, err := d.GetToken(c, token, now)
		ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(res, convey.ShouldNotBeNil)
		})

		str, _ := json.Marshal(res)
		t.Logf("res: %s", str)
	})
}

func TestDaoGetRefresh(t *testing.T) {
	var (
		now        = time.Now()
		c          = context.TODO()
		refresh, _ = hex.DecodeString("aac24b1cccebc4c85fcf5c3b65116ba8")
	)
	convey.Convey("GetRefresh", t, func(ctx convey.C) {
		res, err := d.GetRefresh(c, refresh, now)
		ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(res, convey.ShouldNotBeNil)
		})

		str, _ := json.Marshal(res)
		t.Logf("res: %s", str)
	})
}

func TestDaoDelCookie(t *testing.T) {
	var (
		now        = time.Now()
		c          = context.TODO()
		session, _ = hex.DecodeString("25ded96af42eb61677730d0a74eb4ca1")
	)
	convey.Convey("DelCookie", t, func(ctx convey.C) {
		res, err := d.DelCookie(c, session, now)
		ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(res, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoDelToken(t *testing.T) {
	var (
		now      = time.Now()
		c        = context.TODO()
		token, _ = hex.DecodeString("25ded96af42eb61677730d0a74eb4ca1")
	)
	convey.Convey("DelCookie", t, func(ctx convey.C) {
		res, err := d.DelToken(c, token, now)
		ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(res, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoDelRefresh(t *testing.T) {
	var (
		now        = time.Now()
		c          = context.TODO()
		refresh, _ = hex.DecodeString("aac24b1cccebc4c85fcf5c3b65116ba8")
	)
	convey.Convey("DelRefresh", t, func(ctx convey.C) {
		res, err := d.DelRefresh(c, refresh, now)
		ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(res, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoDelOldCookie(t *testing.T) {
	var (
		c       = context.TODO()
		session = "396c38bb,1519380539,d73804f2"
		mid     = int64(1)
	)
	convey.Convey("DelOldCookie", t, func(ctx convey.C) {
		res, err := d.DelOldCookie(c, session, mid)
		ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(res, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoDelOldToken(t *testing.T) {
	var (
		now   = time.Now()
		c     = context.TODO()
		token = "25ded96af42eb61677730d0a74eb4ca1"
	)
	convey.Convey("DelOldToken", t, func(ctx convey.C) {
		res, err := d.DelOldToken(c, token, now)
		ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(res, convey.ShouldNotBeNil)
		})
	})
}

func TestDao_DelCookieByMid(t *testing.T) {
	var (
		now = time.Now()
		c   = context.TODO()
		mid = int64(1)
	)
	convey.Convey("DelCookieByMid", t, func(ctx convey.C) {
		res, err := d.DelCookieByMid(c, mid, now)
		ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(res, convey.ShouldNotBeNil)
		})
	})
}

func TestDao_DelTokenByMid(t *testing.T) {
	var (
		now = time.Now()
		c   = context.TODO()
		mid = int64(1)
	)
	convey.Convey("DelTokenByMid", t, func(ctx convey.C) {
		res, err := d.DelTokenByMid(c, mid, now)
		ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(res, convey.ShouldNotBeNil)
		})
	})
}

func TestDao_DelRefreshByMid(t *testing.T) {
	var (
		now = time.Now()
		c   = context.TODO()
		mid = int64(1)
	)
	convey.Convey("DelRefreshByMid", t, func(ctx convey.C) {
		res, err := d.DelRefreshByMid(c, mid, now)
		ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(res, convey.ShouldNotBeNil)
		})
	})
}

func TestDao_DelOldCookieByMid(t *testing.T) {
	var (
		c   = context.TODO()
		mid = int64(1)
	)
	convey.Convey("DelOldCookieByMid", t, func(ctx convey.C) {
		res, err := d.DelOldCookieByMid(c, mid)
		ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(res, convey.ShouldNotBeNil)
		})
	})
}

func TestDao_DelOldTokenByMid(t *testing.T) {
	var (
		now = time.Now()
		c   = context.TODO()
		mid = int64(1)
	)
	convey.Convey("TestDao_DelOldTokenByMid", t, func(ctx convey.C) {
		res, err := d.DelOldTokenByMid(c, mid, now)
		ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(res, convey.ShouldNotBeNil)
		})
	})
}
