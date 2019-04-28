package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaokeyWaitBlock(t *testing.T) {
	convey.Convey("keyWaitBlock", t, func(ctx convey.C) {
		var (
			batchNo = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := keyWaitBlock(batchNo)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaolockKey(t *testing.T) {
	convey.Convey("lockKey", t, func(ctx convey.C) {
		var (
			key = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := lockKey(key)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoAddBlockCache(t *testing.T) {
	convey.Convey("AddBlockCache", t, func(ctx convey.C) {
		var (
			c       = context.Background()
			mid     = int64(0)
			score   = int8(0)
			blockNo = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.AddBlockCache(c, mid, score, blockNo)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoBlockMidCache(t *testing.T) {
	convey.Convey("BlockMidCache", t, func(ctx convey.C) {
		var (
			c       = context.Background()
			batchNo = int64(0)
			num     = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			_, err := d.BlockMidCache(c, batchNo, num)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoSetNXLockCache(t *testing.T) {
	convey.Convey("SetNXLockCache", t, func(ctx convey.C) {
		var (
			c = context.Background()
			k = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.SetNXLockCache(c, k)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoPingRedis(t *testing.T) {
	convey.Convey("PingRedis", t, func(ctx convey.C) {
		var (
			c = context.Background()
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.PingRedis(c)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}
