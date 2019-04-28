package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaohashField(t *testing.T) {
	var (
		aid = int64(1)
		tp  = int64(1)
	)
	convey.Convey("hashField", t, func(ctx convey.C) {
		p1 := hashField(aid, tp)
		ctx.Convey("p1 should not be nil", func(ctx convey.C) {
			ctx.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoaddKey2(t *testing.T) {
	var (
		mid = int64(123)
	)
	convey.Convey("addKey2", t, func(ctx convey.C) {
		key := addKey2(mid)
		ctx.Convey("key should not be nil", func(ctx convey.C) {
			ctx.So(key, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoCoinsAddedCache(t *testing.T) {
	var (
		c   = context.TODO()
		mid = int64(123)
		aid = int64(1)
		tp  = int64(60)
	)
	convey.Convey("CoinsAddedCache", t, func(ctx convey.C) {
		added, err := d.CoinsAddedCache(c, mid, aid, tp)
		ctx.Convey("Error should be nil", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
		ctx.Convey("added should not be nil", func(ctx convey.C) {
			ctx.So(added, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoSetCoinAddedCache(t *testing.T) {
	var (
		c     = context.TODO()
		mid   = int64(123)
		aid   = int64(1)
		tp    = int64(60)
		count = int64(20)
	)
	convey.Convey("SetCoinAddedCache", t, func(ctx convey.C) {
		err := d.SetCoinAddedCache(c, mid, aid, tp, count)
		ctx.Convey("Error should be nil", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoSetCoinAddedsCache(t *testing.T) {
	var (
		c      = context.TODO()
		mid    = int64(123)
		counts = map[int64]int64{10: 20, 20: 30, 30: 40}
	)
	convey.Convey("SetCoinAddedsCache", t, func(ctx convey.C) {
		err := d.SetCoinAddedsCache(c, mid, counts)
		ctx.Convey("Error should be nil", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoIncrCoinAddedCache(t *testing.T) {
	var (
		c     = context.TODO()
		mid   = int64(123)
		aid   = int64(1)
		tp    = int64(60)
		count = int64(20)
	)
	convey.Convey("IncrCoinAddedCache", t, func(ctx convey.C) {
		err := d.IncrCoinAddedCache(c, mid, aid, tp, count)
		ctx.Convey("Error should be nil", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoExpireCoinAdded(t *testing.T) {
	var (
		c   = context.TODO()
		mid = int64(123)
	)
	convey.Convey("ExpireCoinAdded", t, func(ctx convey.C) {
		ok, err := d.ExpireCoinAdded(c, mid)
		ctx.Convey("Error should be nil", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
		ctx.Convey("ok should not be nil", func(ctx convey.C) {
			ctx.So(ok, convey.ShouldNotBeNil)
		})
	})
}
