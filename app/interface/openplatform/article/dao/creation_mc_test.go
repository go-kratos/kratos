package dao

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_SubmitCache(t *testing.T) {
	c := context.TODO()
	Convey("add cache", t, func() {
		err := d.AddSubmitCache(c, 100, "title")
		So(err, ShouldBeNil)
		Convey("get cache should work", func() {
			res, err1 := d.SubmitCache(c, 100, "title")
			So(err1, ShouldBeNil)
			So(res, ShouldBeTrue)
			res, err1 = d.SubmitCache(c, 200, "title")
			So(err1, ShouldBeNil)
			So(res, ShouldBeFalse)
		})
		Convey("delete cache should not present", func() {
			err = d.DelSubmitCache(c, 100, "title")
			So(err, ShouldBeNil)
			err = d.DelSubmitCache(c, 200, "title")
			So(err, ShouldBeNil)
		})
	})
}
