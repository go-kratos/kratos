package dao

import (
	"context"
	"go-common/app/job/main/thumbup/model"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaohashStatsKey(t *testing.T) {
	convey.Convey("hashStatsKey", t, func(convCtx convey.C) {
		var (
			businessID = int64(33)
			originID   = int64(7788)
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			p1 := hashStatsKey(businessID, originID)
			convCtx.Convey("Then p1 should not be nil.", func(convCtx convey.C) {
				convCtx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoExpireHashStatsCache(t *testing.T) {
	convey.Convey("ExpireHashStatsCache", t, func(convCtx convey.C) {
		var (
			c          = context.Background()
			businessID = int64(33)
			originID   = int64(7788)
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			ok, err := d.ExpireHashStatsCache(c, businessID, originID)
			convCtx.Convey("Then err should be nil.ok should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(ok, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoAddHashStatsCache(t *testing.T) {
	convey.Convey("AddHashStatsCache", t, func(convCtx convey.C) {
		var (
			c          = context.Background()
			businessID = int64(33)
			originID   = int64(7788)
			stats      = &model.Stats{}
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			err := d.AddHashStatsCache(c, businessID, originID, stats)
			convCtx.Convey("Then err should be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
			})
		})
	})
}
