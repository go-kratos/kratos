package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoBGMCount(t *testing.T) {
	convey.Convey("BGMCount", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(1001)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			Exec(c, "INSERT INTO background_music(sid, mid) VALUES(1000, 1001)")
			count, err := d.BGMCount(c, mid)
			ctx.Convey("Then err should be nil.count should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(count, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoGetBgmTitle(t *testing.T) {
	convey.Convey("GetBgmTitle", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			ids = []int64{1001}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			Exec(c, "INSERT INTO background_music(sid, mid, title) VALUES(1001, 1001, 'test')")
			titles, err := d.GetBgmTitle(c, ids)
			ctx.Convey("Then err should be nil.titles should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(titles, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoListBgmIncome(t *testing.T) {
	convey.Convey("ListBgmIncome", t, func(ctx convey.C) {
		var (
			c         = context.Background()
			mid       = int64(1001)
			startTime = "2018-01-01"
			endTime   = "2019-01-01"
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			Exec(c, "INSERT INTO bgm_income(sid, mid, total_income, date) VALUES(1000, 1001, 100, '2018-06-01')")
			bgms, err := d.ListBgmIncome(c, mid, startTime, endTime)
			ctx.Convey("Then err should be nil.bgms should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(bgms, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoListBgmIncomeByID(t *testing.T) {
	convey.Convey("ListBgmIncomeByID", t, func(ctx convey.C) {
		var (
			c       = context.Background()
			id      = int64(1000)
			endTime = "2019-01-01"
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			Exec(c, "INSERT INTO bgm_income(sid, mid, total_income, date) VALUES(1000, 1001, 100, '2018-06-01')")
			bgms, err := d.ListBgmIncomeByID(c, id, endTime)
			ctx.Convey("Then err should be nil.bgms should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(bgms, convey.ShouldNotBeNil)
			})
		})
	})
}
