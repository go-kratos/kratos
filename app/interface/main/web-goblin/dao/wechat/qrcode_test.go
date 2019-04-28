package wechat

import (
	"context"
	"testing"

	"go-common/app/interface/main/web-goblin/model/wechat"
	"go-common/library/ecode"

	"github.com/smartystreets/goconvey/convey"
)

func TestWechatQrcode(t *testing.T) {
	convey.Convey("Qrcode", t, func(ctx convey.C) {
		var (
			c           = context.Background()
			accessToken = "14_LZVbKTtstzal_T-AfG-EgkUI2WlCdRvKUqhiYKMhNyxsGjzc1K_a1GGWuMPCbX"
			arg         = `{"page":"","scene":"?avid=34188644"}`
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			qrcode, err := d.Qrcode(c, accessToken, arg)
			ctx.Convey("Then err should be nil.qrcode should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldEqual, ecode.RequestErr)
				ctx.Println(qrcode)
			})
		})
	})
}

func TestWechatAddCacheAccessToken(t *testing.T) {
	convey.Convey("AddCacheAccessToken", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			data = &wechat.AccessToken{AccessToken: "string", ExpiresIn: 1111}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.AddCacheAccessToken(c, data)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestWechatCacheAccessToken(t *testing.T) {
	convey.Convey("CacheAccessToken", t, func(ctx convey.C) {
		var (
			c = context.Background()
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			data, err := d.CacheAccessToken(c)
			ctx.Convey("Then err should be nil.data should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(data, convey.ShouldNotBeNil)
			})
		})
	})
}
