package dao

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	token = "5049a45ffc7c49489c14a7677c4548e2"
	ttl   = 150
)

func TestAddTokenCache(t *testing.T) {
	var (
		c = context.TODO()
	)
	Convey("err should return nil, and ttl not -1", t, func() {
		err := d.AddTokenCache(c, token, int32(ttl))
		So(err, ShouldBeNil)
	})
}

func TestDelCaptchaCache(t *testing.T) {
	var (
		c = context.TODO()
	)
	Convey("err should return nil", t, func() {
		err := d.DelCaptchaCache(c, token)
		So(err, ShouldBeNil)
	})
}

func TestCaptchaCache(t *testing.T) {
	var (
		c = context.TODO()
	)
	Convey("err should return nil", t, func() {
		_, _, err := d.CaptchaCache(c, token)
		So(err, ShouldBeNil)
	})
}
