package dao

import (
	"context"
	"testing"
	"time"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoDelTrend(t *testing.T) {
	convey.Convey("DelTrend", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			table = "asc"
			limit = int(100)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			d.db.Exec(c, "INSERT INTO up_rating_trend_asc(mid,tag_id) VALUES(1,2) ON DUPLICATE KEY UPDATE mid=VALUES(mid)")
			rows, err := d.DelTrend(c, table, limit)
			ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rows, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoInsertRatingStatis(t *testing.T) {
	convey.Convey("InsertRatingStatis", t, func(ctx convey.C) {
		var (
			c      = context.Background()
			values = "(1,1,'10-20',160,60,50,50,100,100,100,200,2,1,'2018-06-01')"
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			d.db.Exec(c, "DELETE FROM up_rating_statistics WHERE tag_id=2")
			_, err := d.InsertRatingStatis(c, values)
			ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoInsertTopRating(t *testing.T) {
	convey.Convey("InsertTopRating", t, func(ctx convey.C) {
		var (
			c      = context.Background()
			values = "(1,2,3,100,100,100,'2018-06-01')"
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			d.db.Exec(c, "DELETE FROM up_rating_top WHERE mid=1")
			_, err := d.InsertTopRating(c, values)
			ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoDelRatingCom(t *testing.T) {
	convey.Convey("DelRatingCom", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			table = "up_rating_top"
			date  = time.Now()
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			_, err := d.DelRatingCom(c, table, date, 2000)
			ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}
