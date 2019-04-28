package service

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestUpFilters(t *testing.T) {
	Convey("test update user rule", t, func() {
		rs, _, err := svr.UpFilters(context.TODO(), 27515615, 1, 1, 20)
		So(err, ShouldBeNil)
		So(rs, ShouldNotBeEmpty)
	})
}

func TestEditUpFilters(t *testing.T) {
	Convey("test edit user rule", t, func() {
		_, err := svr.EditUpFilters(context.TODO(), 66, 27515256, 1, 0)
		So(err, ShouldBeNil)
	})
}
