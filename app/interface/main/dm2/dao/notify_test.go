package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaonotifyURI(t *testing.T) {
	convey.Convey("notifyURI", t, func(ctx convey.C) {
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			p1 := testDao.notifyURI()
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoSendNotify(t *testing.T) {
	convey.Convey("SendNotify", t, func(ctx convey.C) {
		var (
			c        = context.Background()
			title    = ""
			content  = ""
			dataType = ""
			mids     = []int64{}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			testDao.SendNotify(c, title, content, dataType, mids)
		})
	})
}
