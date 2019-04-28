package dao

import (
	"context"
	"testing"

	"go-common/app/admin/main/workflow/model"
	"go-common/app/admin/main/workflow/model/param"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoSendMessage(t *testing.T) {
	convey.Convey("SendMessage", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			msg = &param.MessageParam{
				Type:     "json",
				Source:   1,
				DataType: 4,
				MC:       model.WkfNotifyMC,
				Title:    "test title",
				Context:  "test context",
				MidList:  []int64{1},
			}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.SendMessage(c, msg)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err.Error(), convey.ShouldEqual, "-6")
			})
		})
	})
}
