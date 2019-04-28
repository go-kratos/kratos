package service

import (
	"context"
	rpcV1 "go-common/app/service/openplatform/ticket-sales/api/grpc/v1"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestService_GetPromoOrder(t *testing.T) {
	Convey("GetPromoOrder", t, func() {
		res, err := svr.GetPromoOrder(context.TODO(), 1)
		So(err, ShouldBeNil)
		t.Logf("res:%v", res)
	})
}

func TestService_CheckCreateStatus(t *testing.T) {
	Convey("CheckCreateStatus", t, func() {
		res, err := svr.CheckCreateStatus(context.TODO(), &rpcV1.CheckCreatePromoOrderRequest{UID: 112, SKUID: 23, PromoID: 1})
		So(err, ShouldBeNil)
		t.Logf("res:%v", res)
	})
}

func TestService_GetUserJoinPromoOrder(t *testing.T) {
	Convey("GetUserJoinPromoOrder", t, func() {
		res, err := svr.GetUserJoinPromoOrder(context.TODO(), 1, 2, 0, 1112)
		So(err, ShouldBeNil)
		t.Logf("res:%v", res)
	})
}

func TestService_CreatePromoOrder(t *testing.T) {
	Convey("CreatePromoOrder", t, func() {
		res, err := svr.CreatePromoOrder(context.TODO(),
			&rpcV1.CreatePromoOrderRequest{
				PromoID:    1,
				OrderID:    1,
				GroupID:    0,
				UID:        1,
				PromoSKUID: 1,
				Ctime:      1,
				PayMoney:   1,
			})
		So(err, ShouldBeNil)
		t.Logf("res:%v", res)
	})
}

func TestService_CancelOrder(t *testing.T) {
	Convey("CancelOrder", t, func() {
		res, err := svr.CancelOrder(context.TODO(), &rpcV1.OrderID{OrderID: 1})
		So(err, ShouldNotBeNil)
		t.Logf("res:%v", res)
	})
}

func TestService_PromoPayNotify(t *testing.T) {
	Convey("PromoPayNotify", t, func() {
		res, err := svr.PromoPayNotify(context.TODO(), &rpcV1.OrderID{OrderID: 1})
		So(err, ShouldBeNil)
		t.Logf("res:%v", res)
	})
}

func TestService_CheckIssue(t *testing.T) {
	Convey("CheckIssue", t, func() {
		res, err := svr.CheckIssue(context.TODO(), &rpcV1.OrderID{OrderID: 1})
		So(err, ShouldBeNil)
		t.Logf("res:%v", res)
	})
}

func TestService_FinishIssue(t *testing.T) {
	Convey("FinishIssue", t, func() {
		res, err := svr.FinishIssue(context.TODO(), &rpcV1.FinishIssueRequest{PromoID: 1, GroupID: 1})
		So(err, ShouldBeNil)
		t.Logf("res:%v", res)
	})
}

func TestService_PromoRefundNotify(t *testing.T) {
	Convey("PromoRefundNotify", t, func() {
		res, err := svr.PromoRefundNotify(context.TODO(), &rpcV1.OrderID{OrderID: 1})
		So(err, ShouldBeNil)
		t.Logf("res:%v", res)
	})
}
