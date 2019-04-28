package service

import (
	"testing"

	"go-common/app/service/openplatform/ticket-sales/api/grpc/v1"

	. "github.com/smartystreets/goconvey/convey"
)

func TestTicketView(t *testing.T) {
	Convey("TestTicketView", t, func() {
		req := &v1.TicketViewRequest{
			OrderID: 18000700140612,
		}
		tickets, err := svr.TicketView(ctx, req)
		So(tickets, ShouldNotBeNil)
		So(err, ShouldBeNil)
	})
}

func TestTicketSend(t *testing.T) {
	Convey("TestTicketSend", t, func() {
		req := &v1.TicketSendRequest{
			SendTID: []int64{400, 402},
		}
		sends, err := svr.TicketSend(ctx, req)
		So(sends, ShouldNotBeNil)
		So(err, ShouldBeNil)
	})
}
