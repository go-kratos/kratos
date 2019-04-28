package dao

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_Managers(t *testing.T) {
	Convey("return someting", t, func() {
		res, err := d.Managers(context.Background())
		So(err, ShouldBeNil)
		So(res, ShouldNotBeNil)
	})
}

func Test_Manager(t *testing.T) {
	Convey("return someting", t, func() {
		res, err := d.Manager(context.Background(), 1, 2)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeNil)
	})
}

func Test_ManagerTotal(t *testing.T) {
	Convey("return someting", t, func() {
		res, err := d.ManagerTotal(context.Background())
		So(err, ShouldBeNil)
		So(res, ShouldNotBeNil)
	})
}
