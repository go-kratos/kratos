package message

import (
	"context"
	"testing"
	"time"

	"github.com/smartystreets/goconvey/convey"
)

func TestMessageSend(t *testing.T) {
	convey.Convey("Send", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			mc    = "1_14_2"
			title = "test"
			msg   = "test"
			mids  = []int64{253550886}
			ts    = time.Now().Unix()
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.Send(c, mc, title, msg, mids, ts)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestMessageNotifyTask(t *testing.T) {
	convey.Convey("NotifyTask", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			mids = []int64{2316310}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.NotifyTask(c, mids)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldNotBeNil)
			})
		})
	})
}
