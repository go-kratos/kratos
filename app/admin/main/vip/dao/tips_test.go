package dao

import (
	"context"
	"testing"

	"go-common/app/admin/main/vip/model"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoTipList(t *testing.T) {
	convey.Convey("TipList", t, func() {
		rs, err := d.TipList(context.TODO(), 0, 0, 0, 0)
		convey.So(err, convey.ShouldBeNil)
		convey.So(rs, convey.ShouldNotBeNil)
	})
}

func TestDaoAddTip(t *testing.T) {
	var (
		id  int64
		err error
	)
	convey.Convey("AddTip", t, func() {
		id, err = d.AddTip(context.TODO(), &model.Tips{Tip: "test"})
		convey.So(err, convey.ShouldBeNil)
		convey.So(id, convey.ShouldNotBeNil)
	})
	convey.Convey("TipByID", t, func() {
		r, err := d.TipByID(context.TODO(), id)
		convey.So(err, convey.ShouldBeNil)
		convey.So(r, convey.ShouldNotBeNil)
	})
	convey.Convey("TipUpdate", t, func() {
		eff, err := d.TipUpdate(context.TODO(), &model.Tips{ID: id, Tip: "test"})
		convey.So(err, convey.ShouldBeNil)
		convey.So(eff, convey.ShouldNotBeNil)
	})
	convey.Convey("ExpireTip", t, func() {
		eff, err := d.ExpireTip(context.TODO(), id, "", 0)
		convey.So(err, convey.ShouldBeNil)
		convey.So(eff, convey.ShouldNotBeNil)
	})
	convey.Convey("DeleteTip", t, func() {
		eff, err := d.DeleteTip(context.TODO(), id, 0, "")
		convey.So(err, convey.ShouldBeNil)
		convey.So(eff, convey.ShouldNotBeNil)
	})
}
