package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaonotifyURI(t *testing.T) {
	convey.Convey("notifyURI", t, func(convCtx convey.C) {
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			p1 := testDao.notifyURI()
			convCtx.Convey("Then p1 should not be nil.", func(convCtx convey.C) {
				convCtx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoSendNotify(t *testing.T) {
	convey.Convey("SendNotify", t, func(convCtx convey.C) {
		var (
			c       = context.Background()
			title   = ""
			content = ""
			mids    = []int64{1}
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			err := testDao.SendNotify(c, title, content, mids)
			convCtx.Convey("Then err should be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
			})
		})
	})
}
