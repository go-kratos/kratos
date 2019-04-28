package dao

import (
	"context"
	"go-common/app/service/openplatform/ticket-sales/model"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestDao_PromoOrder(t *testing.T) {
	Convey("PromoOrder", t, func() {
		res, err := d.PromoOrder(context.TODO(), 1)
		So(err, ShouldBeNil)
		t.Logf("res:%v", res)
	})
}

func TestDao_CachePromoOrder(t *testing.T) {
	Convey("CachePromoOrder", t, func() {
		res, err := d.CachePromoOrder(context.TODO(), 1)
		So(err, ShouldBeNil)
		t.Logf("res:%v", res)
	})
}

func TestDao_AddCachePromoOrder(t *testing.T) {
	Convey("AddCachePromoOrder", t, func() {
		err := d.AddCachePromoOrder(context.TODO(), 1, &model.PromotionOrder{PromoID: 1})
		So(err, ShouldBeNil)
	})
}

func TestDao_PromoOrderByStatus(t *testing.T) {
	Convey("PromoOrderByStatus", t, func() {
		res, err := d.PromoOrderByStatus(context.TODO(), 1, 1, 1, 1)
		So(err, ShouldBeNil)
		t.Logf("res:%v", res)
	})
}

func TestDao_PromoOrderDoing(t *testing.T) {
	Convey("PromoOrderDoing", t, func() {
		res, err := d.PromoOrderDoing(context.TODO(), 1, 1, 1)
		So(err, ShouldBeNil)
		t.Logf("res:%v", res)
	})
}

func TestDao_AddPromoOrder(t *testing.T) {
	Convey("AddPromoOrder", t, func() {
		res, err := d.AddPromoOrder(context.TODO(), 1, 1, 1, 1, 1, 1, 1, 1)
		So(err, ShouldBeNil)
		t.Logf("res:%v", res)
	})
}

func TestDao_UpdatePromoOrderStatus(t *testing.T) {
	Convey("UpdatePromoOrderStatus", t, func() {
		res, err := d.UpdatePromoOrderStatus(context.TODO(), 1, 1)
		So(err, ShouldBeNil)
		t.Logf("res:%v", res)
	})
}

func TestDao_GroupOrdersByStatus(t *testing.T) {
	Convey("GroupOrdersByStatus", t, func() {
		res, err := d.GroupOrdersByStatus(context.TODO(), 1, 1)
		So(err, ShouldBeNil)
		t.Logf("res:%v", res)
	})
}
