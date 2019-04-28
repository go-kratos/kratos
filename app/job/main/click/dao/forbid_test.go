package dao

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_Forbids(t *testing.T) {
	Convey("Forbids", t, func() {
		fs, err := d.Forbids(context.TODO())
		So(err, ShouldBeNil)
		Println(fs)
	})
}

func Test_UpForbid(t *testing.T) {
	Convey("UpForbid", t, func() {
		_, err := d.UpForbid(context.TODO(), 1, 1, 1, 1)
		So(err, ShouldBeNil)
	})
}

func Test_ForbidMids(t *testing.T) {
	Convey("ForbidMids", t, func() {
		mids, err := d.ForbidMids(context.TODO())
		So(err, ShouldBeNil)
		Println(mids)
	})
}

func Test_UpForbidMid(t *testing.T) {
	Convey("UpForbidMid", t, func() {
		err := d.UpMidForbidStatus(context.TODO(), 1684013, 0)
		So(err, ShouldBeNil)
	})
}
