package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
	gock "gopkg.in/h2non/gock.v1"
)

func TestDaoMutliSendSysMsg(t *testing.T) {
	convey.Convey("MutliSendSysMsg", t, func(ctx convey.C) {
		var (
			c       = context.Background()
			allMids = []int64{27515256}
			title   = "title"
			context = "context"
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			defer gock.OffAll()
			httpMock("GET", d.msgURL).Reply(200).JSON(`{"code":20007}`)
			err := d.MutliSendSysMsg(c, allMids, title, context)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoSendSysMsg(t *testing.T) {
	convey.Convey("SendSysMsg", t, func(ctx convey.C) {
		var (
			c       = context.Background()
			mids    = []int64{27515256}
			title   = "title"
			context = "context"
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			defer gock.OffAll()
			httpMock("GET", d.msgURL).Reply(200).JSON(`{"code":20007}`)
			err := d.SendSysMsg(c, mids, title, context)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldNotBeNil)
			})
		})
	})
}
