package block

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
	"gopkg.in/h2non/gock.v1"
)

func TestBlockSendSysMsg(t *testing.T) {
	convey.Convey("SendSysMsg", t, func(convCtx convey.C) {
		var (
			c        = context.Background()
			code     = ""
			mids     = []int64{}
			title    = ""
			content  = ""
			remoteIP = ""
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			defer gock.OffAll()
			httpMock("POST", d.conf.BlockProperty.MSGURL).Reply(200).JSON(`{"code":0,"data":{"status":1,"remark":"test"}}`)
			err := d.SendSysMsg(c, code, mids, title, content, remoteIP)
			//println("err==",err)
			convCtx.Convey("Then err should be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
			})
		})
		convCtx.Convey("When everything goes negative", func(convCtx convey.C) {
			defer gock.OffAll()
			httpMock("POST", d.conf.BlockProperty.MSGURL).Reply(200).JSON(`{"code":500,"data":{"status":1,"remark":"test"}}`)
			//d.httpClient.SetTransport(gock.DefaultTransport)
			err := d.SendSysMsg(c, code, mids, title, content, remoteIP)
			convCtx.Convey("Then err should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldNotBeNil)
			})
		})

	})
}

func TestBlockmidsToParam(t *testing.T) {
	convey.Convey("midsToParam", t, func(convCtx convey.C) {
		var (
			mids = []int64{}
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			str := midsToParam(mids)
			convCtx.Convey("Then str should not be nil.", func(convCtx convey.C) {
				convCtx.So(str, convey.ShouldNotBeNil)
			})
		})
	})
}
