package service

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestVip(t *testing.T) {
	convey.Convey("Vip", t, func() {
		res, err := s.Vip(context.TODO(), 1)
		convey.So(res, convey.ShouldNotBeNil)
		convey.So(err, convey.ShouldBeNil)
	})
}

func TestVips(t *testing.T) {
	convey.Convey("Vips", t, func() {
		res, err := s.Vips(context.TODO(), []int64{1, 2, 3})
		convey.So(res, convey.ShouldNotBeNil)
		convey.So(err, convey.ShouldBeNil)
	})
}
