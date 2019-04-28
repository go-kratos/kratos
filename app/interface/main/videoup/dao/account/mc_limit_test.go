package account

import (
	"context"
	"go-common/library/cache/memcache"
	"reflect"
	"testing"

	"github.com/bouk/monkey"
	"github.com/smartystreets/goconvey/convey"
)

func TestAccountlimitMidHafMin(t *testing.T) {
	convey.Convey("limitMidHafMin", t, func(ctx convey.C) {
		var (
			mid = int64(2089809)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := limitMidHafMin(mid)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestAccountlimitMidSameTitle(t *testing.T) {
	convey.Convey("limitMidSameTitle", t, func(ctx convey.C) {
		var (
			mid   = int64(2089809)
			title = "iamtitle"
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := limitMidSameTitle(mid, title)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestAccountHalfMin(t *testing.T) {
	convey.Convey("HalfMin", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(2089809)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			connGuard := monkey.PatchInstanceMethod(reflect.TypeOf(d.mc), "Get", func(_ *memcache.Pool, _ context.Context) memcache.Conn {
				return memcache.MockWith(memcache.ErrNotFound)
			})
			defer connGuard.Unpatch()
			exist, ts, err := d.HalfMin(c, mid)
			ctx.Convey("Then err should be nil.exist,ts should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(ts, convey.ShouldNotBeNil)
				ctx.So(exist, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestAccountAddHalfMin(t *testing.T) {
	convey.Convey("AddHalfMin", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(2089809)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			connGuard := monkey.PatchInstanceMethod(reflect.TypeOf(d.mc), "Get", func(_ *memcache.Pool, _ context.Context) memcache.Conn {
				return memcache.MockWith(memcache.ErrNotFound)
			})
			defer connGuard.Unpatch()
			err := d.AddHalfMin(c, mid)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestAccountDelHalfMin(t *testing.T) {
	convey.Convey("DelHalfMin", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(2089809)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			connGuard := monkey.PatchInstanceMethod(reflect.TypeOf(d.mc), "Get", func(_ *memcache.Pool, _ context.Context) memcache.Conn {
				return memcache.MockWith(memcache.ErrNotFound)
			})
			defer connGuard.Unpatch()
			err := d.DelHalfMin(c, mid)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestAccountSubmitCache(t *testing.T) {
	convey.Convey("SubmitCache", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			mid   = int64(2089809)
			title = "iamtitle"
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			connGuard := monkey.PatchInstanceMethod(reflect.TypeOf(d.mc), "Get", func(_ *memcache.Pool, _ context.Context) memcache.Conn {
				return memcache.MockWith(memcache.ErrNotFound)
			})
			defer connGuard.Unpatch()
			exist, err := d.SubmitCache(c, mid, title)
			ctx.Convey("Then err should be nil.exist should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(exist, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestAccountAddSubmitCache(t *testing.T) {
	convey.Convey("AddSubmitCache", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			mid   = int64(2089809)
			title = "iamtitle"
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			connGuard := monkey.PatchInstanceMethod(reflect.TypeOf(d.mc), "Get", func(_ *memcache.Pool, _ context.Context) memcache.Conn {
				return memcache.MockWith(memcache.ErrNotFound)
			})
			defer connGuard.Unpatch()
			err := d.AddSubmitCache(c, mid, title)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestAccountDelSubmitCache(t *testing.T) {
	convey.Convey("DelSubmitCache", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			mid   = int64(2089809)
			title = "iamtitle"
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			connGuard := monkey.PatchInstanceMethod(reflect.TypeOf(d.mc), "Get", func(_ *memcache.Pool, _ context.Context) memcache.Conn {
				return memcache.MockWith(memcache.ErrNotFound)
			})
			defer connGuard.Unpatch()
			err := d.DelSubmitCache(c, mid, title)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestAccountpingMemcache(t *testing.T) {
	convey.Convey("pingMemcache", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			err error
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			connGuard := monkey.PatchInstanceMethod(reflect.TypeOf(d.mc), "Get", func(_ *memcache.Pool, _ context.Context) memcache.Conn {
				return memcache.MockWith(memcache.ErrNotFound)
			})
			defer connGuard.Unpatch()
			err = d.pingMemcache(c)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldNotBeNil)
			})
		})
	})
}
