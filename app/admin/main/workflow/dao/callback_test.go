package dao

import (
	"context"
	"net/url"
	"testing"

	"go-common/app/admin/main/workflow/model"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoAllCallbacks(t *testing.T) {
	convey.Convey("AllCallbacks", t, func(ctx convey.C) {
		var (
			c = context.Background()
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			cbs, err := d.AllCallbacks(c)
			ctx.Convey("Then err should be nil.cbs should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(cbs, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoSendCallback(t *testing.T) {
	convey.Convey("SendCallback", t, func(ctx convey.C) {
		var (
			c  = context.Background()
			cb = &model.Callback{
				URL: "http://uat-manager.bilibili.co/x/admin/reply/internal/callback/del",
			}
			payload = &model.Payload{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.SendCallback(c, cb, payload)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaosign(t *testing.T) {
	convey.Convey("sign", t, func(ctx convey.C) {
		var (
			params url.Values
			appkey = ""
			secret = ""
			lower  bool
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			hexdigest := sign(params, appkey, secret, lower)
			ctx.Convey("Then hexdigest should not be nil.", func(ctx convey.C) {
				ctx.So(hexdigest, convey.ShouldNotBeNil)
			})
		})
	})
}
