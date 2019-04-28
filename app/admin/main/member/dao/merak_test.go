package dao

import (
	"context"
	"strings"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoSendWechat(t *testing.T) {
	convey.Convey("SendWechat", t, func(ctx convey.C) {
		var (
			content = "测试内容"
			title   = "测试标题"
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.MerakNotify(context.Background(), content, title)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoGetMerakSign(t *testing.T) {
	convey.Convey("sign", t, func(ctx convey.C) {
		var (
			params = map[string]string{
				"Action":    "CreateWechatMessage",
				"PublicKey": _publicKey,
				//"UserName":  strings.Join(d.c.ReviewNotify.Users, ","),
				"UserName": strings.Join([]string{"user1", "user2"}, ","),
				"Title":    "测试标题",
				"Content":  "测试内容",
				"TreeId":   "",
			}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1, err := MerakSign(params)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldEqual, "59cd4e74b225a7d326ee7d6c89bf27cf2f6015dc")
			})
		})
	})
}
