package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoorderKey(t *testing.T) {
	convey.Convey("orderKey", t, func(ctx convey.C) {
		var (
			id = ""
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			p1 := orderKey(id)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoassetKey(t *testing.T) {
	convey.Convey("assetKey", t, func(ctx convey.C) {
		var (
			oid      = int64(0)
			otype    = ""
			currency = ""
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			p1 := assetKey(oid, otype, currency)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaotaskKey(t *testing.T) {
	convey.Convey("taskKey", t, func(ctx convey.C) {
		var (
			task = ""
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			p1 := taskKey(task)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoDelCacheOrderUser(t *testing.T) {
	convey.Convey("DelCacheOrderUser", t, func(ctx convey.C) {
		var (
			c  = context.Background()
			id = ""
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := d.DelCacheOrderUser(c, id)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoDelCacheAsset(t *testing.T) {
	convey.Convey("DelCacheAsset", t, func(ctx convey.C) {
		var (
			c        = context.Background()
			oid      = int64(0)
			otype    = ""
			currency = ""
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := d.DelCacheAsset(c, oid, otype, currency)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoAddCacheTask(t *testing.T) {
	convey.Convey("AddCacheTask", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			task = ""
			ttl  = int32(60)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			ok, err := d.AddCacheTask(c, task, ttl)
			ctx.Convey("Then err should be nil.ok should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(ok, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoDelCacheTask(t *testing.T) {
	convey.Convey("DelCacheTask", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			task = ""
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := d.DelCacheTask(c, task)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}
func TestDaoDelCacheUserSetting(t *testing.T) {
	convey.Convey("DelCacheUserSetting", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(46333)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := d.DelCacheUserSetting(c, mid)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}
