package income

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestIncomeGetArchiveStatis(t *testing.T) {
	convey.Convey("GetArchiveStatis", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			table = "av_income_daily_statis"
			query = "id > 0"
			from  = int(0)
			limit = int(10)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			Exec(c, "INSERT INTO av_income_daily_statis(avs,money_section,money_tips,income,category_id,cdate) VALUES(10, 1, '0-3',1, '2018-01-01')")
			avs, err := d.GetArchiveStatis(c, table, query, from, limit)
			ctx.Convey("Then err should be nil.avs should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(avs, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestIncomeGetArchiveIncome(t *testing.T) {
	convey.Convey("GetArchiveIncome av", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			id    = int64(0)
			query = "id > 0"
			from  = "2018-01-01"
			to    = "2019-01-01"
			limit = int(100)
			typ   = int(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			Exec(c, "INSERT INTO av_income(av_id,mid,date) VALUES(1001, 1000, '2018-05-01')")
			archs, err := d.GetArchiveIncome(c, id, query, from, to, limit, typ)
			ctx.Convey("Then err should be nil.archs should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(archs, convey.ShouldNotBeNil)
			})
		})
	})

	convey.Convey("GetArchiveIncome column", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			id    = int64(0)
			query = "id > 0"
			from  = "2018-01-01"
			to    = "2019-01-01"
			limit = int(100)
			typ   = int(2)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			Exec(c, "INSERT INTO column_income(aid,mid,date) VALUES(1002, 1000, '2018-05-01')")
			archs, err := d.GetArchiveIncome(c, id, query, from, to, limit, typ)
			ctx.Convey("Then err should be nil.archs should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(archs, convey.ShouldNotBeNil)
			})
		})
	})

	convey.Convey("GetArchiveIncome bgm", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			id    = int64(0)
			query = "id > 0"
			from  = "2018-01-01"
			to    = "2019-01-01"
			limit = int(100)
			typ   = int(3)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			Exec(c, "INSERT INTO bgm_income(sid,mid,date) VALUES(1003, 1000, '2018-05-01')")
			archs, err := d.GetArchiveIncome(c, id, query, from, to, limit, typ)
			ctx.Convey("Then err should be nil.archs should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(archs, convey.ShouldNotBeNil)
			})
		})
	})

	convey.Convey("GetArchiveIncome error type", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			id    = int64(0)
			query = "id > 0"
			from  = "2018-01-01"
			to    = "2019-01-01"
			limit = int(100)
			typ   = int(4)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			_, err := d.GetArchiveIncome(c, id, query, from, to, limit, typ)
			ctx.Convey("Then err should be nil.archs should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestIncomeGetAvIncome(t *testing.T) {
	convey.Convey("GetAvIncome", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			id    = int64(0)
			query = "id > 0"
			from  = "2018-01-01"
			to    = "2019-01-01"
			limit = int(100)
			typ   = int(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			Exec(c, "INSERT INTO av_income(av_id,mid,date) VALUES(1001, 1000, '2018-05-01')")
			avs, err := d.GetAvIncome(c, id, query, from, to, limit, typ)
			ctx.Convey("Then err should be nil.avs should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(avs, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestIncomeGetColumnIncome(t *testing.T) {
	convey.Convey("GetColumnIncome", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			id    = int64(0)
			query = "id > 0"
			from  = "2018-01-01"
			to    = "2019-01-01"
			limit = int(100)
			typ   = int(2)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			Exec(c, "INSERT INTO column_income(aid,mid,date) VALUES(1002, 1000, '2018-05-01')")
			columns, err := d.GetColumnIncome(c, id, query, from, to, limit, typ)
			ctx.Convey("Then err should be nil.columns should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(columns, convey.ShouldNotBeNil)
			})
		})
	})
}
