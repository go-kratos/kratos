package view

import (
	"context"
	"testing"

	"go-common/app/service/main/archive/api"
	"go-common/app/service/main/archive/model/archive"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_UpArcCache(t *testing.T) {
	Convey("UpArcCache", t, func() {
		err := d.UpArcCache(context.TODO(), &archive.Archive3{Aid: 0})
		So(err, ShouldBeNil)
	})
}

func Test_DelArcCache(t *testing.T) {
	Convey("DelArcCache", t, func() {
		err := d.DelArcCache(context.TODO(), 0)
		So(err, ShouldBeNil)
	})
}

func Test_UpViewCache(t *testing.T) {
	Convey("UpViewCache", t, func() {
		err := d.UpViewCache(context.TODO(), &archive.View3{Archive3: &archive.Archive3{}})
		So(err, ShouldBeNil)
	})
}

func Test_DelViewCache(t *testing.T) {
	Convey("DelViewCache", t, func() {
		err := d.DelViewCache(context.TODO(), 0)
		So(err, ShouldBeNil)
	})
}

func Test_UpStatCache(t *testing.T) {
	Convey("UpStatCache", t, func() {
		err := d.UpStatCache(context.TODO(), &api.Stat{})
		So(err, ShouldBeNil)
	})
}

func Test_UpViewContributeCache(t *testing.T) {
	Convey("UpViewContributeCache", t, func() {
		err := d.UpViewContributeCache(context.TODO(), 1, []int64{})
		So(err, ShouldBeNil)
	})
}
