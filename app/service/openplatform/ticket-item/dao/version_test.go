package dao

import (
	"context"
	"testing"

	"go-common/app/service/openplatform/ticket-item/model"

	. "github.com/smartystreets/goconvey/convey"
)

// TestAddVersion
func TestDao_AddVersion(t *testing.T) {
	Convey("AddVersion", t, func() {
		once.Do(startService)
		err := d.AddVersion(context.TODO(), nil, &model.Version{
			Type:       2,
			Status:     1, // 审核中
			ItemName:   "gotest",
			ParentID:   10164,
			TargetItem: 0,
			AutoPub:    1, // 自动上架
		}, &model.VersionExt{
			Type:     1,
			MainInfo: "{'name':'公告test','introduction':'公告简介','content':'公告内容','pid':10164,'project_name':'删通票删票种'}",
		})
		So(err, ShouldBeNil)
	})
}

// TestUpdateVersion
func TestDao_UpdateVersion(t *testing.T) {
	Convey("UpdateVersion", t, func() {
		once.Do(startService)
		res, err := d.UpdateVersion(context.TODO(), &model.Version{
			VerID:      2691387070776769288,
			Type:       2,
			Status:     2, // 审核中
			ItemName:   "gotest公告",
			ParentID:   0,
			TargetItem: 10164,
			AutoPub:    1, // 自动上架
		})
		So(res, ShouldBeTrue)
		So(err, ShouldBeNil)
	})
}

// TestGetVersion
func TestDao_GetVersion(t *testing.T) {
	Convey("GetVersion", t, func() {
		once.Do(startService)
		verInfo, verExtInfo, err := d.GetVersion(context.TODO(), 153008633987459678, true)
		So(verInfo, ShouldNotBeNil)
		So(verExtInfo, ShouldNotBeNil)
		So(err, ShouldBeNil)
	})
}

// TestRejectVersion
func TestDao_RejectVersion(t *testing.T) {
	Convey("RejectVersion", t, func() {
		once.Do(startService)
		res, err := d.RejectVersion(context.TODO(), 2691387070776769288, 2)
		So(res, ShouldBeTrue)
		So(err, ShouldBeNil)
	})
}

// TestAddVersionLog
func TestDao_AddVersionLog(t *testing.T) {
	Convey("AddVersionLog", t, func() {
		once.Do(startService)
		err := d.AddVersionLog(context.TODO(), &model.VersionLog{
			VerID:  2691387070776769288,
			Type:   1,
			Log:    "reject",
			IsPass: 0,
			Uname:  "tester",
		})
		So(err, ShouldBeNil)
	})
}
