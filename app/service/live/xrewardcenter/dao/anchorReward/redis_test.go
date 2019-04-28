package anchorReward

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestAnchorRewardlockKey(t *testing.T) {
	convey.Convey("lockKey", t, func(ctx convey.C) {
		var (
			key = ""
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			p1 := lockKey(key)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestAnchorRewardexpireCountKey(t *testing.T) {
	convey.Convey("expireCountKey", t, func(ctx convey.C) {
		var (
			key = ""
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			p1 := expireCountKey(key)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestAnchorRewardDelLockCache(t *testing.T) {
	convey.Convey("DelLockCache", t, func(ctx convey.C) {
		var (
			c = context.Background()
			k = ""
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := d.DelLockCache(c, k)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestAnchorRewardGetExpireCountCache(t *testing.T) {
	convey.Convey("GetExpireCountCache", t, func(ctx convey.C) {
		var (
			c = context.Background()
			k = ""
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			count, err := d.GetExpireCountCache(c, k)
			ctx.Convey("Then err should be nil.count should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(count, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestAnchorRewardAddExpireCountCache(t *testing.T) {
	convey.Convey("AddExpireCountCache", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			k     = ""
			times = int64(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := d.AddExpireCountCache(c, k, times)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestAnchorRewardClearExpireCountCache(t *testing.T) {
	convey.Convey("ClearExpireCountCache", t, func(ctx convey.C) {
		var (
			c = context.Background()
			k = ""
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := d.ClearExpireCountCache(c, k)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestAnchorRewardSetNxLock(t *testing.T) {
	convey.Convey("SetNxLock", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			k     = ""
			times = int64(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			res, err := d.SetNxLock(c, k, times)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}
