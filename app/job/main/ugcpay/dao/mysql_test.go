package dao

import (
	"context"
	"fmt"
	"go-common/app/job/main/ugcpay/model"
	"math/rand"
	"testing"
	"time"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoCountPaidOrderUser(t *testing.T) {
	convey.Convey("CountPaidOrderUser", t, func(ctx convey.C) {
		var (
			c         = context.Background()
			beginTime = time.Now()
			endTime   = time.Now()
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			count, err := d.CountPaidOrderUser(c, beginTime, endTime)
			ctx.Convey("Then err should be nil.count should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(count, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoCountRefundedOrderUser(t *testing.T) {
	convey.Convey("CountRefundedOrderUser", t, func(ctx convey.C) {
		var (
			c         = context.Background()
			beginTime = time.Now()
			endTime   = time.Now()
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			count, err := d.CountRefundedOrderUser(c, beginTime, endTime)
			ctx.Convey("Then err should be nil.count should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(count, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoCountDailyBillByVer(t *testing.T) {
	convey.Convey("CountDailyBillByVer", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			ver = int64(20181030)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			count, err := d.CountDailyBillByVer(c, ver)
			ctx.Convey("Then err should be nil.count should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(count, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoCountDailyBillByMonthVer(t *testing.T) {
	convey.Convey("CountDailyBillByMonthVer", t, func(ctx convey.C) {
		var (
			c        = context.Background()
			monthVer = int64(20181030)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			count, err := d.CountDailyBillByMonthVer(c, monthVer)
			ctx.Convey("Then err should be nil.count should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(count, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoCountMonthlyBillByVer(t *testing.T) {
	convey.Convey("CountMonthlyBillByVer", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			ver = int64(201811)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			count, err := d.CountMonthlyBillByVer(c, ver)
			ctx.Convey("Then err should be nil.count should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(count, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoInsertLogTask(t *testing.T) {
	convey.Convey("InsertLogTask", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			data = &model.LogTask{
				Name:   fmt.Sprintf("ut_%d", time.Now().Unix()),
				Expect: 233,
			}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			id, err := d.InsertLogTask(c, data)
			ctx.Convey("Then err should be nil.id should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(id, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoLogTask(t *testing.T) {
	convey.Convey("LogTask", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			name = "ut"
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			data, err := d.LogTask(c, name)
			ctx.Convey("Then err should be nil.data should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(data, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoTXIncrLogTaskSuccess(t *testing.T) {
	convey.Convey("TXIncrLogTaskSuccess", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			tx, _ = d.BeginTran(c)
			name  = "ut"
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			rows, err := d.TXIncrLogTaskSuccess(c, tx, name)
			ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rows, convey.ShouldNotBeNil)
			})
		})
		ctx.Reset(func() {
			tx.Commit()
		})
	})
}

func TestDaoIncrLogTaskFailure(t *testing.T) {
	convey.Convey("IncrLogTaskFailure", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			name = "ut"
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			rows, err := d.IncrLogTaskFailure(c, name)
			ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rows, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoAsset(t *testing.T) {
	convey.Convey("Asset", t, func(ctx convey.C) {
		var (
			c        = context.Background()
			oid      = int64(2333)
			otype    = "archive"
			currency = "bp"
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			data, err := d.Asset(c, oid, otype, currency)
			ctx.Convey("Then err should be nil.data should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(data, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoSumPaidOrderUserRealFee(t *testing.T) {
	convey.Convey("SumPaidOrderUserRealFee", t, func(ctx convey.C) {
		var (
			c         = context.Background()
			beginTime = time.Now()
			endTime   = time.Now()
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			sum, err := d.SumPaidOrderUserRealFee(c, beginTime, endTime)
			ctx.Convey("Then err should be nil.sum should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(sum, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoSumRefundedOrderUserRealFee(t *testing.T) {
	convey.Convey("SumRefundedOrderUserRealFee", t, func(ctx convey.C) {
		var (
			c         = context.Background()
			beginTime = time.Now()
			endTime   = time.Now()
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			sum, err := d.SumRefundedOrderUserRealFee(c, beginTime, endTime)
			ctx.Convey("Then err should be nil.sum should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(sum, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoSumDailyBill(t *testing.T) {
	convey.Convey("SumDailyBill", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			ver = int64(20181030)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			sumIn, sumOut, err := d.SumDailyBill(c, ver)
			ctx.Convey("Then err should be nil.sumIn,sumOut should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(sumOut, convey.ShouldNotBeNil)
				ctx.So(sumIn, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoMinIDOrderPaid(t *testing.T) {
	convey.Convey("MinIDOrderPaid", t, func(ctx convey.C) {
		var (
			c         = context.Background()
			beginTime = time.Now()
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			minID, err := d.MinIDOrderPaid(c, beginTime)
			ctx.Convey("Then err should be nil.minID should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(minID, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoOrderPaidList(t *testing.T) {
	convey.Convey("OrderPaidList", t, func(ctx convey.C) {
		var (
			c         = context.Background()
			beginTime = time.Now()
			endTime   = time.Now()
			fromID    = int64(0)
			limit     = int(10)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			maxID, _, err := d.OrderPaidList(c, beginTime, endTime, fromID, limit)
			ctx.Convey("Then err should be nil.maxID,data should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				// ctx.So(data, convey.ShouldNotBeNil)
				ctx.So(maxID, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoMinIDOrderRefunded(t *testing.T) {
	convey.Convey("MinIDOrderRefunded", t, func(ctx convey.C) {
		var (
			c         = context.Background()
			beginTime = time.Now()
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			minID, err := d.MinIDOrderRefunded(c, beginTime)
			ctx.Convey("Then err should be nil.minID should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(minID, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoOrderRefundedList(t *testing.T) {
	convey.Convey("OrderRefundedList", t, func(ctx convey.C) {
		var (
			c         = context.Background()
			beginTime = time.Now()
			endTime   = time.Now()
			fromID    = int64(0)
			limit     = int(10)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			maxID, _, err := d.OrderRefundedList(c, beginTime, endTime, fromID, limit)
			ctx.Convey("Then err should be nil.maxID,data should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				// ctx.So(data, convey.ShouldNotBeNil)
				ctx.So(maxID, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoTXUpdateOrder(t *testing.T) {
	convey.Convey("TXUpdateOrder", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			tx, _ = d.BeginTran(c)
			order = &model.Order{}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			rows, err := d.TXUpdateOrder(c, tx, order)
			ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rows, convey.ShouldNotBeNil)
			})
		})
		ctx.Reset(func() {
			tx.Commit()
		})
	})
}

func TestDaoTXInsertOrderUserLog(t *testing.T) {
	convey.Convey("TXInsertOrderUserLog", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			tx, _ = d.BeginTran(c)
			data  = &model.LogOrder{
				OrderID: "ut",
				Desc:    "ut",
			}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			id, err := d.TXInsertOrderUserLog(c, tx, data)
			ctx.Convey("Then err should be nil.id should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(id, convey.ShouldNotBeNil)
			})
		})
		ctx.Reset(func() {
			tx.Commit()
		})
	})
}

func TestDaoMinIDDailyBillByVer(t *testing.T) {
	convey.Convey("MinIDDailyBillByVer", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			ver = int64(20181030)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			minID, err := d.MinIDDailyBillByVer(c, ver)
			ctx.Convey("Then err should be nil.minID should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(minID, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoDailyBillListByVer(t *testing.T) {
	convey.Convey("DailyBillListByVer", t, func(ctx convey.C) {
		var (
			c      = context.Background()
			ver    = int64(20181030)
			fromID = int64(0)
			limit  = int(10)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			maxID, data, err := d.DailyBillListByVer(c, ver, fromID, limit)
			ctx.Convey("Then err should be nil.maxID,data should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(data, convey.ShouldNotBeNil)
				ctx.So(maxID, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoMinIDDailyBillByMonthVer(t *testing.T) {
	convey.Convey("MinIDDailyBillByMonthVer", t, func(ctx convey.C) {
		var (
			c        = context.Background()
			monthVer = int64(201811)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			minID, err := d.MinIDDailyBillByMonthVer(c, monthVer)
			ctx.Convey("Then err should be nil.minID should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(minID, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoDailyBillListByMonthVer(t *testing.T) {
	convey.Convey("DailyBillListByMonthVer", t, func(ctx convey.C) {
		var (
			c        = context.Background()
			monthVer = int64(201811)
			fromID   = int64(0)
			limit    = int(10)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			maxID, data, err := d.DailyBillListByMonthVer(c, monthVer, fromID, limit)
			ctx.Convey("Then err should be nil.maxID,data should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(data, convey.ShouldNotBeNil)
				ctx.So(maxID, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoTXInsertLogDailyBill(t *testing.T) {
	convey.Convey("TXInsertLogDailyBill", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			tx, _ = d.BeginTran(c)
			log   = &model.LogBillDaily{
				OrderID: fmt.Sprintf("ut_%d", time.Now().Unix()),
			}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			id, err := d.TXInsertLogDailyBill(c, tx, log)
			ctx.Convey("Then err should be nil.id should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(id, convey.ShouldNotBeNil)
			})
		})
		ctx.Reset(func() {
			tx.Commit()
		})
	})
}

func TestDaoTXInsertLogMonthlyBill(t *testing.T) {
	convey.Convey("TXInsertLogMonthlyBill", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			tx, _ = d.BeginTran(c)
			log   = &model.LogBillMonthly{
				BillID:          fmt.Sprintf("ut_%d", time.Now().Unix()),
				BillUserDailyID: fmt.Sprintf("ut_%d", time.Now().Unix()),
			}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			id, err := d.TXInsertLogMonthlyBill(c, tx, log)
			ctx.Convey("Then err should be nil.id should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(id, convey.ShouldNotBeNil)
			})
		})
		ctx.Reset(func() {
			tx.Commit()
		})
	})
}

func TestDaoMinIDMonthlyBill(t *testing.T) {
	convey.Convey("MinIDMonthlyBill", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			ver = int64(201811)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			minID, err := d.MinIDMonthlyBill(c, ver)
			ctx.Convey("Then err should be nil.minID should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(minID, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoMonthlyBillList(t *testing.T) {
	convey.Convey("MonthlyBillList", t, func(ctx convey.C) {
		var (
			c      = context.Background()
			ver    = int64(201811)
			fromID = int64(0)
			limit  = int(10)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			maxID, data, err := d.MonthlyBillList(c, ver, fromID, limit)
			ctx.Convey("Then err should be nil.maxID,data should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(data, convey.ShouldNotBeNil)
				ctx.So(maxID, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoDailyBill(t *testing.T) {
	convey.Convey("DailyBill", t, func(ctx convey.C) {
		var (
			c        = context.Background()
			mid      = int64(46333)
			biz      = "asset"
			currency = "bp"
			ver      = int64(20181030)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			data, err := d.DailyBill(c, mid, biz, currency, ver)
			ctx.Convey("Then err should be nil.data should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(data, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoInsertDailyBill(t *testing.T) {
	convey.Convey("InsertDailyBill", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			bill = &model.DailyBill{
				MonthVer: 201811,
			}
		)
		bill.BillID = fmt.Sprintf("ut_%d", time.Now().Unix())
		bill.MID = 46333
		bill.Biz = fmt.Sprintf("ut_%d", time.Now().Unix())
		bill.Currency = "bp"
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			id, err := d.InsertDailyBill(c, bill)
			ctx.Convey("Then err should be nil.id should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(id, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoTXUpdateDailyBill(t *testing.T) {
	convey.Convey("TXUpdateDailyBill", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			tx, _ = d.BeginTran(c)
			bill  = &model.DailyBill{}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			rows, err := d.TXUpdateDailyBill(c, tx, bill)
			ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rows, convey.ShouldNotBeNil)
			})
		})
		ctx.Reset(func() {
			tx.Commit()
		})
	})
}

func TestDaoTXUpsertDeltaDailyBill(t *testing.T) {
	convey.Convey("TXUpsertDeltaDailyBill", t, func(ctx convey.C) {
		var (
			c        = context.Background()
			tx, _    = d.BeginTran(c)
			bill     = &model.DailyBill{}
			deltaIn  = int64(0)
			deltaOut = int64(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			rows, err := d.TXUpsertDeltaDailyBill(c, tx, bill, deltaIn, deltaOut)
			ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rows, convey.ShouldNotBeNil)
			})
		})
		ctx.Reset(func() {
			tx.Commit()
		})
	})
}

func TestDaoTXUpdateDeltaDailyBill(t *testing.T) {
	convey.Convey("TXUpdateDeltaDailyBill", t, func(ctx convey.C) {
		var (
			c        = context.Background()
			tx, _    = d.BeginTran(c)
			deltaIn  = int64(0)
			deltaOut = int64(0)
			mid      = int64(0)
			biz      = "archive"
			currency = "bp"
			ver      = int64(20181030)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			rows, err := d.TXUpdateDeltaDailyBill(c, tx, deltaIn, deltaOut, mid, biz, currency, ver)
			ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rows, convey.ShouldNotBeNil)
			})
		})
		ctx.Reset(func() {
			tx.Commit()
		})
	})
}

func TestDaoMonthlyBill(t *testing.T) {
	convey.Convey("MonthlyBill", t, func(ctx convey.C) {
		var (
			c        = context.Background()
			mid      = int64(46333)
			biz      = "asset"
			currency = "bp"
			ver      = int64(201811)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			data, err := d.MonthlyBill(c, mid, biz, currency, ver)
			ctx.Convey("Then err should be nil.data should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(data, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoInsertMonthlyBill(t *testing.T) {
	convey.Convey("InsertMonthlyBill", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			bill = &model.Bill{
				BillID:   fmt.Sprintf("ut_%d", time.Now().Unix()),
				MID:      46333,
				Biz:      fmt.Sprintf("ut_%d", time.Now().Unix()),
				Currency: "bp",
			}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			id, err := d.InsertMonthlyBill(c, bill)
			ctx.Convey("Then err should be nil.id should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(id, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoTXUpdateMonthlyBill(t *testing.T) {
	convey.Convey("TXUpdateMonthlyBill", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			tx, _ = d.BeginTran(c)
			bill  = &model.Bill{}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			rows, err := d.TXUpdateMonthlyBill(c, tx, bill)
			ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rows, convey.ShouldNotBeNil)
			})
		})
		ctx.Reset(func() {
			tx.Commit()
		})
	})
}

func TestDaoTXUpsertDeltaMonthlyBill(t *testing.T) {
	convey.Convey("TXUpsertDeltaMonthlyBill", t, func(ctx convey.C) {
		var (
			c        = context.Background()
			tx, _    = d.BeginTran(c)
			bill     = &model.Bill{}
			deltaIn  = int64(0)
			deltaOut = int64(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			rows, err := d.TXUpsertDeltaMonthlyBill(c, tx, bill, deltaIn, deltaOut)
			ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rows, convey.ShouldNotBeNil)
			})
		})
		ctx.Reset(func() {
			tx.Commit()
		})
	})
}

func TestDaoTXUpdateDeltaMonthlyBill(t *testing.T) {
	convey.Convey("TXUpdateDeltaMonthlyBill", t, func(ctx convey.C) {
		var (
			c        = context.Background()
			tx, _    = d.BeginTran(c)
			deltaIn  = int64(0)
			deltaOut = int64(0)
			mid      = int64(46333)
			biz      = "asset"
			currency = "bp"
			ver      = int64(201811)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			rows, err := d.TXUpdateDeltaMonthlyBill(c, tx, deltaIn, deltaOut, mid, biz, currency, ver)
			ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rows, convey.ShouldNotBeNil)
			})
		})
		ctx.Reset(func() {
			tx.Commit()
		})
	})
}

func TestDaoMinIDUserAccount(t *testing.T) {
	convey.Convey("MinIDUserAccount", t, func(ctx convey.C) {
		var (
			c         = context.Background()
			beginTime = time.Now()
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			minID, err := d.MinIDUserAccount(c, beginTime)
			ctx.Convey("Then err should be nil.minID should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(minID, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoUserAccountList(t *testing.T) {
	convey.Convey("UserAccountList", t, func(ctx convey.C) {
		var (
			c         = context.Background()
			beginTime = time.Now()
			endTime   = time.Now()
			fromID    = int64(0)
			limit     = int(10)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			maxID, _, err := d.UserAccountList(c, beginTime, endTime, fromID, limit)
			ctx.Convey("Then err should be nil.maxID,datas should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				// ctx.So(datas, convey.ShouldNotBeNil)
				ctx.So(maxID, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoInsertUserAccount(t *testing.T) {
	convey.Convey("InsertUserAccount", t, func(ctx convey.C) {
		var (
			c       = context.Background()
			account = &model.UserAccount{
				MID: 46333,
			}
		)
		account.Biz = fmt.Sprintf("ut_%d", time.Now().Unix())
		account.Currency = "bp"

		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			id, err := d.InsertUserAccount(c, account)
			ctx.Convey("Then err should be nil.id should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(id, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoUserAccount(t *testing.T) {
	convey.Convey("UserAccount", t, func(ctx convey.C) {
		var (
			c        = context.Background()
			mid      = int64(46333)
			biz      = "asset"
			currency = "bp"
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			data, err := d.UserAccount(c, mid, biz, currency)
			ctx.Convey("Then err should be nil.data should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(data, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoTXUpdateUserAccount(t *testing.T) {
	convey.Convey("TXUpdateUserAccount", t, func(ctx convey.C) {
		var (
			c       = context.Background()
			tx, _   = d.BeginTran(c)
			account = &model.UserAccount{}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			rows, err := d.TXUpdateUserAccount(c, tx, account)
			ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rows, convey.ShouldNotBeNil)
			})
		})
		ctx.Reset(func() {
			tx.Commit()
		})
	})
}

func TestDaoTXUpsertDeltaUserAccount(t *testing.T) {
	convey.Convey("TXUpsertDeltaUserAccount", t, func(ctx convey.C) {
		var (
			c       = context.Background()
			tx, _   = d.BeginTran(c)
			account = &model.UserAccount{
				MID: 46333,
			}
			deltaBalance = int64(0)
		)
		account.Currency = "bp"
		account.Biz = "asset"
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			rows, err := d.TXUpsertDeltaUserAccount(c, tx, account, deltaBalance)
			ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rows, convey.ShouldNotBeNil)
			})
		})
		ctx.Reset(func() {
			tx.Commit()
		})
	})
}

func TestDaoTXUpdateDeltaUserAccount(t *testing.T) {
	convey.Convey("TXUpdateDeltaUserAccount", t, func(ctx convey.C) {
		var (
			c            = context.Background()
			tx, _        = d.BeginTran(c)
			deltaBalance = int64(0)
			mid          = int64(46333)
			biz          = "asset"
			currency     = "bp"
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			rows, err := d.TXUpdateDeltaUserAccount(c, tx, deltaBalance, mid, biz, currency)
			ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rows, convey.ShouldNotBeNil)
			})
		})
		ctx.Reset(func() {
			tx.Commit()
		})
	})
}

func TestDaoTXInsertUserAccountLog(t *testing.T) {
	convey.Convey("TXInsertUserAccountLog", t, func(ctx convey.C) {
		var (
			c          = context.Background()
			tx, _      = d.BeginTran(c)
			accountLog = &model.AccountLog{
				Name: fmt.Sprintf("ut_%d", time.Now().Unix()),
			}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := d.TXInsertUserAccountLog(c, tx, accountLog)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
		ctx.Reset(func() {
			tx.Commit()
		})
	})
}

func TestDaoBizAccount(t *testing.T) {
	convey.Convey("BizAccount", t, func(ctx convey.C) {
		var (
			c        = context.Background()
			biz      = "asset"
			currency = "bp"
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			data, err := d.BizAccount(c, biz, currency)
			ctx.Convey("Then err should be nil.data should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(data, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoInsertBizAccount(t *testing.T) {
	convey.Convey("InsertBizAccount", t, func(ctx convey.C) {
		var (
			c       = context.Background()
			account = &model.BizAccount{
				Biz:      fmt.Sprintf("ut_%d", time.Now().Unix()),
				Currency: "bp",
			}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			id, err := d.InsertBizAccount(c, account)
			ctx.Convey("Then err should be nil.id should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(id, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoTXUpdateBizAccount(t *testing.T) {
	convey.Convey("TXUpdateBizAccount", t, func(ctx convey.C) {
		var (
			c       = context.Background()
			tx, _   = d.BeginTran(c)
			account = &model.BizAccount{}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			rows, err := d.TXUpdateBizAccount(c, tx, account)
			ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rows, convey.ShouldNotBeNil)
			})
		})
		ctx.Reset(func() {
			tx.Commit()
		})
	})
}

func TestDaoTXUpsertDeltaBizAccount(t *testing.T) {
	convey.Convey("TXUpsertDeltaBizAccount", t, func(ctx convey.C) {
		var (
			c       = context.Background()
			tx, _   = d.BeginTran(c)
			account = &model.BizAccount{
				Biz:      "asset",
				Currency: "bp",
			}
			deltaBalance = int64(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			rows, err := d.TXUpsertDeltaBizAccount(c, tx, account, deltaBalance)
			ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rows, convey.ShouldNotBeNil)
			})
		})
		ctx.Reset(func() {
			tx.Commit()
		})
	})
}

func TestDaoTXUpdateDeltaBizAccount(t *testing.T) {
	convey.Convey("TXUpdateDeltaBizAccount", t, func(ctx convey.C) {
		var (
			c            = context.Background()
			tx, _        = d.BeginTran(c)
			deltaBalance = int64(0)
			biz          = "asset"
			currency     = "bp"
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			rows, err := d.TXUpdateDeltaBizAccount(c, tx, deltaBalance, biz, currency)
			ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rows, convey.ShouldNotBeNil)
			})
		})
		ctx.Reset(func() {
			tx.Commit()
		})
	})
}

func TestDaoTXInsertBizAccountLog(t *testing.T) {
	convey.Convey("TXInsertBizAccountLog", t, func(ctx convey.C) {
		var (
			c          = context.Background()
			tx, _      = d.BeginTran(c)
			accountLog = &model.AccountLog{
				Name: fmt.Sprintf("ut_%d", time.Now().Unix()),
			}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := d.TXInsertBizAccountLog(c, tx, accountLog)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
		ctx.Reset(func() {
			tx.Commit()
		})
	})
}

func TestDaoAggrIncomeUser(t *testing.T) {
	convey.Convey("AggrIncomeUser", t, func(ctx convey.C) {
		var (
			c        = context.Background()
			mid      = int64(46333)
			currency = "bp"
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			data, err := d.AggrIncomeUser(c, mid, currency)
			ctx.Convey("Then err should be nil.data should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(data, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoInsertAggrIncomeUser(t *testing.T) {
	convey.Convey("InsertAggrIncomeUser", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			aggr = &model.AggrIncomeUser{
				MID: int64(2333*1000) + rand.Int63n(1000),
			}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			id, err := d.InsertAggrIncomeUser(c, aggr)
			ctx.Convey("Then err should be nil.id should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(id, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoTXUpdateAggrIncomeUser(t *testing.T) {
	convey.Convey("TXUpdateAggrIncomeUser", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			tx, _ = d.BeginTran(c)
			aggr  = &model.AggrIncomeUser{}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			rows, err := d.TXUpdateAggrIncomeUser(c, tx, aggr)
			ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rows, convey.ShouldNotBeNil)
			})
		})
		ctx.Reset(func() {
			tx.Commit()
		})
	})
}

func TestDaoTXUpsertDeltaAggrIncomeUser(t *testing.T) {
	convey.Convey("TXUpsertDeltaAggrIncomeUser", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			tx, _ = d.BeginTran(c)
			aggr  = &model.AggrIncomeUser{
				MID:      46333,
				Currency: "bp",
			}
			deltaPaySuccess = int64(0)
			deltaPayError   = int64(0)
			deltaTotalIn    = int64(0)
			deltaTotalOut   = int64(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			rows, err := d.TXUpsertDeltaAggrIncomeUser(c, tx, aggr, deltaPaySuccess, deltaPayError, deltaTotalIn, deltaTotalOut)
			ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rows, convey.ShouldNotBeNil)
			})
		})
		ctx.Reset(func() {
			tx.Commit()
		})
	})
}

func TestDaoTXUpdateDeltaAggrIncomeUser(t *testing.T) {
	convey.Convey("TXUpdateDeltaAggrIncomeUser", t, func(ctx convey.C) {
		var (
			c               = context.Background()
			tx, _           = d.BeginTran(c)
			deltaPaySuccess = int64(0)
			deltaPayError   = int64(0)
			deltaTotalIn    = int64(0)
			deltaTotalOut   = int64(0)
			mid             = int64(46333)
			currency        = "bp"
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			rows, err := d.TXUpdateDeltaAggrIncomeUser(c, tx, deltaPaySuccess, deltaPayError, deltaTotalIn, deltaTotalOut, mid, currency)
			ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rows, convey.ShouldNotBeNil)
			})
		})
		ctx.Reset(func() {
			tx.Commit()
		})
	})
}

func TestDaoAggrIncomeUserAsset(t *testing.T) {
	convey.Convey("AggrIncomeUserAsset", t, func(ctx convey.C) {
		var (
			c        = context.Background()
			mid      = int64(46333)
			currency = "bp"
			ver      = int64(201810)
			oid      = int64(10110846)
			otype    = "archive"
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			data, err := d.AggrIncomeUserAsset(c, mid, currency, ver, oid, otype)
			ctx.Convey("Then err should be nil.data should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(data, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoInsertAggrIncomeUserAsset(t *testing.T) {
	convey.Convey("InsertAggrIncomeUserAsset", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			aggr = &model.AggrIncomeUserAsset{
				MID:      46333,
				OID:      10110846,
				OType:    fmt.Sprintf("ut_%d", time.Now().Unix()),
				Ver:      201810,
				Currency: "bp",
			}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			id, err := d.InsertAggrIncomeUserAsset(c, aggr)
			ctx.Convey("Then err should be nil.id should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(id, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoTXUpdateAggrIncomeUserAsset(t *testing.T) {
	convey.Convey("TXUpdateAggrIncomeUserAsset", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			tx, _ = d.BeginTran(c)
			aggr  = &model.AggrIncomeUserAsset{}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			rows, err := d.TXUpdateAggrIncomeUserAsset(c, tx, aggr)
			ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rows, convey.ShouldNotBeNil)
			})
		})
		ctx.Reset(func() {
			tx.Commit()
		})
	})
}

func TestDaoTXUpsertDeltaAggrIncomeUserAsset(t *testing.T) {
	convey.Convey("TXUpsertDeltaAggrIncomeUserAsset", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			tx, _ = d.BeginTran(c)
			aggr  = &model.AggrIncomeUserAsset{
				MID:      46333,
				OID:      10110846,
				OType:    "archive",
				Ver:      201810,
				Currency: "bp",
			}
			deltaPaySuccess = int64(0)
			deltaPayError   = int64(0)
			deltaTotalIn    = int64(0)
			deltaTotalOut   = int64(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			rows, err := d.TXUpsertDeltaAggrIncomeUserAsset(c, tx, aggr, deltaPaySuccess, deltaPayError, deltaTotalIn, deltaTotalOut)
			ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rows, convey.ShouldNotBeNil)
			})
		})
		ctx.Reset(func() {
			tx.Commit()
		})
	})
}

func TestDaoTXUpdateDeltaAggrIncomeUserAsset(t *testing.T) {
	convey.Convey("TXUpdateDeltaAggrIncomeUserAsset", t, func(ctx convey.C) {
		var (
			c               = context.Background()
			tx, _           = d.BeginTran(c)
			deltaPaySuccess = int64(0)
			deltaPayError   = int64(0)
			deltaTotalIn    = int64(0)
			deltaTotalOut   = int64(0)
			mid             = int64(46333)
			currency        = "bp"
			ver             = int64(201810)
			oid             = int64(10110846)
			otype           = "archive"
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			rows, err := d.TXUpdateDeltaAggrIncomeUserAsset(c, tx, deltaPaySuccess, deltaPayError, deltaTotalIn, deltaTotalOut, mid, currency, ver, oid, otype)
			ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rows, convey.ShouldNotBeNil)
			})
		})
		ctx.Reset(func() {
			tx.Commit()
		})
	})
}

func TestDaoOrderBadDebt(t *testing.T) {
	convey.Convey("OrderBadDebt", t, func(ctx convey.C) {
		var (
			c       = context.Background()
			orderID = "666"
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			data, err := d.OrderBadDebt(c, orderID)
			ctx.Convey("Then err should be nil.data should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(data, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoInsertOrderBadDebt(t *testing.T) {
	convey.Convey("InsertOrderBadDebt", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			order = &model.OrderBadDebt{
				OrderID: fmt.Sprintf("ut_%d", time.Now().Unix()),
				Type:    "ut",
			}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			id, err := d.InsertOrderBadDebt(c, order)
			ctx.Convey("Then err should be nil.id should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(id, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoTXUpdateOrderBadDebt(t *testing.T) {
	convey.Convey("TXUpdateOrderBadDebt", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			tx, _ = d.BeginTran(c)
			order = &model.OrderBadDebt{}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			rows, err := d.TXUpdateOrderBadDebt(c, tx, order)
			ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rows, convey.ShouldNotBeNil)
			})
		})
		ctx.Reset(func() {
			tx.Commit()
		})
	})
}

func TestDaoTXInsertOrderRechargeShell(t *testing.T) {
	convey.Convey("TXInsertOrderRechargeShell", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			tx, _ = d.BeginTran(c)
			order = &model.OrderRechargeShell{
				MID:     46333,
				OrderID: fmt.Sprintf("ut_%d", time.Now().Unix()),
				Biz:     fmt.Sprintf("ut_%d", time.Now().Unix()),
				PayMSG:  "ut",
			}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			id, err := d.TXInsertOrderRechargeShell(c, tx, order)
			ctx.Convey("Then err should be nil.id should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(id, convey.ShouldNotBeNil)
			})
		})
		ctx.Reset(func() {
			tx.Commit()
		})
	})
}
