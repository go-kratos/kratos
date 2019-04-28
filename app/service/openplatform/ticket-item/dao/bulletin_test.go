package dao

import (
	"context"
	"testing"

	item "go-common/app/service/openplatform/ticket-item/api/grpc/v1"

	. "github.com/smartystreets/goconvey/convey"
)

// TestGetBulletins
func TestDao_GetBulletins(t *testing.T) {
	Convey("GetBulletins", t, func() {
		once.Do(startService)
		res, err := d.GetBulletins(context.TODO(), 72)
		So(res, ShouldNotBeNil)
		So(err, ShouldBeNil)
	})
}

// TestAddBulletin
func TestDao_AddBulletin(t *testing.T) {
	Convey("AddBulletin", t, func() {
		once.Do(startService)
		res, err := d.AddBulletin(context.TODO(), &item.BulletinInfoRequest{
			ParentID:   72,
			Title:      "go test bulletin",
			Content:    "goooo",
			Detail:     "goo",
			TargetItem: 0,
		})
		So(res, ShouldBeTrue)
		So(err, ShouldBeNil)
	})
}

// TestUpdateBulletin
func TestDao_UpdateBulletin(t *testing.T) {
	Convey("UpdateBulletin", t, func() {
		once.Do(startService)
		res, err := d.UpdateBulletin(context.TODO(), &item.BulletinInfoRequest{
			ParentID:   72,
			Title:      "go test bulletin22",
			Content:    "gooossso",
			Detail:     "goo",
			TargetItem: 0,
			VerID:      2692936350844594047,
		})
		So(res, ShouldBeTrue)
		So(err, ShouldBeNil)
	})
}

// TestPassBulletin
func TestDao_PassBulletin(t *testing.T) {
	Convey("PassBulletin", t, func() {
		once.Do(startService)
		res, err := d.PassBulletin(context.TODO(), 2692936350844594047)
		So(res, ShouldBeTrue)
		So(err, ShouldBeNil)
	})
}

// TestUnpublishBulletin
func TestDao_UnpublishBulletin(t *testing.T) {
	Convey("UnpublishBulletin", t, func() {
		once.Do(startService)
		res, err := d.UnpublishBulletin(context.TODO(), 2692936350844594047, -1)
		So(res, ShouldBeTrue)
		So(err, ShouldBeNil)
	})
}
