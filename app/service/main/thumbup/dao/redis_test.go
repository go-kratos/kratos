package dao

import (
	"context"
	"testing"

	"go-common/app/service/main/thumbup/model"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaohashStatsKey(t *testing.T) {
	var (
		businessID = int64(1)
		originID   = int64(1)
	)
	convey.Convey("hashStatsKey", t, func(ctx convey.C) {
		p1 := hashStatsKey(businessID, originID)
		ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
			ctx.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoExpireHashStatsCache(t *testing.T) {
	var (
		c          = context.TODO()
		businessID = int64(1)
		originID   = int64(1)
	)
	convey.Convey("ExpireHashStatsCache", t, func(ctx convey.C) {
		ok, err := d.ExpireHashStatsCache(c, businessID, originID)
		ctx.Convey("Then err should be nil.ok should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(ok, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoDelHashStatsCache(t *testing.T) {
	var (
		c          = context.TODO()
		businessID = int64(1)
		originID   = int64(1)
	)
	convey.Convey("DelHashStatsCache", t, func(ctx convey.C) {
		err := d.DelHashStatsCache(c, businessID, originID)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoHashStatsCache(t *testing.T) {
	var (
		c          = context.TODO()
		businessID = int64(1)
		originID   = int64(1)
		messageIDs = []int64{1}
	)
	convey.Convey("HashStatsCache", t, func(ctx convey.C) {
		res, err := d.HashStatsCache(c, businessID, originID, messageIDs)
		ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(res, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoAddHashStatsCache(t *testing.T) {
	var (
		c          = context.TODO()
		businessID = int64(1)
		originID   = int64(1)
		stats      = &model.Stats{}
	)
	convey.Convey("AddHashStatsCache", t, func(ctx convey.C) {
		err := d.AddHashStatsCache(c, businessID, originID, stats)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoAddHashStatsCacheMap(t *testing.T) {
	var (
		c          = context.TODO()
		businessID = int64(1)
		originID   = int64(1)
		stats      map[int64]*model.Stats
	)
	convey.Convey("AddHashStatsCacheMap", t, func(ctx convey.C) {
		err := d.AddHashStatsCacheMap(c, businessID, originID, stats)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}
