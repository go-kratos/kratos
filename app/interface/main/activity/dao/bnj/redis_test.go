package bnj

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestBnjresetKey(t *testing.T) {
	convey.Convey("resetKey", t, func(ctx convey.C) {
		var (
			mid = int64(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			p1 := resetKey(mid)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestBnjrewardKey(t *testing.T) {
	convey.Convey("rewardKey", t, func(ctx convey.C) {
		var (
			mid   = int64(0)
			subID = int64(0)
			step  = int(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			p1 := rewardKey(mid, subID, step)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestBnjCacheResetCD(t *testing.T) {
	convey.Convey("CacheResetCD", t, func(ctx convey.C) {
		var (
			c      = context.Background()
			mid    = int64(0)
			expire = int32(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			p1, err := d.CacheResetCD(c, mid, expire)
			ctx.Convey("Then err should be nil.p1 should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestBnjTTLResetCD(t *testing.T) {
	convey.Convey("CacheTTLResetCD", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			p1, err := d.TTLResetCD(c, mid)
			ctx.Convey("Then err should be nil.p1 should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestBnjCacheHasReward(t *testing.T) {
	convey.Convey("CacheHasReward", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			mid   = int64(0)
			subID = int64(0)
			step  = int(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			p1, err := d.CacheHasReward(c, mid, subID, step)
			ctx.Convey("Then err should be nil.p1 should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestBnjDelCacheHasReward(t *testing.T) {
	convey.Convey("DelCacheHasReward", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			mid   = int64(0)
			subID = int64(0)
			step  = int(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := d.DelCacheHasReward(c, mid, subID, step)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestBnjHasReward(t *testing.T) {
	convey.Convey("HasReward", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			mid   = int64(27515251)
			subID = int64(0)
			step  = int(1)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			res, err := d.HasReward(c, mid, subID, step)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
				ctx.Printf("%+v", res)
			})
		})
	})
}

func TestBnjsetNXLockCache(t *testing.T) {
	convey.Convey("setNXLockCache", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			key   = ""
			times = int32(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			res, err := d.setNXLockCache(c, key, times)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestBnjdelNXLockCache(t *testing.T) {
	convey.Convey("delNXLockCache", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			key = ""
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := d.delNXLockCache(c, key)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}
