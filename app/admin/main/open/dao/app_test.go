package dao

import (
	"context"
	"testing"

	"go-common/app/admin/main/open/model"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoAddApp(t *testing.T) {
	convey.Convey("AddApp", t, func() {
		p1 := d.AddApp(context.TODO(), nil)
		convey.So(p1, convey.ShouldBeNil)
	})
}

func TestDaoDelApp(t *testing.T) {
	convey.Convey("DelApp", t, func() {
		p1 := d.DelApp(context.TODO(), 0)
		convey.So(p1, convey.ShouldBeNil)
	})
}

func TestDaoUpdateApp(t *testing.T) {
	bean := &model.AppParams{
		AppID:   123,
		AppName: "xxx",
	}
	convey.Convey("UpdateApp", t, func() {
		p1 := d.UpdateApp(context.TODO(), bean)
		convey.So(p1, convey.ShouldBeNil)
	})
}

func TestDaoListApp(t *testing.T) {
	convey.Convey("ListApp", t, func() {
		res, err := d.ListApp(context.TODO(), nil)
		convey.So(err, convey.ShouldBeNil)
		convey.So(res, convey.ShouldNotBeNil)
	})
}
