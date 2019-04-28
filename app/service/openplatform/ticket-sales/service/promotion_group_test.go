package service

import (
	"context"
	rpcV1 "go-common/app/service/openplatform/ticket-sales/api/grpc/v1"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestService_GetPromoGroupInfo(t *testing.T) {
	Convey("GetPromoGroupInfo", t, func() {
		res, err := svr.GetPromoGroupInfo(context.TODO(), &rpcV1.GetPromoGroupInfoRequest{OrderID: 123})
		So(err, ShouldBeNil)
		t.Logf("res:%v", res)
	})
}

func TestService_GroupFailed(t *testing.T) {
	Convey("GroupFailed", t, func() {
		res, err := svr.GroupFailed(context.TODO(), &rpcV1.GroupFailedRequest{GroupID: 1, CancelNum: 1})
		So(err, ShouldBeNil)
		t.Logf("res:%v", res)
	})
}
