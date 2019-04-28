package dao

import (
	"context"
	"go-common/app/service/main/tv/internal/model"
	"math/rand"
	"strconv"
	"testing"
	"time"

	xtime "go-common/library/time"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoPayOrderById(t *testing.T) {
	convey.Convey("PayOrderById", t, func(ctx convey.C) {
		var (
			c  = context.Background()
			id = 20
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			po, err := d.PayOrderByID(c, id)
			ctx.Convey("Then err should be nil.po should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(po, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoPayOrderByOrderNo(t *testing.T) {
	convey.Convey("PayOrderByOrderNo", t, func(ctx convey.C) {
		var (
			c       = context.Background()
			orderNo = "T123456789"
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			po, err := d.PayOrderByOrderNo(c, orderNo)
			ctx.Convey("Then err should be nil.po should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(po, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoPayOrdersByMid(t *testing.T) {
	convey.Convey("PayOrdersByMid", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int(27515308)
			pn  = int(1)
			ps  = int(10)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			res, total, err := d.PayOrdersByMid(c, mid, pn, ps)
			ctx.Convey("Then err should be nil.res,total should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(total, convey.ShouldNotBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoPayOrdersByMidAndStatus(t *testing.T) {
	convey.Convey("PayOrdersByMidAndStatus", t, func(ctx convey.C) {
		var (
			c      = context.Background()
			mid    = int(27515308)
			status = int8(1)
			pn     = int(1)
			ps     = int(10)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			res, total, err := d.PayOrdersByMidAndStatus(c, mid, status, pn, ps)
			ctx.Convey("Then err should be nil.res,total should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(total, convey.ShouldNotBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoPayOrdersByMidAndStatusAndCtime(t *testing.T) {
	convey.Convey("PayOrdersByMidAndStatusAndCtime", t, func(ctx convey.C) {
		var (
			c      = context.Background()
			mid    = int64(27515308)
			status = int8(1)
			from   = xtime.Time(time.Now().Add(-time.Hour * 24).Unix())
			to     = xtime.Time(time.Now().Unix())
			pn     = int(1)
			ps     = int(10)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			res, total, err := d.PayOrdersByMidAndStatusAndCtime(c, mid, status, from, to, pn, ps)
			ctx.Convey("Then err should be nil.res,total should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(total, convey.ShouldNotBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoTxInsertPayOrder(t *testing.T) {
	convey.Convey("TxInsertPayOrder", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			tx, _ = d.BeginTran(c)
			po    = &model.PayOrder{
				OrderNo:      "T123456789",
				Platform:     1,
				OrderType:    1,
				Mid:          27515308,
				BuyMonths:    1,
				ProductId:    "wx345678",
				Quantity:     1,
				Status:       1,
				PaymentMoney: 100,
				PaymentType:  "wechat",
				Ver:          1,
				Token:        "TOKEN:123456789",
				AppChannel:   "master",
			}
		)
		d.TxInsertPayOrder(c, tx, po)
		po.OrderNo = "T" + strconv.Itoa(rand.Int()/100000)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			id, err := d.TxInsertPayOrder(c, tx, po)
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

func TestDaoUnpaidNotCallbackOrder(t *testing.T) {
	convey.Convey("UnpaidNotCallbackOrder", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			stime xtime.Time
			etime = xtime.Time(time.Now().Unix())
			ps    = int(500)
			pn    = int(1)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			res, err := d.UnpaidNotCallbackOrder(c, stime, etime, pn, ps)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoTxUpdatePayOrder(t *testing.T) {
	convey.Convey("TxUpdatePayOrder", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			tx, _ = d.BeginTran(c)
		)
		po, _ := d.PayOrderByID(c, 60)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := d.TxUpdatePayOrder(c, tx, po)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
		ctx.Reset(func() {
			tx.Commit()
		})
	})
}
