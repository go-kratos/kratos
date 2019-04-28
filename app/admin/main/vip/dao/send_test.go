package dao

import (
	"context"
	"go-common/app/admin/main/vip/model"
	"testing"

	"github.com/smartystreets/goconvey/convey"
	gock "gopkg.in/h2non/gock.v1"
)

func TestDaoSendMultipMsg(t *testing.T) {
	convey.Convey("SendMultipMsg", t, func() {
		defer gock.OffAll()
		httpMock("POST", _SendUserNotify).Reply(200).JSON(`{"code":0}`)
		err := d.SendMultipMsg(context.TODO(), "", "", "", "", "", 0)
		convey.So(err, convey.ShouldBeNil)
	})
}

func TestDaoPayRefund(t *testing.T) {
	var (
		arg                  = &model.PayOrder{}
		refundAmount float64 = 1.0
		refundID             = "test001"
	)
	convey.Convey("PayRefund", t, func() {
		defer gock.OffAll()
		httpMock("POST", _payRefund).Reply(200).JSON(`{"code":0}`)
		err := d.PayRefund(context.Background(), arg, refundAmount, refundID)
		convey.So(err, convey.ShouldBeNil)
	})
}
