package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
	gock "gopkg.in/h2non/gock.v1"
)

func TestDaoSendMessage(t *testing.T) {
	convey.Convey("SendMessage", t, func(convCtx convey.C) {
		var (
			c       = context.Background()
			mids    = ""
			title   = ""
			content = ""
			ip      = ""
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			defer gock.OffAll()
			httpMock("POST", sendMessage).Reply(200).JSON(`{"code":0}`)
			err := d.SendMessage(c, mids, title, content, ip)
			convCtx.Convey("Then err should be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
			})
		})
	})
}
