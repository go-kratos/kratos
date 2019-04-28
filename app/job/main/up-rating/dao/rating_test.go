package dao

import (
	"context"
	"testing"
	"time"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoDelRatings(t *testing.T) {
	convey.Convey("DelRatings", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			date  = time.Date(2018, time.June, 1, 0, 0, 0, 0, time.Local)
			limit = int(10)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			d.db.Exec(c, "INSERT INTO up_rating_06(mid,cdate) VALUES(1,'2018-06-01') ON DUPLICATE KEY UPDATE mid=VALUES(mid)")
			rows, err := d.DelRatings(c, date, limit)
			ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rows, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoRatingStart(t *testing.T) {
	convey.Convey("RatingStart", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			date = time.Date(2018, time.June, 1, 0, 0, 0, 0, time.Local)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			d.db.Exec(c, "INSERT INTO up_rating_06(mid,cdate) VALUES(1,'2018-06-01') ON DUPLICATE KEY UPDATE mid=VALUES(mid)")
			start, err := d.RatingStart(c, date)
			ctx.Convey("Then err should be nil.start should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(start, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoRatingEnd(t *testing.T) {
	convey.Convey("RatingEnd", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			date = time.Date(2018, time.June, 1, 0, 0, 0, 0, time.Local)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			d.db.Exec(c, "INSERT INTO up_rating_06(mid,cdate) VALUES(1,'2018-06-01') ON DUPLICATE KEY UPDATE mid=VALUES(mid)")
			end, err := d.RatingEnd(c, date)
			ctx.Convey("Then err should be nil.end should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(end, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoRatingCount(t *testing.T) {
	convey.Convey("RatingCount", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			date = time.Date(2018, time.June, 1, 0, 0, 0, 0, time.Local)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			d.db.Exec(c, "INSERT INTO up_rating_06(mid,cdate) VALUES(1,'2018-06-01') ON DUPLICATE KEY UPDATE mid=VALUES(mid)")
			count, err := d.RatingCount(c, date)
			ctx.Convey("Then err should be nil.count should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(count, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoGetRatings(t *testing.T) {
	convey.Convey("GetRatings", t, func(ctx convey.C) {
		var (
			c      = context.Background()
			date   = time.Date(2018, time.June, 1, 0, 0, 0, 0, time.Local)
			offset = int(0)
			limit  = int(100)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			d.db.Exec(c, "INSERT INTO up_rating_06(mid,cdate) VALUES(1,'2018-06-01') ON DUPLICATE KEY UPDATE mid=VALUES(mid)")
			rs, last, err := d.GetRatings(c, date, offset, limit)
			ctx.Convey("Then err should be nil.rs,last should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(last, convey.ShouldNotBeNil)
				ctx.So(rs, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoGetRatingsFast(t *testing.T) {
	convey.Convey("GetRatingsFast", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			date  = time.Date(2018, time.June, 1, 0, 0, 0, 0, time.Local)
			start = int(0)
			end   = int(10000)
			limit = int(10)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			d.db.Exec(c, "INSERT INTO up_rating_06(mid,cdate) VALUES(1,'2018-06-01') ON DUPLICATE KEY UPDATE mid=VALUES(mid)")
			rs, id, err := d.GetRatingsFast(c, date, start, end, limit)
			ctx.Convey("Then err should be nil.rs,id should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(id, convey.ShouldNotBeNil)
				ctx.So(rs, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoInsertRatingStat(t *testing.T) {
	convey.Convey("InsertRatingStat", t, func(ctx convey.C) {
		var (
			c                 = context.Background()
			month  time.Month = time.June
			values            = "(1,2,'2018-06-01',100,100,100,600,600,300)"
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			_, err := d.InsertRatingStat(c, month, values)
			ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}
