package service

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_UpdateUserNoticeState(t *testing.T) {
	Convey("update state", t, func() {
		err := s.UpdateUserNoticeState(context.TODO(), 100, "lead")
		So(err, ShouldBeNil)
		Convey("get lead data", func() {
			res, err := s.UserNoticeState(context.TODO(), 100)
			So(err, ShouldBeNil)
			So(res["lead"], ShouldBeTrue)
		})
		Convey("update new and get lead data", func() {
			err := s.UpdateUserNoticeState(context.TODO(), 100, "new")
			res, err := s.UserNoticeState(context.TODO(), 100)
			So(err, ShouldBeNil)
			So(res["lead"], ShouldBeTrue)
			So(res["new"], ShouldBeTrue)
		})
	})
	Convey("update invalid state", t, func() {
		err := s.UpdateUserNoticeState(context.TODO(), 100, "invalid")
		So(err, ShouldNotBeNil)
	})
}
