package service

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestService_Business(t *testing.T) {
	c := context.TODO()
	Convey("test service business", c, func() {
		Convey("test service add business", c, func() {
			_, err := testSvc.AddBusiness(c, -3, "abc", "abc", "abc", "abc")
			So(err, ShouldBeNil)
		})
		Convey("test service get business", c, func() {
			_, err := testSvc.GetBusiness(c, -3)
			So(err, ShouldBeNil)
		})
		Convey("test service list business", c, func() {
			_, err := testSvc.ListBusiness(c, 0)
			So(err, ShouldBeNil)
		})
		Convey("test service update business", c, func() {
			_, err := testSvc.UpBusiness(c, "test", "test", "test", "abc", -3)
			So(err, ShouldBeNil)
		})
		Convey("test service update business state", c, func() {
			_, err := testSvc.UpBusinessState(c, 1, -3)
			So(err, ShouldBeNil)
		})
	})
}
