package bnj

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"

	"gopkg.in/h2non/gock.v1"
)

func TestBnjSendWechat(t *testing.T) {
	convey.Convey("SendWechat", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			title = "【拜年祭必看！】拜年祭预约人数到达预警"
			msg   = "拜年祭预约人数即将到达50w，请及时准备拜年祭抽奖事项。"
			user  = "wuhao02"
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			defer gock.OffAll()
			httpMock("POST", _wechatURL).Reply(200).JSON(`{"RetCode":0}`)
			err := d.SendWechat(c, title, msg, user)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}
