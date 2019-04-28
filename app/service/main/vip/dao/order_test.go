package dao

import (
	"context"
	"fmt"
	"math/rand"
	"testing"
	"time"

	"go-common/app/service/main/vip/model"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoBeginTran(t *testing.T) {
	var (
		c = context.TODO()
	)
	convey.Convey("BeginTran", t, func(ctx convey.C) {
		tx, err := d.BeginTran(c)
		ctx.Convey("Tx Error should be nil", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
		ctx.Convey("tx should not be nil", func(ctx convey.C) {
			ctx.So(tx, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoOrderCount(t *testing.T) {
	var (
		c      = context.TODO()
		mid    = int64(20606508)
		status = int8(1)
	)
	convey.Convey("OrderCount", t, func(ctx convey.C) {
		count, err := d.OrderCount(c, mid, status)
		ctx.Convey("Error should be nil", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
		ctx.Convey("count should not be nil", func(ctx convey.C) {
			ctx.So(count, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoOrderList(t *testing.T) {
	var (
		c      = context.TODO()
		mid    = int64(20606508)
		status = int8(1)
		pn     = int(1)
		ps     = int(20)
	)
	convey.Convey("OrderList", t, func(ctx convey.C) {
		_, err := d.OrderList(c, mid, status, pn, ps)
		ctx.Convey("Error should be nil", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoOrderInfo(t *testing.T) {
	var (
		c       = context.TODO()
		orderNo = "12345"
	)
	convey.Convey("OrderInfo", t, func(ctx convey.C) {
		_, err := d.OrderInfo(c, orderNo)
		ctx.Convey("Error should be nil", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoDiscountSQL(t *testing.T) {
	var (
		c          = context.TODO()
		mid        = int64(20606508)
		discountID = int64(2)
	)
	convey.Convey("DiscountSQL", t, func(ctx convey.C) {
		_, err := d.DiscountSQL(c, mid, discountID)
		ctx.Convey("Error should be nil", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoPriceMapping(t *testing.T) {
	var (
		c        = context.TODO()
		monthID  = int64(1)
		platform = int8(1)
	)
	convey.Convey("PriceMapping", t, func(ctx convey.C) {
		_, err := d.PriceMapping(c, monthID, platform)
		ctx.Convey("Error should be nil", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoTxAddOrder(t *testing.T) {
	var (
		c = context.TODO()
		p = &model.PayOrder{
			OrderNo:      fmt.Sprintf("%d_%d", time.Now().Unix(), rand.Int31()),
			AppID:        1,
			Platform:     1,
			OrderType:    1,
			AppSubID:     "a",
			Mid:          20606508,
			ToMid:        1,
			BuyMonths:    1,
			Money:        1.0,
			Status:       1,
			PayType:      1,
			RechargeBp:   1.0,
			ThirdTradeNo: "209",
			Ver:          1,
			UserIP:       []byte("127.0.0.1"),
		}
	)
	convey.Convey("TxAddOrder", t, func(ctx convey.C) {
		tx, err := d.StartTx(c)
		ctx.Convey("Tx Error should be nil", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
		ctx.Convey("TxAddOrder Error should be nil", func(ctx convey.C) {
			id, err := d.TxAddOrder(tx, p)
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(id, convey.ShouldNotBeNil)
			err = tx.Commit()
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoAllMonthByOrder(t *testing.T) {
	var (
		c        = context.TODO()
		orderStr = ""
	)
	convey.Convey("AllMonthByOrder", t, func(ctx convey.C) {
		res, err := d.AllMonthByOrder(c, orderStr)
		ctx.Convey("Error should be nil", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
		ctx.Convey("res should not be nil", func(ctx convey.C) {
			ctx.So(res, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoTxUpdateOrderStatus(t *testing.T) {
	var (
		c            = context.TODO()
		status       = int8(0)
		payType      = ""
		thirdTradeNO = ""
		orderNO      = ""
	)
	convey.Convey("TxUpdateOrderStatus", t, func(ctx convey.C) {
		tx, err := d.StartTx(c)
		ctx.Convey("Tx Error should be nil", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
		defer tx.Commit()
		err = d.TxUpdateOrderStatus(c, tx, status, payType, thirdTradeNO, orderNO)
		ctx.Convey("Error should be nil", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoTxUpdatePayOrder(t *testing.T) {
	var (
		c = context.TODO()
		o = &model.OrderInfo{
			OrderNo:      "12345",
			AppID:        1,
			OrderType:    2,
			Platform:     1,
			Mid:          20606508,
			ToMid:        20606509,
			BuyMonths:    1,
			Money:        3.0,
			Status:       1,
			PayType:      "1",
			RechargeBP:   3.0,
			ThirdTradeNo: "456",
			Ver:          23,
			AppSubID:     "1234",
		}
		ver = int64(0)
	)
	convey.Convey("TxUpdatePayOrder", t, func(ctx convey.C) {
		tx, err := d.StartTx(c)
		ctx.Convey("Tx Error should be nil", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
		defer tx.Commit()
		err = d.TxUpdatePayOrder(tx, o, ver)
		ctx.Convey("Error should be nil", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoTxUpdateIosPayOrder(t *testing.T) {
	var (
		c = context.TODO()
		o = &model.OrderInfo{
			OrderNo:      "12345",
			AppID:        1,
			OrderType:    2,
			Platform:     1,
			Mid:          20606508,
			ToMid:        20606509,
			BuyMonths:    1,
			Money:        3.0,
			Status:       1,
			PayType:      "1",
			RechargeBP:   3.0,
			ThirdTradeNo: "456",
			Ver:          23,
			AppSubID:     "1234",
		}
		ver = int64(0)
	)
	convey.Convey("TxUpdateIosPayOrder", t, func(ctx convey.C) {
		tx, err := d.StartTx(c)
		ctx.Convey("Tx Error should be nil", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
		defer tx.Commit()
		err = d.TxUpdateIosPayOrder(tx, o, ver)
		ctx.Convey("Error should be nil", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoTxUpdatePayOrderStatus(t *testing.T) {
	var (
		c      = context.TODO()
		status = int8(0)
		id     = int64(0)
		ver    = int64(0)
	)
	convey.Convey("TxUpdatePayOrderStatus", t, func(ctx convey.C) {
		tx, err := d.StartTx(c)
		ctx.Convey("Tx Error should be nil", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
		defer tx.Commit()
		a, err := d.TxUpdatePayOrderStatus(tx, status, id, ver)
		ctx.Convey("Error should be nil", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
		ctx.Convey("a should not be nil", func(ctx convey.C) {
			ctx.So(a, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoTxAddOrderLog(t *testing.T) {
	var (
		c       = context.TODO()
		mid     = int64(0)
		orderNo = ""
		status  = int8(0)
	)
	convey.Convey("TxAddOrderLog", t, func(ctx convey.C) {
		tx, err := d.StartTx(c)
		ctx.Convey("Tx Error should be nil", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
		defer tx.Commit()
		olog := new(model.VipPayOrderLog)
		olog.Mid = mid
		olog.Status = status
		olog.OrderNo = orderNo
		err = d.TxAddOrderLog(tx, olog)
		ctx.Convey("Error should be nil", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoPayOrderLast(t *testing.T) {
	var (
		c          = context.TODO()
		mid        = int64(0)
		status     = int8(0)
		orderTypes = int64(0)
	)
	convey.Convey("PayOrderLast", t, func(ctx convey.C) {
		_, err := d.PayOrderLast(c, mid, status, orderTypes)
		ctx.Convey("Error should be nil", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoSelOldPayOrder(t *testing.T) {
	var (
		c       = context.TODO()
		orderNo = ""
	)
	convey.Convey("SelOldPayOrder", t, func(ctx convey.C) {
		_, err := d.SelOldPayOrder(c, orderNo)
		ctx.Convey("Error should be nil", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}
