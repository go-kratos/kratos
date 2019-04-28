package service

import (
	"testing"

	"github.com/smartystreets/goconvey/convey"
	"go-common/app/service/main/videoup/model/archive"
)

func TestServicesendMsg(t *testing.T) {
	convey.Convey("sendMsg", t, func(ctx convey.C) {
		var (
			arg = &archive.ArgMsg{}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			svr.sendMsg(arg)
			ctx.Convey("No return values", func(ctx convey.C) {
			})
		})
	})
}
