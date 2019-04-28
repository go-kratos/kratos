package dao

import (
	"context"
	"go-common/app/service/openplatform/ticket-sales/model"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestDao_RawPromo(t *testing.T) {
	Convey("RawPromo", t, func() {
		res, err := d.RawPromo(context.TODO(), 1)
		So(err, ShouldBeNil)
		t.Logf("res:%v", res)
	})
}

func TestDao_CachePromo(t *testing.T) {
	Convey("CachePromo", t, func() {
		res, err := d.CachePromo(context.TODO(), 1)
		So(err, ShouldBeNil)
		t.Logf("res:%v", res)
	})
}

func TestDao_AddCachePromo(t *testing.T) {
	Convey("AddCachePromo", t, func() {
		err := d.AddCachePromo(context.TODO(), 1, &model.Promotion{PromoID: 1})
		So(err, ShouldBeNil)
	})
}

func TestDao_DelCachePromo(t *testing.T) {
	Convey("DelCachePromo", t, func() {
		d.DelCachePromo(context.TODO(), 1)
	})
}

func TestDao_CreatePromo(t *testing.T) {
	Convey("CreatePromo", t, func() {
		res, err := d.CreatePromo(context.TODO(), &model.Promotion{PromoID: 1})
		So(err, ShouldBeNil)
		t.Logf("res:%v", res)
	})
}
