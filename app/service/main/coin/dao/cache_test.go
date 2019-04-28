package dao

import (
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaocountKey(t *testing.T) {
	var (
		mid = int64(123)
	)
	convey.Convey("countKey", t, func(ctx convey.C) {
		p1 := countKey(mid)
		ctx.Convey("p1 should not be nil", func(ctx convey.C) {
			ctx.So(p1, convey.ShouldNotBeEmpty)
		})
	})
}

func TestDaocacheSFUserCoin(t *testing.T) {
	var (
		mid = int64(123)
	)
	convey.Convey("cacheSFUserCoin", t, func(ctx convey.C) {
		p1 := d.cacheSFUserCoin(mid)
		ctx.Convey("p1 should not be nil", func(ctx convey.C) {
			ctx.So(p1, convey.ShouldNotBeEmpty)
		})
	})
}

func TestDaoitemCoinKey(t *testing.T) {
	var (
		aid = int64(1)
		tp  = int64(60)
	)
	convey.Convey("itemCoinKey", t, func(ctx convey.C) {
		p1 := itemCoinKey(aid, tp)
		ctx.Convey("p1 should not be nil", func(ctx convey.C) {
			ctx.So(p1, convey.ShouldNotBeEmpty)
		})
	})
}

func TestDaocacheSFItemCoin(t *testing.T) {
	var (
		aid = int64(1)
		tp  = int64(60)
	)
	convey.Convey("cacheSFItemCoin", t, func(ctx convey.C) {
		p1 := d.cacheSFItemCoin(aid, tp)
		ctx.Convey("p1 should not be nil", func(ctx convey.C) {
			ctx.So(p1, convey.ShouldNotBeEmpty)
		})
	})
}

func TestDaoexpKey(t *testing.T) {
	var (
		mid = int64(123)
	)
	convey.Convey("expKey", t, func(ctx convey.C) {
		key := expKey(mid)
		ctx.Convey("key should not be nil", func(ctx convey.C) {
			ctx.So(key, convey.ShouldNotBeEmpty)
		})
	})
}
