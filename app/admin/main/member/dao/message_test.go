package dao

import (
	"context"
	"gopkg.in/h2non/gock.v1"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoMessage(t *testing.T) {
	convey.Convey("Message", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			title = "abc test"
			msg   = "abc"
			mids  = []int64{112}
			mc    = "2_2_2"
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			defer gock.OffAll()
			httpMock("POST", d.msgURL).Reply(200).JSON(`{"code":0}`)
			err := d.RawMessage(c, mc, title, msg, mids)
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoRawMessage(t *testing.T) {
	convey.Convey("RawMessage", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			mc    = "2_2_2"
			title = "abc test"
			msg   = "abc test"
			mids  = []int64{112}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			defer gock.OffAll()
			httpMock("POST", d.msgURL).Reply(200).JSON(`{"code":0}`)
			err := d.RawMessage(c, mc, title, msg, mids)
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}
