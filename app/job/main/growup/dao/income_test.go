package dao

import (
	"context"
	"testing"
	"time"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoGetAvTagRatio(t *testing.T) {
	convey.Convey("GetAvTagRatio", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			from  = int64(0)
			limit = int64(10)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			Exec(c, "INSERT INTO av_charge_ratio(tag_id,av_id) VALUES(1,2) ON DUPLICATE KEY UPDATE tag_id=VALUES(tag_id), av_id=VALUES(av_id)")
			infos, err := d.GetAvTagRatio(c, from, limit)
			ctx.Convey("Then err should be nil.infos should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(infos, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoGetAvIncomeInfo(t *testing.T) {
	convey.Convey("GetAvIncomeInfo", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			avID = int64(1)
			date = time.Date(2018, 6, 24, 0, 0, 0, 0, time.Local)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			Exec(c, "INSERT INTO av_income(av_id,mid,income,date) VALUES(1,2,3,'2018-06-24') ON DUPLICATE KEY UPDATE income=VALUES(income), av_id=VALUES(av_id)")
			info, err := d.GetAvIncomeInfo(c, avID, date)
			ctx.Convey("Then err should be nil.info should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(info, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoTxInsertTagIncome(t *testing.T) {
	convey.Convey("TxInsertTagIncome", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			tx, _ = d.BeginTran(c)
			sql   = "(1, 2, 3, 4, 5, '2018-06-23')"
		)
		defer tx.Commit()
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			rows, err := d.TxInsertTagIncome(tx, sql)
			ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rows, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoGetTagAvTotalIncome(t *testing.T) {
	convey.Convey("GetTagAvTotalIncome", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			tagID = int64(2)
			avID  = int64(1)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			Exec(c, "INSERT INTO up_tag_income (tag_id, av_id, total_income, date) VALUES(2, 1, 100,'2018-06-24') ON DUPLICATE KEY UPDATE tag_id=VALUES(tag_id),av_id=VALUES(av_id)")
			infos, err := d.GetTagAvTotalIncome(c, tagID, avID)
			ctx.Convey("Then err should be nil.infos should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(infos, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoListAvIncome(t *testing.T) {
	convey.Convey("ListAvIncome", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			id    = int64(0)
			limit = int(100)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			avIncome, err := d.ListAvIncome(c, id, limit)
			ctx.Convey("Then err should be nil.avIncome should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(avIncome, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoListUpAccount(t *testing.T) {
	convey.Convey("ListUpAccount", t, func(ctx convey.C) {
		var (
			c            = context.Background()
			withdrawDate = "2018-06"
			ctime        = "2018-06-23"
			from         = int(0)
			limit        = int(100)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			upAct, err := d.ListUpAccount(c, withdrawDate, ctime, from, limit)
			ctx.Convey("Then err should be nil.upAct should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(upAct, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoListUpIncome(t *testing.T) {
	convey.Convey("ListUpIncome", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			table = "up_income"
			date  = "2018-06-23"
			id    = int64(0)
			limit = int(100)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			um, err := d.ListUpIncome(c, table, date, id, limit)
			ctx.Convey("Then err should be nil.um should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(um, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoListUpWithdraw(t *testing.T) {
	convey.Convey("ListUpWithdraw", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			date  = "2018-06-23"
			from  = int(0)
			limit = int(100)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			ups, err := d.ListUpWithdraw(c, date, from, limit)
			ctx.Convey("Then err should be nil.ups should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(ups, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoGetUpTotalIncome(t *testing.T) {
	convey.Convey("GetUpTotalIncome", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			from  = int64(0)
			limit = int64(100)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			infos, err := d.GetUpTotalIncome(c, from, limit)
			ctx.Convey("Then err should be nil.infos should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(infos, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoGetUpIncome(t *testing.T) {
	convey.Convey("GetUpIncome", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			date  = time.Now()
			from  = int64(0)
			limit = int64(100)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			infos, err := d.GetUpIncome(c, date, from, limit)
			ctx.Convey("Then err should be nil.infos should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(infos, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoGetAvIncome(t *testing.T) {
	convey.Convey("GetAvIncome", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			date  = time.Now()
			id    = int64(0)
			limit = int64(10)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			infos, err := d.GetAvIncome(c, date, id, limit)
			ctx.Convey("Then err should be nil.infos should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(infos, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoGetUpTotalIncomeCnt(t *testing.T) {
	convey.Convey("GetUpTotalIncomeCnt", t, func(ctx convey.C) {
		var (
			c = context.Background()
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			upCnt, err := d.GetUpTotalIncomeCnt(c)
			ctx.Convey("Then err should be nil.upCnt should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(upCnt, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoGetAvStatisCount(t *testing.T) {
	convey.Convey("GetAvStatisCount", t, func(ctx convey.C) {
		var (
			c = context.Background()
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			cnt, err := d.GetAvStatisCount(c)
			ctx.Convey("Then err should be nil.cnt should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(cnt, convey.ShouldNotBeNil)
			})
		})
	})
}
