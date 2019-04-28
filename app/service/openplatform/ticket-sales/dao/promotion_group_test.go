package dao

import (
	"context"
	"go-common/app/service/openplatform/ticket-sales/model"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestDao_AddPromoGroup(t *testing.T) {
	Convey("AddPromoGroup", t, func() {
		res, err := d.AddPromoGroup(context.TODO(), 1, 1, 1, 1, 1, 1)
		So(err, ShouldBeNil)
		t.Logf("res:%v", res)
	})
}

func TestDao_RawPromoGroup(t *testing.T) {
	Convey("RawPromoGroup", t, func() {
		res, err := d.PromoGroup(context.TODO(), 1)
		So(err, ShouldBeNil)
		t.Logf("res:%v", res)
	})
}

func TestDao_CachePromoGroup(t *testing.T) {
	Convey("CachePromoGroup", t, func() {
		res, err := d.CachePromoGroup(context.TODO(), 1)
		So(err, ShouldBeNil)
		t.Logf("res:%v", res)
	})
}

func TestDao_AddCachePromoGroup(t *testing.T) {
	Convey("AddCachePromoGroup", t, func() {
		err := d.AddCachePromoGroup(context.TODO(), 1, &model.PromotionGroup{PromoID: 1})
		So(err, ShouldBeNil)
	})
}

func TestDao_DelCacheGroup(t *testing.T) {
	Convey("DelCacheGroup", t, func() {
		d.DelCachePromoGroup(context.TODO(), 1)
	})
}

func TestDao_GetUserGroupDoing(t *testing.T) {
	Convey("GetUserGroupDoing", t, func() {
		res, err := d.GetUserGroupDoing(context.TODO(), 1, 1, 1)
		So(err, ShouldBeNil)
		t.Logf("res:%v", res)
	})
}

func TestDao_UpdateGroupStatusAndOrderCount(t *testing.T) {
	Convey("UpdateGroupStatusAndOrderCount", t, func() {
		res, err := d.UpdateGroupStatusAndOrderCount(context.TODO(), 1, 1, 1, 1)
		So(err, ShouldBeNil)
		t.Logf("res:%v", res)
	})
}
