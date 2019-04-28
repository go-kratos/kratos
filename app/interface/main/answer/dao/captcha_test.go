package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestCaptcha(t *testing.T) {
	convey.Convey("Captcha", t, func(ctx convey.C) {
		var (
			c = context.Background()
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1, p2, err := d.Captcha(c)
			ctx.Convey("Then err should be nil.p1,p2 should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(p2, convey.ShouldNotBeNil)
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestVerify(t *testing.T) {
	convey.Convey("Verify", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			token = ""
			code  = ""
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.Verify(c, token, code, "")
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldNotBeNil)
			})
		})
	})
}
