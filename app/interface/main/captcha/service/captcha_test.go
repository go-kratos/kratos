package service

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	bid   = "account"
	token = "5049a45ffc7c49489c14a7677c4548e2"
)

func TestToken(t *testing.T) {
	var (
		c = context.Background()
	)
	Convey("err should return nil", t, func() {
		_, t, err := svr.Token(c, bid)
		So(err, ShouldBeNil)
		So(t, ShouldNotBeNil)
	})
	Convey("err should return nil", t, func() {
		err := svr.VerifyCaptcha(c, token, "test")
		So(err, ShouldNotBeNil)
	})
	Convey("err should return nil", t, func() {
		business := svr.LookUp(bid)
		So(business, ShouldNotBeEmpty)
	})
}
