package dao

import (
	"context"
	"go-common/library/ecode"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestVerifyDsb(t *testing.T) {
	convey.Convey("VerifyDsb", t, func() {
		sid := "1234567890"
		_, err := d.VerifyDsb(context.Background(), sid)
		if err == ecode.NoLogin {
			err = nil
		}
		convey.So(err, convey.ShouldBeNil)
	})
}

func TestNewSession(t *testing.T) {
	convey.Convey("NewSession", t, func() {
		res := d.NewSession(context.Background())
		convey.So(res, convey.ShouldNotBeNil)
	})
}
