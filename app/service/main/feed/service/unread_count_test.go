package service

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_UnreadCount(t *testing.T) {
	Convey("app should return without err", t, WithService(t, func(svf *Service) {
		_, err := svf.UnreadCount(context.TODO(), true, false, _mid, _ip)
		So(err, ShouldBeNil)
	}))

	Convey("app without bangumi should return without err", t, WithService(t, func(svf *Service) {
		_, err := svf.UnreadCount(context.TODO(), true, true, _mid, _ip)
		So(err, ShouldBeNil)
	}))

	Convey("web should return without err", t, WithService(t, func(svf *Service) {
		_, err := svf.UnreadCount(context.TODO(), false, false, _mid, _ip)
		So(err, ShouldBeNil)
	}))
}
