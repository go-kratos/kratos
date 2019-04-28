package charge

import (
	"context"
	"testing"
	"time"

	"github.com/smartystreets/goconvey/convey"
)

func TestChargeAvDailyCharge(t *testing.T) {
	convey.Convey("AvDailyCharge", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			date  = time.Date(2018, 6, 24, 0, 0, 0, 0, time.Local)
			id    = int64(0)
			limit = int(100)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			d.db.Exec(c, "INSERT INTO av_daily_charge_06(av_id,mid,date,inc_charge) VALUES(1,2,'2018-06-24',100) ON DUPLICATE KEY UPDATE inc_charge=VALUES(inc_charge)")
			data, err := d.AvDailyCharge(c, date, id, limit)
			ctx.Convey("Then err should be nil.data should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(data, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestChargeAvWeeklyCharge(t *testing.T) {
	convey.Convey("AvWeeklyCharge", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			date  = time.Date(2018, 6, 24, 0, 0, 0, 0, time.Local)
			id    = int64(0)
			limit = int(100)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			d.db.Exec(c, "INSERT INTO av_weekly_charge(av_id,mid,date) VALUES(1,2,'2018-06-24') ON DUPLICATE KEY UPDATE mid=VALUES(mid)")
			data, err := d.AvWeeklyCharge(c, date, id, limit)
			ctx.Convey("Then err should be nil.data should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(data, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestChargeAvMonthlyCharge(t *testing.T) {
	convey.Convey("AvMonthlyCharge", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			date  = time.Date(2018, 6, 24, 0, 0, 0, 0, time.Local)
			id    = int64(0)
			limit = int(100)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			d.db.Exec(c, "INSERT INTO av_monthly_charge(av_id,mid,date) VALUES(1,2,'2018-06-24') ON DUPLICATE KEY UPDATE mid=VALUES(mid)")
			data, err := d.AvMonthlyCharge(c, date, id, limit)
			ctx.Convey("Then err should be nil.data should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(data, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestChargeGetAvCharge(t *testing.T) {
	convey.Convey("GetAvCharge", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			sql   = "SELECT id,av_id,mid,tag_id,is_original,upload_time,danmaku_count,comment_count,collect_count,coin_count,share_count,elec_pay_count,total_play_count,web_play_count,app_play_count,h5_play_count,lv_unknown,lv_0,lv_1,lv_2,lv_3,lv_4,lv_5,lv_6,v_score,inc_charge,total_charge,date,is_deleted,ctime,mtime FROM av_monthly_charge WHERE id > ? AND date = ? ORDER BY id LIMIT ?"
			date  = "2018-06-24"
			id    = int64(0)
			limit = int(100)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			d.db.Exec(c, "INSERT INTO av_daily_charge_06(av_id,mid,date) VALUES(1,2,'2018-06-24') ON DUPLICATE KEY UPDATE mid=VALUES(mid)")
			data, err := d.GetAvCharge(c, sql, date, id, limit)
			ctx.Convey("Then err should be nil.data should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(data, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestChargeInsertAvChargeTable(t *testing.T) {
	convey.Convey("InsertAvChargeTable", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			vals  = "(1,2,3,1,5,2,2,2,2,2,2,2,2,2,1,1,1,1,1,1,1,1,1,1,1,'2018-06-24', '2018-06-24')"
			table = "av_weekly_charge"
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			rows, err := d.InsertAvChargeTable(c, vals, table)
			ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rows, convey.ShouldNotBeNil)
			})
		})
	})
}
