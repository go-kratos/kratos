package dao

import (
	"context"
	"go-common/app/service/main/thumbup/model"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaostatsKey(t *testing.T) {
	var (
		businessID = int64(1)
		messageID  = int64(1)
	)
	convey.Convey("statsKey", t, func(ctx convey.C) {
		p1 := statsKey(businessID, messageID)
		ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
			ctx.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func TestDaorecoverStatsValue(t *testing.T) {
	var (
		c = context.TODO()
		s = ""
	)
	convey.Convey("recoverStatsValue", t, func(ctx convey.C) {
		res := recoverStatsValue(c, s)
		ctx.Convey("Then res should not be nil.", func(ctx convey.C) {
			ctx.So(res, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoAddStatsCache(t *testing.T) {
	var (
		c          = context.TODO()
		businessID = int64(1)
		vs         = &model.Stats{}
	)
	convey.Convey("AddStatsCache", t, func(ctx convey.C) {
		err := d.AddStatsCache(c, businessID, vs)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoDelStatsCache(t *testing.T) {
	var (
		c          = context.TODO()
		businessID = int64(1)
		messageID  = int64(1)
	)
	convey.Convey("DelStatsCache", t, func(ctx convey.C) {
		err := d.DelStatsCache(c, businessID, messageID)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoAddStatsCacheMap(t *testing.T) {
	var (
		c          = context.TODO()
		businessID = int64(1)
		stats      map[int64]*model.Stats
	)
	convey.Convey("AddStatsCacheMap", t, func(ctx convey.C) {
		err := d.AddStatsCacheMap(c, businessID, stats)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoStatsCache(t *testing.T) {
	var (
		c          = context.TODO()
		businessID = int64(1)
		messageIDs = []int64{1}
	)
	convey.Convey("StatsCache", t, func(ctx convey.C) {
		cached, missed, err := d.StatsCache(c, businessID, messageIDs)
		ctx.Convey("Then err should be nil.cached,missed should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(missed, convey.ShouldNotBeNil)
			ctx.So(cached, convey.ShouldBeEmpty)
		})
	})
}
