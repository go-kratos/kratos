package dao

import (
	"context"
	"testing"

	"go-common/library/ecode"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoCaptchaToken(t *testing.T) {
	convey.Convey("CaptchaToken", t, func(convCtx convey.C) {
		var (
			c   = context.Background()
			bid = d.c.Property.CaptchaBID
			ip  = "127.0.0.1"
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			res, err := d.CaptchaToken(c, bid, ip)
			convCtx.Convey("Then err should be nil.res should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoCaptchaVerify(t *testing.T) {
	convey.Convey("CaptchaVerify", t, func(convCtx convey.C) {
		var (
			c     = context.Background()
			code  = "xxxx"
			token = "xxxx"
			ip    = ""
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			err := d.CaptchaVerify(c, code, token, ip)
			if err == ecode.CouponCodeVerifyFaildErr {
				err = nil
			}
			convCtx.So(err, convey.ShouldBeNil)
		})
	})
}
