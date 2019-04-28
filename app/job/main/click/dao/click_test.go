package dao

import (
	"context"
	"testing"

	"go-common/app/job/main/click/model"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_Click(t *testing.T) {
	Convey("Click", t, func() {
		_, err := d.Click(context.TODO(), 1)
		So(err, ShouldBeNil)
	})
}

func Test_AddClick(t *testing.T) {
	Convey("AddClick", t, func() {
		_, err := d.AddClick(context.TODO(), 3, 1, 1, 1, 1, 1, 100)
		So(err, ShouldBeNil)
	})
}

func Test_UpClick(t *testing.T) {
	Convey("UpClick", t, func() {
		rows, err := d.UpClick(context.TODO(), &model.ClickInfo{Aid: 2, AndroidTV: 22222})
		Println(rows, err)
	})
}

func Test_UpSpecial(t *testing.T) {
	Convey("UpSpecial", t, func() {
		_, err := d.UpSpecial(context.TODO(), 1, model.TypeForAndroid, 1)
		So(err, ShouldBeNil)
	})
}
