package dao

import (
	"context"
	"testing"

	"go-common/app/admin/main/workflow/model/param"
	"go-common/app/job/main/workflow/model"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoSendMessage(t *testing.T) {
	convey.Convey("SendMessage", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			chs = []*model.ChallRes{}
			msg = &param.MessageParam{}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := d.SendMessage(c, chs, msg)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err.Error(), convey.ShouldEqual, "-6")
			})
		})
	})
}
