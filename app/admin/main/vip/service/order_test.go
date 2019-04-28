package service

import (
	"context"
	"testing"

	"go-common/app/admin/main/vip/model"

	. "github.com/smartystreets/goconvey/convey"
)

func TestService_OrderList(t *testing.T) {
	Convey("test order list", t, func() {
		order := new(model.ArgPayOrder)
		order.Mid = 1001
		res, count, err := s.OrderList(context.TODO(), order)
		t.Logf("orderlist len:%+v count:%v", res, count)
		So(err, ShouldBeNil)
	})
}

func TestService_Refund(t *testing.T) {
	Convey("test refund ", t, func() {
		err := s.Refund(context.TODO(), "93035846180822184113", "zhaozhihao", 10)
		So(err, ShouldBeNil)
	})
}
