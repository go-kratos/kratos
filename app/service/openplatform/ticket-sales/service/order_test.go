package service

import (
	"context"
	"fmt"
	"testing"
	"time"

	rpc "go-common/app/service/openplatform/ticket-sales/api/grpc/v1"

	"go-common/app/service/openplatform/ticket-sales/api/grpc/type"

	"github.com/davecgh/go-spew/spew"
	"github.com/smartystreets/goconvey/convey"
)

func TestListOrders(t *testing.T) {
	convey.Convey("ListOrders", t, func() {
		req := &rpc.ListOrdersRequest{
			OrderID: []int64{100000003240240},
		}
		res, _ := svr.ListOrders(context.TODO(), req)
		convey.So(res.Count, convey.ShouldEqual, 1)
		o := res.List[0]
		convey.So(o.ItemInfo.Name, convey.ShouldEqual, "小凡の订单流程")
		convey.So(len(o.SKUs[0].SeatIDs), convey.ShouldBeGreaterThan, 1)
		convey.So(o.PayCharge.ChargeID, convey.ShouldEqual, "2946652108982083584")
		convey.So(o.Detail.Buyers[0].Tel, convey.ShouldEqual, "13524568885")
		req = &rpc.ListOrdersRequest{
			UID:    "27515328",
			ItemID: 662,
			Status: []int16{2},
		}
		res, _ = svr.ListOrders(context.TODO(), req)
		o = res.List[0]
		convey.So(o.Detail.Deliver.Tel, convey.ShouldEqual, "15129075612")
	})
}

func TestCreateOrders(t *testing.T) {
	convey.Convey("CreateOrders", t, func() {
		req := &rpc.CreateOrdersRequest{
			Orders: []*rpc.CreateOrderRequest{
				{
					ProjectID: 1191,
					ScreenID:  1738,
					SKUs: []*rpc.CreateOrderSKU{
						{
							SKUID: 13760,
							Count: 1,
						},
					},
					UID:       27515305,
					PayMoney:  30,
					OrderType: 1,
					TS:        time.Now().Unix(),
					Buyers: []*_type.OrderBuyer{
						{
							Name:       "Lisi",
							Tel:        "13800138000",
							PersonalID: "123456",
						},
					},
					DeliverDetail: &_type.OrderDeliver{
						Name: "Lisi",
						Addr: "Beijing",
						Tel:  "13800138000",
					},
				},
			},
		}
		res, err := svr.CreateOrders(context.TODO(), req)
		spew.Dump(res, err)
	})
}

func TestListOrderLogs(t *testing.T) {
	convey.Convey("ListOrders", t, func() {
		req := &rpc.ListOrderLogRequest{

			OrderID: 15198751542471,
			Limit:   0,
			Offset:  10,
			OrderBy: "",
		}

		res, _ := svr.ListOrderLogs(context.TODO(), req)
		fmt.Println(res)

	})
}
