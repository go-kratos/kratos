package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoInsertBGMIncomeStatis(t *testing.T) {
	convey.Convey("InsertBGMIncomeStatis", t, func(ctx convey.C) {
		var (
			c      = context.Background()
			sid    = int64(100)
			income = int64(1090)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			Exec(c, "DELETE FROM bgm_income_statis WHERE sid=100")
			rows, err := d.InsertBGMIncomeStatis(c, sid, income)
			ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rows, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoGetBGMIncome(t *testing.T) {
	convey.Convey("GetBGMIncome", t, func(ctx convey.C) {
		var (
			c = context.Background()
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			statis, err := d.GetBGMIncome(c)
			ctx.Convey("Then err should be nil.statis should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(statis, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoGetCreditScore(t *testing.T) {
	convey.Convey("GetCreditScore", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			table = "video"
			id    = int64(100)
			limit = int64(200)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			scores, last, err := d.GetCreditScore(c, table, id, limit)
			ctx.Convey("Then err should be nil.scores,last should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(last, convey.ShouldNotBeNil)
				ctx.So(scores, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoSyncCreditScore(t *testing.T) {
	convey.Convey("SyncCreditScore", t, func(ctx convey.C) {
		var (
			c      = context.Background()
			values = "(100, 100)"
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			rows, err := d.SyncCreditScore(c, values)
			ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rows, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoGetAvBaseIncome(t *testing.T) {
	convey.Convey("GetAvBaseIncome", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			table = "up_income"
			id    = int64(0)
			limit = int64(10)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			d.db.Exec(c, "INSERT INTO up_income(mid,date,income) VALUES(1,'2018-06-24',2) ON DUPLICATE KEY UPDATE income=VALUES(income)")
			abs, last, err := d.GetAvBaseIncome(c, table, id, limit)
			ctx.Convey("Then err should be nil.abs,last should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(last, convey.ShouldNotBeNil)
				ctx.So(abs, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoBatchUpdateUpIncome(t *testing.T) {
	convey.Convey("BatchUpdateUpIncome", t, func(ctx convey.C) {
		var (
			c      = context.Background()
			table  = "up_income"
			values = "(100, '2018-06-23', 100)"
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			rows, err := d.BatchUpdateUpIncome(c, table, values)
			ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rows, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoGetAvs(t *testing.T) {
	convey.Convey("GetAvs", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			date = "2018-06-23"
			mid  = int64(100)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			avs, err := d.GetAvs(c, date, mid)
			ctx.Convey("Then err should be nil.avs should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(avs, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoGetAvCharges(t *testing.T) {
	convey.Convey("GetAvCharges", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			avIds = []int64{100, 200}
			date  = "2018-06-23"
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			charges, err := d.GetAvCharges(c, avIds, date)
			ctx.Convey("Then err should be nil.charges should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(charges, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoGetUpChargeRatio(t *testing.T) {
	convey.Convey("GetUpChargeRatio", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			tagID = int64(10)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			ups, err := d.GetUpChargeRatio(c, tagID)
			ctx.Convey("Then err should be nil.ups should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(ups, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoGetUpIncomeStatis(t *testing.T) {
	convey.Convey("GetUpIncomeStatis", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			mids = []int64{100, 200}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			ups, err := d.GetUpIncomeStatis(c, mids)
			ctx.Convey("Then err should be nil.ups should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(ups, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoGetUpIncomeDate(t *testing.T) {
	convey.Convey("GetUpIncomeDate", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			mids  = []int64{100, 200}
			table = "up_income"
			date  = "2018-06-23"
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			ups, err := d.GetUpIncomeDate(c, mids, table, date)
			ctx.Convey("Then err should be nil.ups should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(ups, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoUpdateDate(t *testing.T) {
	convey.Convey("UpdateDate", t, func(ctx convey.C) {
		var (
			tx, _ = d.BeginTran(context.Background())
			stmt  = "UPDATE up_info_video SET account_state=1 WHERE mid=10"
		)
		defer tx.Commit()
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			count, err := d.UpdateDate(tx, stmt)
			ctx.Convey("Then err should be nil.count should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(count, convey.ShouldNotBeNil)
			})
		})
	})
}
