package vip

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestVipCodeVerify(t *testing.T) {
	convey.Convey("CodeVerify", t, func(convCtx convey.C) {
		var (
			c = context.Background()
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			token, err := d.CodeVerify(c)
			convCtx.Convey("Then err should be nil.token should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(token, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestVipCodeOpen(t *testing.T) {
	convey.Convey("CodeOpen", t, func(convCtx convey.C) {
		var (
			c      = context.Background()
			mid    = int64(0)
			code   = ""
			token  = ""
			verify = ""
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			data, err := d.CodeOpen(c, mid, code, token, verify)
			convCtx.Convey("Then err should be nil.data should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(data, convey.ShouldNotBeNil)
			})
		})
	})
}
