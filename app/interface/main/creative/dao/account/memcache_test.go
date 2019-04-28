package account

import (
	"context"
	accmdl "go-common/app/interface/main/creative/model/account"
	"go-common/library/cache/memcache"
	"reflect"
	"testing"

	"github.com/bouk/monkey"
	"github.com/smartystreets/goconvey/convey"
)

func TestAccountlimitMidHafMin(t *testing.T) {
	var (
		mid = int64(2089809)
	)
	convey.Convey("limitMidHafMin", t, func(ctx convey.C) {
		p1 := limitMidHafMin(mid)
		ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
			ctx.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func TestAccountkeyUpInfo(t *testing.T) {
	var (
		mid = int64(2089809)
	)
	convey.Convey("keyUpInfo", t, func(ctx convey.C) {
		p1 := keyUpInfo(mid)
		ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
			ctx.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func TestAccountHalfMin(t *testing.T) {
	var (
		c   = context.TODO()
		mid = int64(2089809)
	)
	convey.Convey("HalfMin", t, func(ctx convey.C) {
		exist, ts, err := d.HalfMin(c, mid)
		ctx.Convey("Then err should be nil.exist,ts should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(ts, convey.ShouldNotBeNil)
			ctx.So(exist, convey.ShouldNotBeNil)
		})
	})
}

func TestAccountUpInfoCache(t *testing.T) {
	var (
		c   = context.TODO()
		mid = int64(2089809)
	)
	convey.Convey("UpInfoCache", t, func(ctx convey.C) {
		connGuard := monkey.PatchInstanceMethod(reflect.TypeOf(d.mc), "Get", func(_ *memcache.Pool, _ context.Context) memcache.Conn {
			return memcache.MockWith(memcache.ErrNotFound)
		})
		defer connGuard.Unpatch()
		st, err := d.UpInfoCache(c, mid)
		ctx.Convey("Then err should be nil.st should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(st, convey.ShouldBeNil)
		})
	})
}

func TestAccountAddUpInfoCache(t *testing.T) {
	var (
		c   = context.TODO()
		mid = int64(2089809)
		st  = &accmdl.UpInfo{}
	)
	convey.Convey("AddUpInfoCache", t, func(ctx convey.C) {
		connGuard := monkey.PatchInstanceMethod(reflect.TypeOf(d.mc), "Get", func(_ *memcache.Pool, _ context.Context) memcache.Conn {
			return memcache.MockWith(nil)
		})
		defer connGuard.Unpatch()
		err := d.AddUpInfoCache(c, mid, st)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}
