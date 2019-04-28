package dao

import (
	"context"
	"testing"

	"go-common/app/job/main/thumbup/model"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaostatsKey(t *testing.T) {
	convey.Convey("statsKey", t, func(convCtx convey.C) {
		var (
			businessID = int64(33)
			messageID  = int64(2233)
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			p1 := statsKey(businessID, messageID)
			convCtx.Convey("Then p1 should not be nil.", func(convCtx convey.C) {
				convCtx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoAddStatsCache(t *testing.T) {
	convey.Convey("AddStatsCache", t, func(convCtx convey.C) {
		var (
			c          = context.Background()
			businessID = int64(33)
			ms         = &model.Stats{}
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			err := d.AddStatsCache(c, businessID, ms)
			convCtx.Convey("Then err should be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
			})
		})
	})
}
