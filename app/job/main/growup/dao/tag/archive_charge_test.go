package tag

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestTagGetAvDailyCharge(t *testing.T) {
	convey.Convey("GetAvDailyCharge", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			month = "01"
			query = "id > 0"
			from  = int(0)
			limit = int(10)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			Exec(c, "INSERT INTO av_daily_charge_01(av_id, inc_charge) VALUES(111, 100)")
			avs, err := d.GetAvDailyCharge(c, month, query, from, limit)
			ctx.Convey("Then err should be nil.avs should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(avs, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestTagGetCmDailyCharge(t *testing.T) {
	convey.Convey("GetCmDailyCharge", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			query = "id > 0"
			from  = int(0)
			limit = int(10)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			Exec(c, "INSERT INTO column_daily_charge(aid, inc_charge) VALUES(111, 100)")
			cms, err := d.GetCmDailyCharge(c, query, from, limit)
			ctx.Convey("Then err should be nil.cms should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(cms, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestTagGetBgm(t *testing.T) {
	convey.Convey("GetBgm", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			id    = int64(0)
			limit = int64(10)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			Exec(c, "INSERT INTO background_music(sid, mid) VALUES(111, 1001)")
			bs, last, err := d.GetBgm(c, id, limit)
			ctx.Convey("Then err should be nil.bs,last should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(last, convey.ShouldNotBeNil)
				ctx.So(bs, convey.ShouldNotBeNil)
			})
		})
	})
}
