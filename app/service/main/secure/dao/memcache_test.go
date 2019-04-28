package dao

import (
	"context"
	"reflect"
	"testing"

	"go-common/app/service/main/secure/model"
	"go-common/library/cache/memcache"

	"github.com/bouk/monkey"
	. "github.com/smartystreets/goconvey/convey"
)

func TestAddLocsCache(t *testing.T) {
	Convey("TestAddLocsCache", t, func() {
		err := d.AddLocsCache(context.TODO(), 3, &model.Locs{LocsCount: map[int64]int64{3: 2, 4: 3}})
		So(err, ShouldBeNil)
	})
	Convey("TestAddLocsCacheErr", t, func() {
		connGuard := monkey.PatchInstanceMethod(reflect.TypeOf(d.mc), "Get", func(_ *memcache.Pool, _ context.Context) memcache.Conn {
			return memcache.MockWith(memcache.ErrItemObject)
		})
		defer connGuard.Unpatch()
		err := d.AddLocsCache(context.TODO(), 3, &model.Locs{LocsCount: map[int64]int64{3: 2, 4: 3}})
		So(err, ShouldEqual, memcache.ErrItemObject)
	})
}

func TestLocsCache(t *testing.T) {
	Convey("TestGetLocsCache", t, func() {
		locs, err := d.LocsCache(context.TODO(), 3)
		So(err, ShouldBeNil)
		So(locs, ShouldNotBeNil)
	})
	Convey("TestGetLocsCacheGetErr", t, func() {
		connGuard := monkey.PatchInstanceMethod(reflect.TypeOf(d.mc), "Get", func(_ *memcache.Pool, _ context.Context) memcache.Conn {
			return memcache.MockWith(memcache.ErrNotFound)
		})
		defer connGuard.Unpatch()
		locs, err := d.LocsCache(context.TODO(), 3)
		So(err, ShouldBeNil)
		So(locs, ShouldBeNil)
	})

	Convey("TestGetLocsCacheScanErr", t, func() {
		connGuard := monkey.PatchInstanceMethod(reflect.TypeOf(d.mc), "Get", func(_ *memcache.Pool, _ context.Context) memcache.Conn {
			return memcache.MockWith(memcache.ErrItemObject)
		})
		defer connGuard.Unpatch()
		locs, err := d.LocsCache(context.TODO(), 3)
		So(err, ShouldEqual, memcache.ErrItemObject)
		So(locs, ShouldBeNil)
	})
}

func TestMcPing(t *testing.T) {
	Convey("TestMcPing", t, func() {
		err := d.pingMC(context.Background())
		So(err, ShouldBeNil)
	})

	Convey("TestMcPingErr", t, func() {
		connGuard := monkey.PatchInstanceMethod(reflect.TypeOf(d.mc), "Get", func(_ *memcache.Pool, _ context.Context) memcache.Conn {
			return memcache.MockWith(memcache.ErrConnClosed)
		})
		defer connGuard.Unpatch()
		err := d.pingMC(context.Background())
		So(err, ShouldNotBeNil)
		So(err, ShouldEqual, memcache.ErrConnClosed)
	})
}
