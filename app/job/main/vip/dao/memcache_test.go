package dao

import (
	"testing"

	"context"
	"go-common/app/job/main/vip/model"

	. "github.com/smartystreets/goconvey/convey"
)

func TestDao_DelVipInfoCache(t *testing.T) {
	Convey("should return true where err != nil and res not empty", t, func() {
		err := d.DelVipInfoCache(context.TODO(), 1234)
		So(err, ShouldBeNil)
	})
}

func TestDao_SetVipInfoCache(t *testing.T) {
	Convey("set vip info cache", t, func() {
		err := d.SetVipInfoCache(context.TODO(), 1234, &model.VipInfo{Mid: 1234, VipType: model.Vip, VipStatus: model.VipStatusNotOverTime})
		So(err, ShouldBeNil)
	})
}
func TestDao_AddPayOrderLog(t *testing.T) {
	Convey("add pay order log", t, func() {
		_, err := d.AddPayOrderLog(context.TODO(), &model.VipPayOrderLog{Mid: 1234, Status: 1, OrderNo: "12891723894189"})
		So(err, ShouldBeNil)
	})
}

func TestDao_GetVipMadelCache(t *testing.T) {
	Convey("get vip madel cache ", t, func() {
		val, err := d.GetVipMadelCache(context.TODO(), 2089801)
		t.Logf("val %v \n", val)
		So(err, ShouldBeNil)
	})
}

func TestDao_SetVipMadelCache(t *testing.T) {
	Convey("set vip madel", t, func() {
		err := d.SetVipMadelCache(context.TODO(), 2089801, 1)
		So(err, ShouldBeNil)
	})
}

func TestDao_SetVipFrozen(t *testing.T) {
	Convey("set vip frozen", t, func() {
		err := d.SetVipFrozen(context.TODO(), 2089809)
		So(err, ShouldBeNil)
	})
}

func TestDao_GetVipBuyCache(t *testing.T) {
	var val int64 = 1
	Convey("SetVipBuyCache", t, func() {
		err := d.SetVipBuyCache(context.TODO(), 2089809, val)
		So(err, ShouldBeNil)
	})
	Convey("GetVipBuyCache", t, func() {
		res, err := d.GetVipBuyCache(context.TODO(), 2089809)
		So(err, ShouldBeNil)
		So(val, ShouldEqual, res)
	})
}

func TestDao_AddTransferLock(t *testing.T) {
	Convey("AddTransferLock", t, func() {
		err := d.delCache(context.TODO(), "2089809")
		So(err, ShouldBeNil)
		bool := d.AddTransferLock(context.TODO(), "2089809")
		So(bool, ShouldBeTrue)
	})
}
