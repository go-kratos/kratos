package dao

import (
	"context"
	"fmt"
	"testing"

	"go-common/app/service/main/passport-sns/model"

	"github.com/smartystreets/goconvey/convey"
)

func TestDao_WeiboAuthorize(t *testing.T) {
	var (
		c           = context.Background()
		AppID       = "101135748"
		RedirectUrl = "https://passport.bilibili.com/login/snsback?sns=weibo"
		Display     = "mobile"
		// AppID       : "1108092926",
		// RedirectUrl : "https://passport.bilibili.com/web/sns/bind/callback/weibo",
	)
	convey.Convey("WeiboAuthorize", t, func(ctx convey.C) {
		res := d.WeiboAuthorize(c, AppID, RedirectUrl, Display)
		ctx.Convey("Then res should not be nil.", func(ctx convey.C) {
			ctx.So(res, convey.ShouldNotBeNil)
		})
		fmt.Println(res)
	})
}

func TestDao_WeiboOauth2Info(t *testing.T) {
	var (
		c           = context.Background()
		code        = "C4946CD493AEEDE67C574DFE2C756D09"
		redirectUrl = "https://passport.bilibili.com/web/sns/bind/callback"
		app         = &model.SnsApps{
			AppID:     "",
			AppSecret: "",
			Business:  model.BusinessMall,
		}
	)
	convey.Convey("WeiboOauth2Info", t, func(ctx convey.C) {
		res, err := d.WeiboOauth2Info(c, code, redirectUrl, app)
		ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldNotBeNil)
			ctx.So(res, convey.ShouldBeNil)
		})
		fmt.Printf("(%+v) error(%+v)", res, err)
	})
}

func TestDao_weiboAccessToken(t *testing.T) {
	var (
		c           = context.Background()
		code        = "CF8CE1408E8E43E4CD2DC778B5993FBB"
		appID       = ""
		appSecret   = ""
		redirectUrl = "https://passport.bilibili.com/web/sns/bind/callback"
	)
	convey.Convey("weiboAccessToken", t, func(ctx convey.C) {
		res, err := d.weiboAccessToken(c, code, appID, appSecret, redirectUrl)
		ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldNotBeNil)
			ctx.So(res, convey.ShouldBeNil)
		})
		fmt.Printf("(%+v) error(%+v)", res, err)
	})
}
