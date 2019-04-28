package income

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestIncomeGetAvDailyCharge(t *testing.T) {
	convey.Convey("GetAvDailyCharge", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			month = int(1)
			avID  = int64(1001)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			Exec(c, "INSERT INTO av_daily_charge_01(av_id, date) VALUES(1001, '2018-01-01')")
			_, err := d.GetAvDailyCharge(c, month, avID)
			ctx.Convey("Then err should be nil.avs should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})

	convey.Convey("GetAvDailyCharge month error", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			month = int(13)
			avID  = int64(1)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			_, err := d.GetAvDailyCharge(c, month, avID)
			ctx.Convey("Then err should be nil.avs should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestIncomeGetColumnCharges(t *testing.T) {
	convey.Convey("GetColumnCharges", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			aid = int64(1002)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			Exec(c, "INSERT INTO column_daily_charge(aid, date) VALUES(1002, '2018-01-01')")
			cms, err := d.GetColumnCharges(c, aid)
			ctx.Convey("Then err should be nil.cms should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(cms, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestIncomeGetBgmCharges(t *testing.T) {
	convey.Convey("GetBgmCharges", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			aid = int64(1003)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			Exec(c, "INSERT INTO bgm_daily_charge(sid, date) VALUES(1003, '2018-01-01')")
			bgms, err := d.GetBgmCharges(c, aid)
			ctx.Convey("Then err should be nil.bgms should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(bgms, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestIncomeGetArchiveChargeStatis(t *testing.T) {
	convey.Convey("GetArchiveChargeStatis", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			table = "av_charge_daily_statis"
			query = "cdate = '2018-01-01'"
			from  = int(0)
			limit = int(10)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			Exec(c, "INSERT INTO av_charge_daily_statis(avs, cdate) VALUES(10, '2018-01-01')")
			archs, err := d.GetArchiveChargeStatis(c, table, query, from, limit)
			ctx.Convey("Then err should be nil.archs should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(archs, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestIncomeGetTotalCharge(t *testing.T) {
	convey.Convey("GetTotalCharge", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			table = "av_charge_statis"
			query = "av_id = 1001"
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			Exec(c, "INSERT INTO av_charge_statis(av_id, total_income) VALUES(1001, 10)")
			total, err := d.GetTotalCharge(c, table, query)
			ctx.Convey("Then err should be nil.total should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(total, convey.ShouldNotBeNil)
			})
		})
	})
}
