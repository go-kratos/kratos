package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoCacheUserCoin(t *testing.T) {
	var (
		c  = context.TODO()
		id = int64(1)
	)
	convey.Convey("CacheUserCoin", t, func(ctx convey.C) {
		res, err := d.CacheUserCoin(c, id)
		ctx.Convey("Error should be nil", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
		ctx.Convey("res should not be nil", func(ctx convey.C) {
			ctx.So(res, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoAddCacheUserCoin(t *testing.T) {
	var (
		c   = context.TODO()
		id  = int64(10)
		val = float64(10)
	)
	convey.Convey("AddCacheUserCoin", t, func(ctx convey.C) {
		err := d.AddCacheUserCoin(c, id, val)
		ctx.Convey("Error should be nil", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoCacheItemCoin(t *testing.T) {
	var (
		c  = context.TODO()
		id = int64(1)
		tp = int64(60)
	)
	convey.Convey("CacheItemCoin", t, func(ctx convey.C) {
		res, err := d.CacheItemCoin(c, id, tp)
		ctx.Convey("Error should be nil", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
		ctx.Convey("res should not be nil", func(ctx convey.C) {
			ctx.So(res, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoAddCacheItemCoin(t *testing.T) {
	var (
		c   = context.TODO()
		id  = int64(1)
		val = int64(10)
		tp  = int64(60)
	)
	convey.Convey("AddCacheItemCoin", t, func(ctx convey.C) {
		err := d.AddCacheItemCoin(c, id, val, tp)
		ctx.Convey("Error should be nil", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoExp(t *testing.T) {
	var (
		c  = context.TODO()
		id = int64(1)
	)
	convey.Convey("Exp", t, func(ctx convey.C) {
		res, err := d.Exp(c, id)
		ctx.Convey("Error should be nil", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
		ctx.Convey("res should not be nil", func(ctx convey.C) {
			ctx.So(res, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoSetTodayExpCache(t *testing.T) {
	var (
		c   = context.TODO()
		id  = int64(1)
		val = int64(10)
	)
	convey.Convey("SetTodayExpCache", t, func(ctx convey.C) {
		err := d.SetTodayExpCache(c, id, val)
		ctx.Convey("Error should be nil", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}
