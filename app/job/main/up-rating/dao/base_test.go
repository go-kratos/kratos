package dao

import (
	"context"
	"testing"
	"time"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoGetBaseInfo(t *testing.T) {
	convey.Convey("GetBaseInfo", t, func(ctx convey.C) {
		var (
			c                = context.Background()
			month time.Month = time.June
			start            = int(0)
			end              = int(1000)
			limit            = int(1)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			d.db.Exec(c, "INSERT INTO up_level_info_06(mid) VALUES(100) ON DUPLICATE KEY UPDATE mid=VALUES(mid)")
			bs, id, err := d.GetBaseInfo(c, month, start, end, limit)
			ctx.Convey("Then err should be nil.bs,id should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(id, convey.ShouldNotBeNil)
				ctx.So(bs, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoGetBaseTotal(t *testing.T) {
	convey.Convey("GetBaseTotal", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			date  = time.Date(2018, time.June, 1, 0, 0, 0, 0, time.Local)
			id    = int64(0)
			limit = int64(100)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			d.db.Exec(c, "INSERT INTO up_level_info_06(mid) VALUES(101) ON DUPLICATE KEY UPDATE mid=VALUES(mid)")
			bs, err := d.GetBaseTotal(c, date, id, limit)
			ctx.Convey("Then err should be nil.bs should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(bs, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoBaseInfoStart(t *testing.T) {
	convey.Convey("BaseInfoStart", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			date = time.Date(2018, time.June, 1, 0, 0, 0, 0, time.Local)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			d.db.Exec(c, "INSERT INTO up_level_info_06(mid) VALUES(102) ON DUPLICATE KEY UPDATE mid=VALUES(mid)")
			start, err := d.BaseInfoStart(c, date)
			ctx.Convey("Then err should be nil.start should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(start, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoBaseInfoEnd(t *testing.T) {
	convey.Convey("BaseInfoEnd", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			date = time.Date(2018, time.June, 1, 0, 0, 0, 0, time.Local)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			d.db.Exec(c, "INSERT INTO up_level_info_06(mid) VALUES(103) ON DUPLICATE KEY UPDATE mid=VALUES(mid)")
			end, err := d.BaseInfoEnd(c, date)
			ctx.Convey("Then err should be nil.end should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(end, convey.ShouldNotBeNil)
			})
		})
	})
}
