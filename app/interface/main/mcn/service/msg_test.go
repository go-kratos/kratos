package service

import (
	"testing"

	"go-common/app/admin/main/mcn/model"

	"github.com/smartystreets/goconvey/convey"
)

func TestServicesendMsg(t *testing.T) {
	convey.Convey("sendMsg", t, func(ctx convey.C) {
		var (
			arg = &model.ArgMsg{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			s.sendMsg(arg)
			ctx.Convey("No return values", func(ctx convey.C) {
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
