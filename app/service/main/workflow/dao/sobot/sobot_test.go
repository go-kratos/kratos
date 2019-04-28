package sobot

import (
	"context"
	"testing"

	"go-common/app/service/main/workflow/model/sobot"

	"github.com/smartystreets/goconvey/convey"
)

func TestSobotSobotTicketInfo(t *testing.T) {
	var (
		c        = context.Background()
		ticketID = int32(0)
	)
	convey.Convey("SobotTicketInfo", t, func(ctx convey.C) {
		res, err := d.SobotTicketInfo(c, ticketID)
		ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(res, convey.ShouldNotBeNil)
		})
	})
}

func TestSobotSobotAddTicket(t *testing.T) {
	var (
		c  = context.Background()
		tp = &sobot.TicketParam{
			TicketTitle:   "我是202工单",
			TicketID:      202,
			TicketContent: "233333333",
			CustomerEmail: "1107691251@qq.com",
		}
	)
	convey.Convey("SobotAddTicket", t, func(ctx convey.C) {
		err := d.SobotAddTicket(c, tp)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err.Error(), convey.ShouldEqual, "2000101")
		})
	})
}

func TestSobotSobotAddReply(t *testing.T) {
	var (
		c  = context.Background()
		rp = &sobot.ReplyParam{
			TicketID:      202,
			ReplyContent:  "reply_test",
			CustomerEmail: "1107691251@qq.com",
			StartType:     1,
			ReplyType:     1,
		}
	)
	convey.Convey("SobotAddReply", t, func(ctx convey.C) {
		err := d.SobotAddReply(c, rp)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestSobotSobotTicketModify(t *testing.T) {
	var (
		c  = context.Background()
		tp = &sobot.TicketParam{
			TicketID:      202,
			CustomerEmail: "1107691251@qq.com",
			StartType:     1,
		}
	)
	convey.Convey("SobotTicketModify", t, func(ctx convey.C) {
		err := d.SobotTicketModify(c, tp)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}
