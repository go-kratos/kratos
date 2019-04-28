package service

import (
	"context"
	rpcV1 "go-common/app/service/openplatform/ticket-sales/api/grpc/v1"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestService_GetPromo(t *testing.T) {
	Convey("GetPromo", t, func() {
		res, err := svr.GetPromo(context.TODO(), &rpcV1.PromoID{PromoID: 1})
		So(err, ShouldBeNil)
		t.Logf("res:%v", res)
	})
}

func TestService_CreatePromo(t *testing.T) {
	Convey("CreatePromo", t, func() {
		res, err := svr.CreatePromo(context.TODO(), &rpcV1.CreatePromoRequest{PromoID: 1})
		So(err, ShouldNotBeNil)
		t.Logf("res:%v", res)
	})
}

func TestService_HasPromoOfSKU(t *testing.T) {
	Convey("HasPromoOfSKU", t, func() {
		res, err := svr.HasPromoOfSKU(context.TODO(), 1, 1, 1)
		So(err, ShouldBeNil)
		t.Logf("res:%v", res)
	})
}

func TestService_OperatePromo(t *testing.T) {
	Convey("OperatePromo", t, func() {
		res, err := svr.OperatePromo(context.TODO(), &rpcV1.OperatePromoRequest{PromoID: 1, OperateType: 1})
		So(err, ShouldNotBeNil)
		t.Logf("res:%v", res)
	})
}

func TestService_CheckPromoStatus(t *testing.T) {
	Convey("CheckPromoStatus", t, func() {
		res, err := svr.CheckPromoStatus(context.TODO(), 1, 1)
		So(err, ShouldNotBeNil)
		t.Logf("res:%v", res)
	})
}

func TestService_EditPromo(t *testing.T) {
	Convey("EditPromo", t, func() {
		res, err := svr.EditPromo(context.TODO(), &rpcV1.EditPromoRequest{PromoID: 1, Amount: 100})
		So(err, ShouldNotBeNil)
		t.Logf("res:%v", res)
	})
}
