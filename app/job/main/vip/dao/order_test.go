package dao

import (
	"context"
	"testing"

	"go-common/app/job/main/vip/model"
	xtime "go-common/library/time"

	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestDao_AddPayOrder(t *testing.T) {
	Convey("add pay order", t, func() {
		d.AddPayOrder(context.TODO(), &model.VipPayOrder{Mid: 123})
	})
}

func TestDao_UpdatePayOrderStatus(t *testing.T) {
	Convey("update order", t, func() {
		r := new(model.VipPayOrder)
		r.OrderType = 1
		r.PayType = 6
		r.Status = 2
		r.Ver = 2
		r.OrderNo = "1807021929153562939"
		r.Mtime = xtime.Time(time.Now().Unix())
		_, err := d.UpdatePayOrderStatus(context.TODO(), r)
		So(err, ShouldBeNil)

	})
}

func TestDao_SelPayOrderByMid(t *testing.T) {
	Convey("SelPayOrderByMid", t, func() {
		_, err := d.SelPayOrderByMid(context.TODO(), 7593623, 1, 1)
		So(err, ShouldBeNil)
	})
}

func TestDao_SSelOrderByOrderNo(t *testing.T) {
	Convey("SelOrderByOrderNo", t, func() {
		_, err := d.SelOrderByOrderNo(context.TODO(), "7593623")
		So(err, ShouldBeNil)
	})
}

func TestDao_SelPayOrderLog(t *testing.T) {
	Convey("SelPayOrderLog", t, func() {
		_, err := d.SelPayOrderLog(context.TODO(), 7593623, 1)
		So(err, ShouldBeNil)
	})
}

func TestDao_SelPayOrder(t *testing.T) {
	Convey("SelPayOrder", t, func() {
		_, err := d.SelPayOrder(context.TODO(), 1, 1)
		So(err, ShouldBeNil)
	})
}

func TestDao_SelOldPayOrder(t *testing.T) {
	Convey("SelOldPayOrder", t, func() {
		_, err := d.SelOldPayOrder(context.TODO(), 1, 1)
		So(err, ShouldBeNil)
	})
}

func TestDao_SelOldRechargeOrder(t *testing.T) {
	Convey("SelOldRechargeOrder", t, func() {
		_, err := d.SelOldRechargeOrder(context.TODO(), []string{"test"})
		So(err, ShouldBeNil)
	})
}
