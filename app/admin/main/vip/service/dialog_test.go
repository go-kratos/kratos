package service

import (
	"context"
	"go-common/app/admin/main/vip/model"
	"go-common/library/ecode"
	xtime "go-common/library/time"
	"testing"
	"time"

	"github.com/smartystreets/goconvey/convey"
)

func TestServiceDialogSave(t *testing.T) {
	curr := xtime.Time(time.Now().Unix())
	convey.Convey("DialogSave", t, func() {
		dlg1 := &model.ConfDialog{ID: 1, StartTime: curr, Operator: "tommy"}
		eff, err := s.DialogSave(context.TODO(), dlg1)
		convey.So(err, convey.ShouldBeNil)
		convey.So(eff, convey.ShouldNotBeNil)
	})
	convey.Convey("DialogSave2", t, func() {
		dlg2 := &model.ConfDialog{ID: 2, StartTime: xtime.Time(time.Now().Unix() + 1), Operator: "tommy"}
		_, err := s.DialogSave(context.TODO(), dlg2)
		convey.So(err, convey.ShouldEqual, ecode.VipDialogConflictErr)
	})
	convey.Convey("DialogByID", t, func() {
		dlg, err := s.DialogByID(context.TODO(), &model.ArgID{ID: 1})
		convey.So(err, convey.ShouldBeNil)
		convey.So(dlg, convey.ShouldNotBeNil)
	})
	convey.Convey("DialogEnable", t, func() {
		eff, err := s.DialogEnable(context.TODO(), &model.ConfDialog{ID: 1, Stage: true})
		convey.So(err, convey.ShouldBeNil)
		convey.So(eff, convey.ShouldNotBeNil)
	})
	convey.Convey("DialogAll", t, func() {
		res, err := s.DialogAll(context.TODO(), 0, 0, "active")
		convey.So(err, convey.ShouldBeNil)
		convey.So(res, convey.ShouldNotBeNil)
	})
}

func TestServiceDialogDel(t *testing.T) {
	convey.Convey("DialogDel", t, func() {
		eff, err := s.DialogDel(context.TODO(), nil, "")
		convey.So(err, convey.ShouldBeNil)
		convey.So(eff, convey.ShouldNotBeNil)
	})
}

func TestServiceDialogStatus(t *testing.T) {
	convey.Convey("dialogStatus padding", t, func() {
		v := &model.ConfDialog{Stage: true, StartTime: xtime.Time(time.Now().AddDate(1, 1, 1).Unix())}
		res := dialogStatus(v)
		convey.So(res, convey.ShouldEqual, "padding")
	})
	convey.Convey("dialogStatus active", t, func() {
		v := &model.ConfDialog{Stage: true}
		res := dialogStatus(v)
		convey.So(res, convey.ShouldEqual, "active")
		v.StartTime = xtime.Time(time.Now().AddDate(-1, 1, 1).Unix())
		res = dialogStatus(v)
		convey.So(res, convey.ShouldEqual, "active")
		v.EndTime = xtime.Time(time.Now().AddDate(1, 1, 1).Unix())
		res = dialogStatus(v)
		convey.So(res, convey.ShouldEqual, "active")
	})
	convey.Convey("dialogStatus inactive", t, func() {
		v := &model.ConfDialog{Stage: false}
		res := dialogStatus(v)
		convey.So(res, convey.ShouldEqual, "inactive")
		v.EndTime = xtime.Time(time.Now().AddDate(-1, 1, 1).Unix())
		res = dialogStatus(v)
		convey.So(res, convey.ShouldEqual, "inactive")
	})

}
