package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoAvDailyIncCharge(t *testing.T) {
	convey.Convey("AvDailyIncCharge", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			avID = int64(10)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			Exec(c, "INSERT INTO av_daily_charge_06(av_id,inc_charge,date,upload_time) VALUES(10, 100, '2018-06-24', '2018-06-24') ON DUPLICATE KEY UPDATE av_id=VALUES(av_id)")
			incCharge, tagID, err := d.AvDailyIncCharge(c, avID)
			ctx.Convey("Then err should be nil.incCharge,tagID should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(tagID, convey.ShouldNotBeNil)
				ctx.So(incCharge, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoAvChargeRatio(t *testing.T) {
	convey.Convey("AvChargeRatio", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			id    = int64(10)
			limit = int64(1000)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			m, last, err := d.AvChargeRatio(c, id, limit)
			ctx.Convey("Then err should be nil.m,last should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(last, convey.ShouldNotBeNil)
				ctx.So(m, convey.ShouldNotBeNil)
			})
		})
	})
}
