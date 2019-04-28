package block

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoSendSysMsg(t *testing.T) {
	convey.Convey("SendSysMsg", t, func(ctx convey.C) {
		var (
			c        = context.Background()
			mids     = []int64{1, 2, 3}
			content  = "账号违规处理通知-test-content"
			remoteIP = "127.0.0.1"
			code     = "2_3_2"
			title    = "账号违规处理通知-test"
		)
		ctx.Convey("When everything right", func(ctx convey.C) {
			err := d.SendSysMsg(c, code, mids, title, content, remoteIP)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaomidsToParam(t *testing.T) {
	convey.Convey("midsToParam", t, func(ctx convey.C) {
		var (
			mids = []int64{46333, 35858}
		)
		ctx.Convey("When everything right", func(ctx convey.C) {
			str := midsToParam(mids)
			ctx.Convey("Then str should equal mids[0],mids[1],....", func(ctx convey.C) {
				ctx.So(str, convey.ShouldEqual, "46333,35858")
			})
		})
	})
}
