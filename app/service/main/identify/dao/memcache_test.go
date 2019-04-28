package dao

import (
	"context"
	"testing"

	"go-common/app/service/main/identify/model"
	"go-common/library/cache/memcache"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoSetAccessCache(t *testing.T) {
	var (
		c   = context.Background()
		key = "123"
		res = &model.IdentifyInfo{Mid: 1, Csrf: "csrf", Expires: 1577811661}
	)
	convey.Convey("SetAccessCache", t, func(ctx convey.C) {
		d.SetAccessCache(c, key, res)
		ctx.Convey("No return values", func(ctx convey.C) {
		})
	})
}

func TestDaoAccessCache(t *testing.T) {
	var (
		c       = context.Background()
		hitKey  = "123"
		missKey = "miss"
	)
	convey.Convey("AccessCache", t, func(ctx convey.C) {
		res, err := d.AccessCache(c, hitKey)
		ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(res, convey.ShouldNotBeNil)
			ctx.So(res.Mid, convey.ShouldEqual, 1)
		})
	})
	convey.Convey("AccessCache miss", t, func(ctx convey.C) {
		res, err := d.AccessCache(c, missKey)
		ctx.Convey("Then err and res should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(res, convey.ShouldBeNil)
		})
	})
}

func TestDaoDelCache(t *testing.T) {
	var (
		c       = context.Background()
		hitKey  = "123"
		missKey = "miss"
	)
	convey.Convey("DelCache", t, func(ctx convey.C) {
		err := d.DelCache(c, hitKey)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
	convey.Convey("DelCache miss", t, func(ctx convey.C) {
		err := d.DelCache(c, missKey)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaocacheKey(t *testing.T) {
	var (
		key = "1"
	)
	convey.Convey("cacheKey", t, func(ctx convey.C) {
		p1 := cacheKey(key)
		ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
			ctx.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func TestDaologinCacheKey(t *testing.T) {
	var (
		mid = int64(1)
		ip  = "127.0.0.1"
	)
	convey.Convey("loginCacheKey", t, func(ctx convey.C) {
		p1 := loginCacheKey(mid, ip)
		ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
			ctx.So(p1, convey.ShouldEqual, "l1127.0.0.1")
		})
	})
}

func TestDao_SetLoginCache(t *testing.T) {
	var (
		c       = context.Background()
		mid     = int64(1)
		ip      = "127.0.0.1"
		expires = int32(1577811661)
	)
	convey.Convey("SetLoginCache", t, func(ctx convey.C) {
		err := d.SetLoginCache(c, mid, ip, expires)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeIn, nil, memcache.ErrNotStored)
		})
	})
}

func TestDao_ExistMIDAndIP(t *testing.T) {
	var (
		c      = context.Background()
		mid    = int64(1)
		hitIP  = "127.0.0.1"
		missIP = "127.0.0.2"
	)
	convey.Convey("IsExistMID", t, func(ctx convey.C) {
		ok, err := d.ExistMIDAndIP(c, mid, hitIP)
		ctx.Convey("Then err should be nil.ok should be true.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(ok, convey.ShouldEqual, true)
		})
	})
	convey.Convey("IsExistMID miss", t, func(ctx convey.C) {
		ok, err := d.ExistMIDAndIP(c, mid, missIP)
		ctx.Convey("Then err should be nil and ok should be false.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(ok, convey.ShouldEqual, false)
		})
	})
}
