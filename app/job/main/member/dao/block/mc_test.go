package block

import (
	"context"
	"reflect"
	"testing"

	"go-common/library/cache/memcache"

	"github.com/bouk/monkey"
	"github.com/smartystreets/goconvey/convey"
)

func TestBlockuserKey(t *testing.T) {
	convey.Convey("userKey", t, func(convCtx convey.C) {
		var (
			mid = int64(0)
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			key := userKey(mid)
			convCtx.Convey("Then key should not be nil.", func(convCtx convey.C) {
				convCtx.So(key, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestBlocksyncBlockTypeID(t *testing.T) {
	convey.Convey("syncBlockTypeID", t, func(convCtx convey.C) {
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			key := syncBlockTypeID()
			convCtx.Convey("Then key should not be nil.", func(convCtx convey.C) {
				convCtx.So(key, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestBlockSetSyncBlockTypeID(t *testing.T) {
	convey.Convey("SetSyncBlockTypeID", t, func(convCtx convey.C) {
		var (
			c  = context.Background()
			id = int64(0)
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			err := d.SetSyncBlockTypeID(c, id)
			convCtx.Convey("Then err should be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestBlockSyncBlockTypeID(t *testing.T) {
	convey.Convey("SyncBlockTypeID", t, func(convCtx convey.C) {
		var (
			c = context.Background()
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			id, err := d.SyncBlockTypeID(c)
			convCtx.Convey("Then err should be nil.id should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(id, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestBlockDeleteUserCache(t *testing.T) {
	convey.Convey("DeleteUserCache", t, func(convCtx convey.C) {
		var (
			c   = context.Background()
			mid = int64(0)
		)
		convCtx.Convey("connect memcache failed", func(convCtx convey.C) {
			guard := monkey.PatchInstanceMethod(reflect.TypeOf(d.mc), "Get", func(_ *memcache.Pool, _ context.Context) memcache.Conn {
				return memcache.MockWith(memcache.ErrConnClosed)
			})
			defer guard.Unpatch()
			err := d.DeleteUserCache(c, mid)
			convCtx.Convey("Then err should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldNotBeNil)
			})
		})
	})
}
