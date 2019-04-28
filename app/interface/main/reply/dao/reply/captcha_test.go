package reply

import (
	"context"
	httpx "go-common/library/net/http/blademaster"
	xtime "go-common/library/time"
	"testing"
	"time"

	"github.com/smartystreets/goconvey/convey"
)

func TestReplyNewCaptchaDao(t *testing.T) {
	convey.Convey("NewCaptchaDao", t, func(ctx convey.C) {
		var (
			c = &httpx.ClientConfig{
				App: &httpx.App{
					Key:    "test",
					Secret: "test",
				},
				Dial:    xtime.Duration(time.Second),
				Timeout: xtime.Duration(time.Second),
			}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := NewCaptchaDao(c)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestReplyCaptcha(t *testing.T) {
	convey.Convey("Captcha", t, func(ctx convey.C) {
		var (
			c = context.Background()
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1, p2, err := d.Captcha.Captcha(c)
			ctx.Convey("Then err should be nil.p1,p2 should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(p2, convey.ShouldNotBeNil)
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestReplyVerify(t *testing.T) {
	convey.Convey("Verify", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			token = ""
			code  = ""
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.Captcha.Verify(c, token, code)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldNotBeNil)
			})
		})
	})
}
