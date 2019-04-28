package dao

import (
	"context"
	"fmt"
	"net"
	"strings"
	"testing"
	"time"

	"go-common/app/service/openplatform/ticket-sales/api/grpc/type"
	"go-common/app/service/openplatform/ticket-sales/model"
	"go-common/app/service/openplatform/ticket-sales/model/consts"

	"github.com/smartystreets/goconvey/convey"
)

var (
	oids []int64
)

func TestTxInsertOrders(t *testing.T) {
	convey.Convey("OrderID", t, func() {
		oids, _ = d.OrderID(context.TODO(), 3)
		convey.So(len(oids), convey.ShouldEqual, 3)
	})
	convey.Convey("TxInsertOrders", t, func() {
		o1 := &model.OrderMain{
			OrderID:   oids[0],
			UID:       "TEST",
			OrderType: 1,
			ItemID:    1,
			ItemInfo: &_type.OrderItemInfo{
				Name:           "TEST",
				Img:            "//img.bilibili.com",
				ScreenID:       1,
				ScreenName:     "TEST",
				ScreenType:     1,
				DeliverType:    1,
				ExpressFee:     1,
				VIPExpressFree: 1,
				VerID:          1,
			},
			Count:        1,
			TotalMoney:   2,
			ExpressFee:   1,
			PayMoney:     2,
			PayChannel:   1,
			PayTime:      time.Now().Unix(),
			Source:       "TEST",
			Status:       1,
			SubStatus:    1,
			RefundStatus: 1,
		}
		tx, _ := d.BeginTx(context.TODO())
		c, _ := d.TxInsertOrders(tx, []*model.OrderMain{o1})
		convey.So(c, convey.ShouldEqual, 1)
		o2, o3 := &model.OrderMain{}, &model.OrderMain{}
		*o2 = *o1
		o2.OrderID = oids[1]
		o2.ItemID = 2
		*o3 = *o1
		o3.OrderID = oids[2]
		o3.ItemID = 2
		o3.RefundStatus = consts.RefundStatusPtRefunded
		c, _ = d.TxInsertOrders(tx, []*model.OrderMain{o2, o3})
		convey.So(c, convey.ShouldEqual, 2)
		tx.Commit()
	})
}

func TestOrders(t *testing.T) {
	convey.Convey("UserOrders", t, func() {
		orders, _ := d.CacheOrders(context.TODO(), &model.OrderMainQuerier{UID: "TEST", OrderID: oids[:1]})
		convey.So(len(orders), convey.ShouldEqual, 0)
		orders, _ = d.Orders(context.TODO(), &model.OrderMainQuerier{UID: "TEST"})
		convey.So(len(orders), convey.ShouldEqual, 3)
		for i := 0; i < 10; i++ {
			orders, _ = d.CacheOrders(context.TODO(), &model.OrderMainQuerier{UID: "TEST"})
			if len(orders) > 0 {
				break
			}
			time.Sleep(100 * time.Millisecond)
		}
		convey.So(len(orders), convey.ShouldEqual, 3)
		for i := 0; i < 10; i++ {
			orders, _ = d.CacheOrders(context.TODO(), &model.OrderMainQuerier{UID: "TEST", OrderID: oids[1:2]})
			if len(orders) > 0 {
				break
			}
			time.Sleep(100 * time.Millisecond)
		}
		convey.So(len(orders), convey.ShouldEqual, 1)
		orders, _ = d.Orders(context.TODO(), &model.OrderMainQuerier{UID: "TEST", ItemID: 2, Status: []int16{1}})
		convey.So(len(orders), convey.ShouldEqual, 2)
	})
}

func TestOrderCount(t *testing.T) {
	convey.Convey("OrderCount", t, func() {
		cnt, _ := d.CacheOrderCount(context.TODO(), &model.OrderMainQuerier{UID: "TEST"})
		convey.So(cnt, convey.ShouldEqual, 0)
		cnt, _ = d.OrderCount(context.TODO(), &model.OrderMainQuerier{UID: "TEST"})
		convey.So(cnt, convey.ShouldEqual, 3)
		for i := 0; i < 10; i++ {
			cnt, _ = d.CacheOrderCount(context.TODO(), &model.OrderMainQuerier{UID: "TEST"})
			if cnt > 0 {
				break
			}
			time.Sleep(100 * time.Millisecond)
		}
		convey.So(cnt, convey.ShouldEqual, 3)
	})
}

func TestTxInsertOrderDetails(t *testing.T) {
	convey.Convey("TxInsertOrderDetails", t, func() {
		dt1 := &model.OrderDetail{
			OrderID:    oids[0],
			Buyer:      "TEST",
			Tel:        "13800138000",
			PersonalID: "342921",
			ExpressCO:  "shunfeng",
			ExpressNO:  "000",
			Remark:     "TEST",
			DeviceType: 1,
			IP:         net.ParseIP("127.0.0.1"),
			DeliverDetail: &_type.OrderDeliver{
				AddrID: 1,
				Name:   "张三",
				Tel:    "13810559189",
				Addr:   "北京市",
			},
			Detail: &_type.OrderExtra{
				AutoRecvTime:   time.Now().Unix(),
				DelayRecvTimes: 1,
			},
		}
		dt2 := &model.OrderDetail{}
		*dt2 = *dt1
		dt2.OrderID = oids[1]
		tx, _ := d.BeginTx(context.TODO())
		c, _ := d.TxInsertOrderDetails(tx, []*model.OrderDetail{dt1, dt2})
		tx.Commit()
		convey.So(c, convey.ShouldEqual, 2)
	})
}

func TestOrderDetails(t *testing.T) {
	convey.Convey("OrderDetails", t, func() {
		dts, _ := d.CacheOrderDetails(context.TODO(), oids)
		convey.So(len(dts), convey.ShouldEqual, 0)
		dts, _ = d.OrderDetails(context.TODO(), oids)
		convey.So(len(dts), convey.ShouldEqual, 2)
		for i := 0; i < 10; i++ {
			dts, _ = d.CacheOrderDetails(context.TODO(), oids)
			if len(dts) > 0 {
				break
			}
			time.Sleep(100 * time.Millisecond)
		}
		convey.So(len(dts), convey.ShouldEqual, 3)
		dts, _ = d.OrderDetails(context.TODO(), oids)
		convey.So(len(dts), convey.ShouldEqual, 2)
	})
}

func TestTxInsertOrderSKUs(t *testing.T) {
	convey.Convey("TxInsertOrderSKUs", t, func() {
		sku1 := &model.OrderSKU{
			OrderID:     oids[0],
			SKUID:       1,
			Count:       1,
			OriginPrice: 1,
			Price:       1,
			TicketType:  1,
			Discounts: &_type.OrderSKUDiscounts{
				Platform: map[int32]int64{1: 1},
			},
		}
		sku2 := &model.OrderSKU{
			OrderID:     oids[0],
			SKUID:       2,
			Count:       1,
			OriginPrice: 2,
			Price:       2,
			TicketType:  1,
		}
		tx, _ := d.BeginTx(context.TODO())
		cnt, _ := d.TxInsertOrderSKUs(tx, []*model.OrderSKU{sku1, sku2})
		tx.Commit()
		convey.So(cnt, convey.ShouldEqual, 2)
	})
}

func TestOrderSKUs(t *testing.T) {
	convey.Convey("OrderSKUs", t, func() {
		skus, _ := d.CacheOrderSKUs(context.TODO(), oids)
		convey.So(len(skus), convey.ShouldEqual, 0)
		skus, _ = d.OrderSKUs(context.TODO(), oids)
		convey.So(len(skus), convey.ShouldEqual, 1)
		for i := 0; i < 10; i++ {
			skus, _ = d.CacheOrderSKUs(context.TODO(), oids)
			if len(skus) > 0 {
				break
			}
			time.Sleep(100 * time.Millisecond)
		}
		convey.So(len(skus), convey.ShouldEqual, 3)
		skus, _ = d.OrderSKUs(context.TODO(), oids)
		convey.So(len(skus), convey.ShouldEqual, 1)
	})
}

func TestOrderPayCharges(t *testing.T) {
	convey.Convey("OrderPayCharges", t, func() {
		var oid int64 = 1519626683221383
		chs, _ := d.OrderPayCharges(context.TODO(), []int64{oid})
		convey.So(chs, convey.ShouldContainKey, oid)
	})
}

func TestGetBoughtCount(t *testing.T) {
	convey.Convey("GetBoughtCount", t, func() {
		c, _ := d.RawBoughtCount(context.TODO(), "TEST", 2, nil)
		convey.So(c, convey.ShouldEqual, 2)
		c, _ = d.RawBoughtCount(context.TODO(), "TEST", 1, []int64{1})
		convey.So(c, convey.ShouldEqual, 1)
	})
}

//删除测试数据
func TestDelOrders(t *testing.T) {
	tx, _ := d.BeginTx(ctx)
	tables := []string{"order_main", "order_detail", "order_sku", "order_pay_charge", "ticket"}
	w := strings.Repeat(",?", len(oids))[1:]
	a := make([]interface{}, len(oids))
	for k, v := range oids {
		a[k] = v
	}
	for _, t := range tables {
		f := "order_id"
		if t == "ticket" {
			f = "oid"
		}
		_, err := tx.Exec(fmt.Sprintf("DELETE FROM `%s` WHERE %s IN (%s)", t, f, w), a...)
		if err != nil {
			tx.Rollback()
			return
		}
	}
	tx.Commit()
	d.DelCacheOrders(context.TODO(), &model.OrderMainQuerier{UID: "TEST"})
}
