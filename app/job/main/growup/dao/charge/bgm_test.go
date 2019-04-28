package charge

import (
	"context"
	"testing"
	"time"

	"github.com/smartystreets/goconvey/convey"
)

func TestChargeGetBgm(t *testing.T) {
	convey.Convey("GetBgm", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			id    = int64(0)
			limit = int64(10)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			d.db.Exec(c, "INSERT INTO backgroud_music(aid,sid) VALUES(1,2)")
			bs, last, err := d.GetBgm(c, id, limit)
			ctx.Convey("Then err should be nil.bs,last should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(last, convey.ShouldNotBeNil)
				ctx.So(bs, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestChargeBgmCharge(t *testing.T) {
	convey.Convey("BgmCharge", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			date  = time.Date(2018, 6, 24, 0, 0, 0, 0, time.Local)
			id    = int64(0)
			limit = int(100)
			table = "bgm_daily_charge"
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			d.db.Exec(c, "INSERT INTO bgm_daily_charge(aid,sid,inc_charge) VALUES(1,2,3)")
			bgms, err := d.BgmCharge(c, date, id, limit, table)
			ctx.Convey("Then err should be nil.bgms should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(bgms, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestChargeBgmStatis(t *testing.T) {
	convey.Convey("BgmStatis", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			id    = int64(0)
			limit = int(100)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			d.db.Exec(c, "INSERT INTO bgm_charge_statis(aid,sid,total_charge) VALUES(1,2,3)")
			bgms, err := d.BgmStatis(c, id, limit)
			ctx.Convey("Then err should be nil.bgms should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(bgms, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestChargeInsertBgmChargeTable(t *testing.T) {
	convey.Convey("InsertBgmChargeTable", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			vals  = "(1,2,3,4,'test',100, '2018-06-24','2018-06-24')"
			table = "bgm_weekly_charge"
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			rows, err := d.InsertBgmChargeTable(c, vals, table)
			ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rows, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestChargeInsertBgmStatisBatch(t *testing.T) {
	convey.Convey("InsertBgmStatisBatch", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			vals = "(1,2,3,4,'test',100,'2018-06-24')"
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			rows, err := d.InsertBgmStatisBatch(c, vals)
			ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rows, convey.ShouldNotBeNil)
			})
		})
	})
}
