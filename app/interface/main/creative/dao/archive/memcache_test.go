package archive

import (
	"context"
	arcmdl "go-common/app/interface/main/creative/model/archive"
	"go-common/library/cache/memcache"
	"reflect"
	"testing"

	"github.com/bouk/monkey"
	"github.com/smartystreets/goconvey/convey"
)

func TestArchivekeyPorder(t *testing.T) {
	var (
		aid = int64(10110560)
	)
	convey.Convey("keyPorder", t, func(ctx convey.C) {
		p1 := keyPorder(aid)
		ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
			ctx.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func TestArchivekeyArcCM(t *testing.T) {
	var (
		aid = int64(10110560)
	)
	convey.Convey("keyArcCM", t, func(ctx convey.C) {
		p1 := keyArcCM(aid)
		ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
			ctx.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func TestArchivePOrderCache(t *testing.T) {
	var (
		c   = context.TODO()
		aid = int64(10110560)
	)
	convey.Convey("POrderCache", t, func(ctx convey.C) {
		connGuard := monkey.PatchInstanceMethod(reflect.TypeOf(d.mc), "Get", func(_ *memcache.Pool, _ context.Context) memcache.Conn {
			return memcache.MockWith(memcache.ErrNotFound)
		})
		defer connGuard.Unpatch()
		st, err := d.POrderCache(c, aid)
		ctx.Convey("Then err should be nil.st should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(st, convey.ShouldBeNil)
		})
	})
}

func TestArchiveAddPOrderCache(t *testing.T) {
	var (
		c   = context.TODO()
		aid = int64(10110560)
		st  = &arcmdl.Porder{}
	)
	convey.Convey("AddPOrderCache", t, func(ctx convey.C) {
		connGuard := monkey.PatchInstanceMethod(reflect.TypeOf(d.mc), "Get", func(_ *memcache.Pool, _ context.Context) memcache.Conn {
			return memcache.MockWith(memcache.ErrNotFound)
		})
		defer connGuard.Unpatch()
		err := d.AddPOrderCache(c, aid, st)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldNotBeNil)
		})
	})
}

func TestArchiveArcCMCache(t *testing.T) {
	var (
		c   = context.TODO()
		aid = int64(10110560)
	)
	convey.Convey("ArcCMCache", t, func(ctx convey.C) {
		connGuard := monkey.PatchInstanceMethod(reflect.TypeOf(d.mc), "Get", func(_ *memcache.Pool, _ context.Context) memcache.Conn {
			return memcache.MockWith(memcache.ErrNotFound)
		})
		defer connGuard.Unpatch()
		st, err := d.ArcCMCache(c, aid)
		ctx.Convey("Then err should be nil.st should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(st, convey.ShouldBeNil)
		})
	})
}

func TestArchiveAddArcCMCache(t *testing.T) {
	var (
		c   = context.TODO()
		aid = int64(10110560)
		st  = &arcmdl.Commercial{}
	)
	convey.Convey("AddArcCMCache", t, func(ctx convey.C) {
		connGuard := monkey.PatchInstanceMethod(reflect.TypeOf(d.mc), "Get", func(_ *memcache.Pool, _ context.Context) memcache.Conn {
			return memcache.MockWith(memcache.ErrNotFound)
		})
		defer connGuard.Unpatch()
		err := d.AddArcCMCache(c, aid, st)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldNotBeNil)
		})
	})
}
