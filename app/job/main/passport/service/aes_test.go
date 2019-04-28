package service

import (
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestService_encrypt(t *testing.T) {
	once.Do(startService)
	convey.Convey("", t, func() {
		text := "123456"
		et, err := s.encrypt(text)
		convey.So(err, convey.ShouldBeNil)
		dt, err := s.decrypt(et)
		convey.So(err, convey.ShouldBeNil)
		convey.So(dt, convey.ShouldEqual, text)
	})
}
