package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoTotal(t *testing.T) {
	convey.Convey("Total", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			mon   = int(6)
			date  = "2018-06-01"
			where = ""
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			d.db.Exec(c, "INSERT INTO up_rating_06(mid,cdate,creativity_score) VALUES(1001, '2018-06-01', 100) ON DUPLICATE KEY UPDATE creativity_score=VALUES(creativity_score)")
			total, err := d.Total(c, mon, date, where)
			ctx.Convey("Then err should be nil.total should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(total, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoScoreList(t *testing.T) {
	convey.Convey("ScoreList", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			mon   = int(6)
			date  = "2018-06-01"
			where = ""
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			d.db.Exec(c, "INSERT INTO up_rating_06(mid,cdate,creativity_score) VALUES(1001, '2018-06-01', 100) ON DUPLICATE KEY UPDATE creativity_score=VALUES(creativity_score)")
			list, err := d.ScoreList(c, mon, date, where)
			ctx.Convey("Then err should be nil.list should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(list, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoLevelList(t *testing.T) {
	convey.Convey("LevelList", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			mon  = int(6)
			date = "2018-06-01"
			mids = []int64{1001}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			d.db.Exec(c, "INSERT INTO up_rating_06(mid,cdate,creativity_score) VALUES(1001, '2018-06-01', 100) ON DUPLICATE KEY UPDATE creativity_score=VALUES(creativity_score)")
			list, err := d.LevelList(c, mon, date, mids)
			ctx.Convey("Then err should be nil.list should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(list, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoUpScore(t *testing.T) {
	convey.Convey("UpScore", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			mon  = int(6)
			mid  = int64(1001)
			date = "2018-06-01"
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.UpScore(c, mon, mid, date)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoTaskStatus(t *testing.T) {
	convey.Convey("TaskStatus", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			date = "2018-06-01"
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			d.db.Exec(c, "INSERT INTO task_status(date,type) VALUES('2018-06-01', 2) ON DUPLICATE KEY UPDATE type=VALUES(type)")
			status, err := d.TaskStatus(c, date)
			ctx.Convey("Then err should be nil.status should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(status, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoUpScores(t *testing.T) {
	convey.Convey("UpScores", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mon = int(6)
			mid = int64(1001)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			d.db.Exec(c, "INSERT INTO up_rating_06(mid,cdate,creativity_score) VALUES(1001, '2018-06-01', 100) ON DUPLICATE KEY UPDATE creativity_score=VALUES(creativity_score)")
			list, err := d.UpScores(c, mon, mid)
			ctx.Convey("Then err should be nil.list should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(list, convey.ShouldNotBeNil)
			})
		})
	})
}
