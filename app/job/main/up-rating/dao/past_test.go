package dao

import (
	"context"
	"testing"
	"time"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoDelPastRecord(t *testing.T) {
	convey.Convey("DelPastRecord", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			date = time.Date(2018, time.June, 1, 0, 0, 0, 0, time.Local)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			d.db.Exec(c, "INSERT INTO past_rating_record(times,date) VALUES(1, '2018-06-01') ON DUPLICATE KEY UPDATE mid=VALUES(mid)")
			rows, err := d.DelPastRecord(c, date)
			ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rows, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoGetPastRecord(t *testing.T) {
	convey.Convey("GetPastRecord", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			cdate = "2018-06-01"
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			d.db.Exec(c, "INSERT INTO past_rating_record(times,date) VALUES(1, '2018-06-01') ON DUPLICATE KEY UPDATE mid=VALUES(mid)")
			times, err := d.GetPastRecord(c, cdate)
			ctx.Convey("Then err should be nil.times should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(times, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoInsertPastRecord(t *testing.T) {
	convey.Convey("InsertPastRecord", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			times = int(1)
			cdate = "2018-06-01"
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			rows, err := d.InsertPastRecord(c, times, cdate)
			ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rows, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoInsertPastScoreStat(t *testing.T) {
	convey.Convey("InsertPastScoreStat", t, func(ctx convey.C) {
		var (
			c      = context.Background()
			values = "(1, 100, 100, 100)"
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			rows, err := d.InsertPastScoreStat(c, values)
			ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rows, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoGetPasts(t *testing.T) {
	convey.Convey("GetPasts", t, func(ctx convey.C) {
		var (
			c      = context.Background()
			offset = int64(0)
			limit  = int64(1000)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			d.db.Exec(c, "INSERT INTO past_score_statistics(mid) VALUES(1) ON DUPLICATE KEY UPDATE mid=VALUES(mid)")
			past, last, err := d.GetPasts(c, offset, limit)
			ctx.Convey("Then err should be nil.past,last should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(last, convey.ShouldNotBeNil)
				ctx.So(past, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoDelPastStat(t *testing.T) {
	convey.Convey("DelPastStat", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			limit = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			d.db.Exec(c, "INSERT INTO past_score_statistics(mid) VALUES(1) ON DUPLICATE KEY UPDATE mid=VALUES(mid)")
			rows, err := d.DelPastStat(c, limit)
			ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rows, convey.ShouldNotBeNil)
			})
		})
	})
}
