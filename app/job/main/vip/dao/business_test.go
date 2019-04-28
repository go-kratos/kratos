package dao

import (
	"context"
	"testing"
	"time"

	"go-common/app/job/main/vip/model"
	xtime "go-common/library/time"

	"github.com/smartystreets/goconvey/convey"
	gock "gopkg.in/h2non/gock.v1"
)

//  go test  -test.v -test.run TestDaoSalaryCoupon
func TestDaoSalaryCoupon(t *testing.T) {
	convey.Convey("TestDaoSalaryCoupon salary coupon", t, func() {
		var (
			c           = context.TODO()
			mid   int64 = 123
			ct    int8  = 2
			count int64 = 2
			err   error
		)
		err = d.SalaryCoupon(c, mid, ct, count, "cartoon_1_2018_06")
		convey.So(err, convey.ShouldBeNil)
	})
}

func TestDao_SendMultipMsg(t *testing.T) {
	convey.Convey("send multipmsg", t, func() {
		defer gock.OffAll()
		httpMock("POST", _message).Reply(200).JSON(`{"code":0,"data":1}`)
		err := d.SendMultipMsg(context.TODO(), "27515256", "test", "test", "10_1_2", 4)
		convey.So(err, convey.ShouldBeNil)
	})
}

func TestDaoPushData(t *testing.T) {
	pushData := &model.VipPushData{
		Title:         "TEST",
		PushStartTime: "15:04:05",
		PushEndTime:   "15:04:05",
	}
	convey.Convey("PushData", t, func() {
		defer gock.OffAll()
		httpMock("POST", _pushData).Reply(200).JSON(`{"code":0,"data":1}`)
		rel, err := d.PushData(context.TODO(), []int64{7593623}, pushData, "2006-01-02")
		convey.So(err, convey.ShouldBeNil)
		convey.So(rel, convey.ShouldNotBeNil)
	})
}

func TestDaoSendMedal(t *testing.T) {
	convey.Convey("SendMedal", t, func() {
		defer gock.OffAll()
		httpMock("GET", _sendMedal).Reply(200).JSON(`{"code":0,"data":{"status":1}}`)
		status := d.SendMedal(context.TODO(), 0, 0)
		convey.So(status, convey.ShouldNotBeNil)
	})
}

func TestDaoSendCleanCache(t *testing.T) {
	hv := &model.HandlerVip{Mid: 7593623}
	convey.Convey("SendCleanCache", t, func() {
		defer gock.OffAll()
		httpMock("GET", _cleanCache).Reply(200).JSON(`{"code":0,"data":{"status":1}}`)
		err := d.SendCleanCache(context.TODO(), hv)
		convey.So(err, convey.ShouldBeNil)
	})
}

func TestDaoSendBcoin(t *testing.T) {
	convey.Convey("SendBcoin", t, func() {
		defer gock.OffAll()
		httpMock("POST", _addBcoin).Reply(200).JSON(`{"code":0,"data":{"status":1}}`)
		err := d.SendBcoin(context.TODO(), []int64{7593623}, 0, xtime.Time(time.Now().Unix()), "")
		convey.So(err, convey.ShouldBeNil)
	})
}

func TestDaoSendAppCleanCache(t *testing.T) {
	var (
		hv  = &model.HandlerVip{Mid: 7593623}
		app = &model.VipAppInfo{
			PurgeURL: "http://bilibili.com/test",
		}
	)
	convey.Convey("SendAppCleanCache", t, func() {
		defer gock.OffAll()
		httpMock("GET", app.PurgeURL).Reply(200).JSON(`{"code":0}`)
		err := d.SendAppCleanCache(context.TODO(), hv, app)
		convey.So(err, convey.ShouldBeNil)
	})
}

func TestDaosortParamsKey(t *testing.T) {
	var v map[string]string
	convey.Convey("sortParamsKey", t, func() {
		p1 := d.sortParamsKey(v)
		convey.So(p1, convey.ShouldNotBeNil)
	})
}

func TestDaoPaySign(t *testing.T) {
	var params map[string]string
	convey.Convey("PaySign", t, func() {
		sign := d.PaySign(params, "test")
		convey.So(sign, convey.ShouldNotBeNil)
	})
}

func TestDaodoNomalSend(t *testing.T) {
	var path = "/x/internal/vip/user/info"
	convey.Convey("doNomalSend", t, func() {
		defer gock.OffAll()
		httpMock("POST", path).Reply(200).JSON(`{"code":0}`)
		err := d.doNomalSend(context.TODO(), "http://api.bilibili.com", path, "", nil, nil, new(model.VipPushResq))
		convey.So(err, convey.ShouldBeNil)
	})
}
