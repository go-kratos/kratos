package alarm

import (
	"context"
	"testing"

	"go-common/app/service/main/resource/model"

	"github.com/smartystreets/goconvey/convey"
)

func TestAlarmstateDescribe(t *testing.T) {
	convey.Convey("stateDescribe", t, func(ctx convey.C) {
		ctx.Convey("When state = 0 is in the const map", func(ctx convey.C) {
			res := stateDescribe(0)
			ctx.Convey("Then res should equal 开放浏览", func(ctx convey.C) {
				ctx.So(res, convey.ShouldEqual, "开放浏览")
			})
		})
		ctx.Convey("When state = -99 is not in the const map", func(ctx convey.C) {
			res := stateDescribe(-99)
			ctx.Convey("Then res should equal -99", func(ctx convey.C) {
				ctx.So(res, convey.ShouldEqual, "-99")
			})
		})
	})
}

func TestAlarmSendWeChart(t *testing.T) {
	convey.Convey("SendWeChart", t, func() {
		convey.Convey("When everything is correct", func(ctx convey.C) {
			httpMock("POST", "http://bap.bilibili.co/api/v1/message/add").Reply(200).JSON("{}")
			err := d.SendWeChart(context.Background(), 1, "", []*model.ResWarnInfo{{}}, "unit test msg")
			ctx.Convey("Then err should be nil", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
		convey.Convey("When set http request gets 404", func(ctx convey.C) {
			httpMock("POST", "http://bap.bilibili.co/api/v1/message/add").Reply(404)
			err := d.SendWeChart(context.Background(), 1, "", []*model.ResWarnInfo{{}}, "unit test msg")
			ctx.Convey("Then err should not be nil", func(ctx convey.C) {
				ctx.So(err, convey.ShouldNotBeNil)
			})
		})
		convey.Convey("When set titleType == \"warn\"", func(ctx convey.C) {
			httpMock("POST", "http://bap.bilibili.co/api/v1/message/add").Reply(200).JSON("{}")
			err := d.SendWeChart(context.Background(), 1, "", []*model.ResWarnInfo{{}}, "warn")
			ctx.Convey("Then err should be nil", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestAlarmsendWeChartURL(t *testing.T) {
	convey.Convey("sendWeChartURL", t, func() {
		convey.Convey("When everything is correct", func(ctx convey.C) {
			httpMock("POST", "http://bap.bilibili.co/api/v1/message/add").Reply(200).JSON("{}")
			err := d.sendWeChartURL(context.Background(), 1, "", []*model.ResWarnInfo{{}})
			ctx.Convey("Then err should be nil", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
		convey.Convey("When set http request gets 404", func(ctx convey.C) {
			httpMock("POST", "http://bap.bilibili.co/api/v1/message/add").Reply(404)
			err := d.sendWeChartURL(context.Background(), 1, "", []*model.ResWarnInfo{{}})
			ctx.Convey("Then err should not be nil", func(ctx convey.C) {
				ctx.So(err, convey.ShouldNotBeNil)
			})
		})
	})
}
