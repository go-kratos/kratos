package service

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestServiceCheckBirthday(t *testing.T) {
	convey.Convey("CheckBirthday", t, func() {
		ok := s.CheckBirthday(context.Background(), 0)
		convey.So(ok, convey.ShouldNotBeNil)
	})
}
func TestServiceaccInfo(t *testing.T) {
	convey.Convey("accInfo", t, func() {
		ai, err := s.accInfo(context.Background(), 0)
		convey.So(err, convey.ShouldBeNil)
		convey.So(ai, convey.ShouldNotBeNil)
	})
}
