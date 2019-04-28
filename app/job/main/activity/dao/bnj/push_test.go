package bnj

import (
	"context"
	"testing"

	"gopkg.in/h2non/gock.v1"

	"github.com/smartystreets/goconvey/convey"
)

func TestBnjPushAll(t *testing.T) {
	convey.Convey("PushAll", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			msg = `{"second":100,"name":"啊*"}`
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := d.PushAll(c, msg)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestBnjSendMessage(t *testing.T) {
	convey.Convey("SendMessage", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			mids  = []int64{2089809}
			mc    = "1_21_1"
			title = "【bilibili2019拜年祭档案揭秘】001"
			msg   = "飞雪连天射白鹿，笑书神侠倚碧鸳。当V家碰到金庸，会碰撞出怎样的火花？来拜年祭后台看看吧~"
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			defer gock.OffAll()
			httpMock("POST", d.messageURL).Reply(200).JSON(`{"code":0}`)
			err := d.SendMessage(c, mids, mc, title, msg)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}
