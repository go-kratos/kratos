package vip

import (
	"net/url"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestVipsign(t *testing.T) {
	convey.Convey("sign", t, func(convCtx convey.C) {
		var (
			params url.Values
			path   = ""
			secret = ""
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			p1 := sign(params, path, secret)
			convCtx.Convey("Then p1 should not be nil.", func(convCtx convey.C) {
				convCtx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}
