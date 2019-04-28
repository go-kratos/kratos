package service

import (
	"context"
	"testing"

	"go-common/app/admin/main/mcn/model"

	"github.com/smartystreets/goconvey/convey"
)

func TestServicesendMsg(t *testing.T) {
	convey.Convey("sendMsg", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			arg = &model.ArgMsg{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := s.sendMsg(c, arg)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestServicesetMsgTypeMap(t *testing.T) {
	convey.Convey("setMsgTypeMap", t, func(ctx convey.C) {
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			s.setMsgTypeMap()
			ctx.Convey("No return values", func(ctx convey.C) {
			})
		})
	})
}
