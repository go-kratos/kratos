package dao

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"testing"

	"go-common/app/service/main/ugcpay/model"
	xsql "go-common/library/database/sql"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoBeginTran(t *testing.T) {
	convey.Convey("BeginTran", t, func(ctx convey.C) {
		var (
			c = context.Background()
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			tx, err := d.BeginTran(c)
			ctx.Convey("Then err should be nil.tx should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(tx, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoRawOrderUser(t *testing.T) {
	convey.Convey("RawOrderUser", t, func(ctx convey.C) {
		var (
			c  = context.Background()
			id = ""
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			data, err := d.RawOrderUser(c, id)
			ctx.Convey("Then err should be nil.data should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(data, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoInsertOrderUser(t *testing.T) {
	convey.Convey("InsertOrderUser", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			data = &model.Order{
				OrderID:  strconv.FormatInt(time.Now().Unix(), 10),
				MID:      35858,
				Biz:      "asset",
				Platform: "ios",
				OID:      2333,
				OType:    "archive",
				Fee:      1300,
				Currency: "bp",
				State:    "created",
			}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			id, err := d.InsertOrderUser(c, data)
			ctx.Convey("Then err should be nil.id should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(id, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoRawAsset(t *testing.T) {
	convey.Convey("RawAsset", t, func(ctx convey.C) {
		var (
			c        = context.Background()
			oid      = int64(2333)
			otype    = "archive"
			currency = "bp"
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			data, err := d.RawAsset(c, oid, otype, currency)
			ctx.Convey("Then err should be nil.data should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(data, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoUpsertAsset(t *testing.T) {
	convey.Convey("UpsertAsset", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			data = &model.Asset{
				MID:      46333,
				OID:      2333,
				OType:    "archive",
				Currency: "bp",
				Price:    1000,
				State:    "valid",
			}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.UpsertAsset(c, data)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoRawAssetRelation(t *testing.T) {
	convey.Convey("RawAssetRelation", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			mid   = int64(35858)
			oid   = int64(2333)
			otype = "archive"
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			data, err := d.RawAssetRelation(c, mid, oid, otype)
			ctx.Convey("Then err should be nil.data should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(data, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoUpsertAssetRelation(t *testing.T) {
	convey.Convey("UpsertAssetRelation", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			data = &model.AssetRelation{
				OID:   2333,
				OType: "archive",
				MID:   35858,
				State: "paid",
			}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			_, err := d.UpsertAssetRelation(c, data)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoRawAggrIncomeUser(t *testing.T) {
	convey.Convey("RawAggrIncomeUser", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(46333)
			cur = "bp"
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			data, err := d.RawAggrIncomeUser(c, mid, cur)
			ctx.Convey("Then err should be nil.data should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(data, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoRawAggrIncomeUserMonthly(t *testing.T) {
	convey.Convey("RawAggrIncomeUserMonthly", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(46333)
			cur = "bp"
			ver = int64(201809)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			data, err := d.RawAggrIncomeUserAssetList(c, mid, cur, ver, 1000)
			ctx.Convey("Then err should be nil.data should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(data, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoRawAggrIncomeUserAsset(t *testing.T) {
	convey.Convey("RawAggrIncomeUserAsset", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			mid   = int64(46333)
			cur   = "bp"
			ver   = int64(201809)
			oid   = int64(10110842)
			otype = "archive"
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			data, err := d.RawAggrIncomeUserAsset(c, mid, cur, oid, otype, ver)
			ctx.Convey("Then err should be nil.data should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(data, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoRawBillUserDaily(t *testing.T) {
	convey.Convey("RawBillUserDaily", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(46333)
			biz = "asset"
			cur = "bp"
			ver = int64(20181030)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			data, err := d.RawBillUserDaily(c, mid, biz, cur, ver)
			ctx.Convey("Then err should be nil.data should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(data, convey.ShouldNotBeNil)
				t.Log(data)
			})
		})
	})
}

func TestDaoRawBillUserDailyListByMonthVer(t *testing.T) {
	convey.Convey("RawBillUserDailyListByMonthVer", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(46333)
			biz = "asset"
			cur = "bp"
			ver = int64(201811)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			datas, err := d.RawBillUserDailyByMonthVer(c, mid, biz, cur, ver)
			ctx.Convey("Then err should be nil.data should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(datas, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoTXInsertOrderUserLog(t *testing.T) {
	convey.Convey("TestDaoTXInsertOrderUserLog", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			log = &model.LogOrder{
				OrderID:   "test",
				FromState: "ut1",
				ToState:   "ut2",
				Desc:      "hahah",
			}
			tx *xsql.Tx
		)
		tx, err := d.BeginTran(c)
		convey.So(err, convey.ShouldBeNil)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			id, err := d.TXInsertOrderUserLog(c, tx, log)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(id, convey.ShouldBeGreaterThan, 0)
			})
		})
	})
}

func TestDaoOrderRechargeShell(t *testing.T) {
	convey.Convey("TestDaoOrderRechargeShell", t, func(ctx convey.C) {
		var (
			c       = context.Background()
			orderID = "23656019181109210220"
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			data, err := d.RawOrderRechargeShell(c, orderID)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(data, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoTXUpdateOrderRechargeShell(t *testing.T) {
	convey.Convey("TestDaoTXOrderRechargeShell", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			data = &model.OrderRechargeShell{
				MID:     46333,
				OrderID: "123",
				Biz:     "asset",
				Amount:  1,
				PayMSG:  "ut test",
				State:   "test",
				Ver:     201811,
			}
		)
		tx, err := d.BeginTran(c)
		convey.So(err, convey.ShouldBeNil)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.TXUpdateOrderRechargeShell(c, tx, data)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoTXInsertOrderRechargeShellLog(t *testing.T) {
	convey.Convey("TestDaoTXInsertOrderRechargeShellLog", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			data = &model.OrderRechargeShellLog{
				OrderID:           "233",
				FromState:         "ut1",
				ToState:           "ut2",
				Desc:              "ut test",
				BillUserMonthlyID: fmt.Sprintf("ut_%d", time.Now().Unix()),
			}
		)
		tx, err := d.BeginTran(c)
		convey.So(err, convey.ShouldBeNil)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			id, err := d.TXInsertOrderRechargeShellLog(c, tx, data)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(id, convey.ShouldBeGreaterThan, 0)
			})
		})
	})
}
func TestDaoTXUpdateOrderUser(t *testing.T) {
	convey.Convey("TXUpdateOrderUser", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			tx, _ = d.BeginTran(c)
			data  = &model.Order{
				OrderID: "96841175181102170810",
				State:   "ut test",
			}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			_, err := d.TXUpdateOrderUser(c, tx, data)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
		ctx.Reset(func() {
			tx.Commit()
		})
	})
}
