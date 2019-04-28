package dao

import (
	"context"
	"testing"

	gock "gopkg.in/h2non/gock.v1"

	"github.com/smartystreets/goconvey/convey"
)

// go test  -test.v -test.run TestDaoNotifyRet
func TestDaoNotifyRet(t *testing.T) {
	convey.Convey("NotifyRet", t, func(convCtx convey.C) {
		var (
			c         = context.Background()
			notifyURL = "http://bangumi.bilibili.com/pay/inner/notify_ticket"
			ticketNO  = "706821058124120180326154"
			orderNO   = "20180326153703795"
			ip        = "127.0.0.1"
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			data, err := d.NotifyRet(c, notifyURL, ticketNO, orderNO, ip)
			t.Logf("data (%v)", data)
			convCtx.Convey("Then err should be nil.data should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(data, convey.ShouldNotBeNil)
			})
			convCtx.Convey("res.Code!=0 Then err should not be nil.data should be nil.", func(convCtx convey.C) {
				defer gock.OffAll()
				httpMock("POST", notifyURL).Reply(200).JSON(`{"code":-1}`)
				data, err = d.NotifyRet(c, notifyURL, ticketNO, orderNO, ip)
				convCtx.So(err, convey.ShouldNotBeNil)
				convCtx.So(data, convey.ShouldBeNil)
			})
			convCtx.Convey("call service error Then err should not be nil.data should be nil.", func(convCtx convey.C) {
				defer gock.OffAll()
				httpMock("POST", "").Reply(-400).JSON(`{"code":-1}`)
				data, err = d.NotifyRet(c, notifyURL, ticketNO, orderNO, ip)
				convCtx.So(err, convey.ShouldNotBeNil)
				convCtx.So(data, convey.ShouldBeNil)
			})
		})
	})
}
