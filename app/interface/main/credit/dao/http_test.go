package dao

import (
	"context"
	"testing"

	gock "gopkg.in/h2non/gock.v1"

	"github.com/smartystreets/goconvey/convey"
)

// TestReplysCount .
func TestSendSysMsg(t *testing.T) {
	convey.Convey("return someting", t, func(convCtx convey.C) {
		defer gock.OffAll()
		httpMock("POST", d.sendMsgURL).Reply(200).JSON(`{"code": 0}`)
		err := d.SendSysMsg(context.Background(), 1, "test", "test")
		convCtx.So(err, convey.ShouldBeNil)
	})
}
func TestGetQS(t *testing.T) {
	convey.Convey("return someting", t, func(convCtx convey.C) {
		defer gock.OffAll()
		httpMock("GET", d.getQSURL).Reply(200).JSON(`{"code": 0,"data":{}}`)
		res, err := d.GetQS(context.Background(), 1)
		convCtx.So(err, convey.ShouldBeNil)
		convCtx.So(res, convey.ShouldNotBeNil)
	})
}
func TestReplysCount(t *testing.T) {
	convey.Convey("return someting", t, func(convCtx convey.C) {
		res, err := d.ReplysCount(context.Background(), []int64{2, 3, 4})
		convCtx.So(err, convey.ShouldBeNil)
		convCtx.So(res, convey.ShouldNotBeNil)
	})
}

func TestSendMedal(t *testing.T) {
	convey.Convey("return someting", t, func(convCtx convey.C) {
		defer gock.OffAll()
		httpMock("POST", d.sendMedalURL).Reply(200).JSON(`{"code": 0}`)
		err := d.SendMedal(context.Background(), 1, 1)
		convCtx.So(err, convey.ShouldBeNil)
	})
}
