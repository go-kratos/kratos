package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoGetUpDailyCharge(t *testing.T) {
	convey.Convey("GetUpDailyCharge", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			mid   = int64(1001)
			begin = "2018-01-01"
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			Exec(c, "INSERT INTO up_daily_charge(mid, date, inc_charge) VALUES(1001, '2018-06-01', 100)")
			incs, err := d.GetUpDailyCharge(c, mid, begin)
			ctx.Convey("Then err should be nil.incs should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(incs, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoListAvIncome(t *testing.T) {
	convey.Convey("ListAvIncome", t, func(ctx convey.C) {
		var (
			c         = context.Background()
			mid       = int64(1001)
			startTime = "2018-01-01"
			endTime   = "2019-01-01"
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			Exec(c, "INSERT INTO av_income(av_id, mid, date, income) VALUES(1000, 1001, '2018-06-01', 100)")
			avs, err := d.ListAvIncome(c, mid, startTime, endTime)
			ctx.Convey("Then err should be nil.avs should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(avs, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoListAvIncomeByID(t *testing.T) {
	convey.Convey("ListAvIncomeByID", t, func(ctx convey.C) {
		var (
			c       = context.Background()
			avID    = int64(1000)
			endTime = "2019-01-01"
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			Exec(c, "INSERT INTO av_income(av_id, mid, date, income) VALUES(1000, 1001, '2018-06-01', 100)")
			avs, err := d.ListAvIncomeByID(c, avID, endTime)
			ctx.Convey("Then err should be nil.avs should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(avs, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoListAvBlackList(t *testing.T) {
	convey.Convey("ListAvBlackList", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			avIds = []int64{1000}
			typ   = int(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			Exec(c, "INSERT INTO av_black_list(av_id, ctype) VALUES(1000, 0) ON DUPLICATE KEY UPDATE ctype = 0")
			avb, err := d.ListAvBlackList(c, avIds, typ)
			ctx.Convey("Then err should be nil.avb should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(avb, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoListActiveInfo(t *testing.T) {
	convey.Convey("ListActiveInfo", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			avIds = []int64{1000}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			Exec(c, "INSERT INTO activity_info(archive_id, tag_id) VALUES(1000, 1)")
			acM, err := d.ListActiveInfo(c, avIds)
			ctx.Convey("Then err should be nil.acM should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(acM, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoListTagInfo(t *testing.T) {
	convey.Convey("ListTagInfo", t, func(ctx convey.C) {
		var (
			c      = context.Background()
			tagIds = []int64{1000}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			Exec(c, "INSERT INTO tag_info(id, ratio, icon) VALUES(1000, 10, 'aaaaaa')")
			tagM, err := d.ListTagInfo(c, tagIds)
			ctx.Convey("Then err should be nil.tagM should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(tagM, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoListUpIncome(t *testing.T) {
	convey.Convey("ListUpIncome", t, func(ctx convey.C) {
		var (
			c         = context.Background()
			mid       = int64(1001)
			table     = "up_income"
			startTime = "2018-01-01"
			endTime   = "2019-01-01"
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			Exec(c, "INSERT INTO up_income(mid, date) VALUES(1001, '2018-06-01')")
			ups, err := d.ListUpIncome(c, mid, table, startTime, endTime)
			ctx.Convey("Then err should be nil.ups should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(ups, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoListUpAccount(t *testing.T) {
	convey.Convey("ListUpAccount", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(1001)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			Exec(c, "INSERT INTO up_account(mid, is_deleted) VALUES(1001, 0) ON DUPLICATE KEY UPDATE is_deleted = 0")
			up, err := d.ListUpAccount(c, mid)
			ctx.Convey("Then err should be nil.up should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(up, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoGetUpIncome(t *testing.T) {
	convey.Convey("GetUpIncome", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			mid   = int64(1001)
			begin = "2018-01-01"
			end   = "2019-01-01"
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			Exec(c, "INSERT INTO up_income(mid, date) VALUES(1001, '2018-06-01')")
			result, err := d.GetUpIncome(c, mid, begin, end)
			ctx.Convey("Then err should be nil.result should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(result, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoGetUpIncomeCount(t *testing.T) {
	convey.Convey("GetUpIncomeCount", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			date = "2018-06-01"
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			Exec(c, "INSERT INTO up_income(mid, date) VALUES(1001, '2018-06-01')")
			count, err := d.GetUpIncomeCount(c, date)
			ctx.Convey("Then err should be nil.count should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(count, convey.ShouldNotBeNil)
			})
		})
	})
}
