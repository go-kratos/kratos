package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoAscTrendCount(t *testing.T) {
	convey.Convey("AscTrendCount", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			date  = "2018-06-01"
			query = ""
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			d.db.Exec(c, "INSERT INTO up_rating_trend_asc(mid) VALUES(1001) ON DUPLICATE KEY UPDATE mid=VALUES(mid)")
			count, err := d.AscTrendCount(c, date, query)
			ctx.Convey("Then err should be nil.count should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(count, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoDescTrendCount(t *testing.T) {
	convey.Convey("DescTrendCount", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			date  = "2018-06-01"
			query = ""
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			d.db.Exec(c, "INSERT INTO up_rating_trend_desc(mid) VALUES(1001) ON DUPLICATE KEY UPDATE mid=VALUES(mid)")
			count, err := d.DescTrendCount(c, date, query)
			ctx.Convey("Then err should be nil.count should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(count, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoGetRatingStatis(t *testing.T) {
	convey.Convey("GetRatingStatis", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			ctype = int64(0)
			date  = "2018-06-01"
			query = "2"
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			d.db.Exec(c, "INSERT INTO up_rating_statistics(cdate,tag_id,ctype,section) VALUES('2018-06-01',2,2,1) ON DUPLICATE KEY UPDATE tag_id=VALUES(tag_id)")
			statis, err := d.GetRatingStatis(c, ctype, date, query)
			ctx.Convey("Then err should be nil.statis should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(statis, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoGetTrendAsc(t *testing.T) {
	convey.Convey("GetTrendAsc", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			ctype = "magnetic"
			date  = "2018-06-01"
			query = ""
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			d.db.Exec(c, "INSERT INTO up_rating_trend_asc(mid,date) VALUES(1001,'2018-06-01') ON DUPLICATE KEY UPDATE mid=VALUES(mid)")
			ts, err := d.GetTrendAsc(c, ctype, date, query)
			ctx.Convey("Then err should be nil.ts should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(ts, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoGetTrendDesc(t *testing.T) {
	convey.Convey("GetTrendDesc", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			ctype = "magnetic"
			date  = "2018-06-01"
			query = ""
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			d.db.Exec(c, "INSERT INTO up_rating_trend_desc(mid,date) VALUES(1001,'2018-06-01') ON DUPLICATE KEY UPDATE mid=VALUES(mid)")
			ts, err := d.GetTrendDesc(c, ctype, date, query)
			ctx.Convey("Then err should be nil.ts should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(ts, convey.ShouldNotBeNil)
			})
		})
	})
}
