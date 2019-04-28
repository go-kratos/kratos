package dao

import (
	"context"
	"go-common/app/admin/main/vip/model"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoAddPushData(t *testing.T) {
	var (
		id  int64
		err error
	)
	convey.Convey("AddPushData", t, func() {
		id, err = d.AddPushData(context.TODO(), &model.VipPushData{Title: "test push"})
		convey.So(err, convey.ShouldBeNil)
		convey.So(id, convey.ShouldNotBeNil)
	})
	convey.Convey("GetPushData", t, func() {
		r, err := d.GetPushData(context.TODO(), id)
		convey.So(err, convey.ShouldBeNil)
		convey.So(r, convey.ShouldNotBeNil)
	})
	convey.Convey("UpdatePushData", t, func() {
		eff, err := d.UpdatePushData(context.TODO(), &model.VipPushData{ID: id, Title: "push test"})
		convey.So(err, convey.ShouldBeNil)
		convey.So(eff, convey.ShouldNotBeNil)
	})
	convey.Convey("PushDataCount", t, func() {
		count, err := d.PushDataCount(context.TODO(), &model.ArgPushData{})
		convey.So(err, convey.ShouldBeNil)
		convey.So(count, convey.ShouldNotBeNil)
	})
	convey.Convey("PushDatas", t, func() {
		res, err := d.PushDatas(context.TODO(), &model.ArgPushData{})
		convey.So(err, convey.ShouldBeNil)
		convey.So(res, convey.ShouldNotBeNil)
	})
	convey.Convey("DisablePushData", t, func() {
		err := d.DisablePushData(context.TODO(), &model.VipPushData{ID: id, Title: "push test"})
		convey.So(err, convey.ShouldBeNil)
	})
	convey.Convey("DelPushData", t, func() {
		err := d.DelPushData(context.TODO(), id)
		convey.So(err, convey.ShouldBeNil)
	})
}
