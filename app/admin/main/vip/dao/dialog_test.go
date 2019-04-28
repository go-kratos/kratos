package dao

import (
	"context"
	"fmt"
	"testing"
	"time"

	"go-common/app/admin/main/vip/model"
	xtime "go-common/library/time"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoDialogSave(t *testing.T) {
	curr := xtime.Time(time.Now().Unix())
	convey.Convey("DialogSave", t, func() {
		dlg := &model.ConfDialog{ID: 1, StartTime: curr, Operator: "tommy"}
		eff, err := d.DialogSave(context.TODO(), dlg)
		convey.So(err, convey.ShouldBeNil)
		convey.So(eff, convey.ShouldNotBeNil)
	})
	convey.Convey("DialogByID", t, func() {
		dlg, err := d.DialogByID(context.TODO(), 1)
		convey.So(err, convey.ShouldBeNil)
		convey.So(dlg, convey.ShouldNotBeNil)
	})
	convey.Convey("DialogEnable", t, func() {
		dlg := &model.ConfDialog{ID: 1, Stage: false}
		eff, err := d.DialogEnable(context.TODO(), dlg)
		convey.So(err, convey.ShouldBeNil)
		convey.So(eff, convey.ShouldNotBeNil)
	})
	convey.Convey("DialogEnable", t, func() {
		dlg := &model.ConfDialog{ID: 1, Stage: true}
		eff, err := d.DialogEnable(context.TODO(), dlg)
		convey.So(err, convey.ShouldBeNil)
		convey.So(eff, convey.ShouldNotBeNil)
	})
	convey.Convey("DialogAll", t, func() {
		res, err := d.DialogAll(context.TODO(), 0, 0, "active")
		convey.So(err, convey.ShouldBeNil)
		convey.So(res, convey.ShouldNotBeNil)
	})
}

func TestDaoCountDialogByPlatID(t *testing.T) {
	convey.Convey("TestDaoCountDialogByPlatID", t, func() {
		res, err := d.CountDialogByPlatID(context.TODO(), 1)
		fmt.Println(res)
		convey.So(err, convey.ShouldBeNil)
		convey.So(res, convey.ShouldBeGreaterThanOrEqualTo, 0)
	})
}
