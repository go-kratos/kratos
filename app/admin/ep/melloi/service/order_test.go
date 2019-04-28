package service

import (
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
	"go-common/app/admin/ep/melloi/model"
)

var (
	testOrder = model.QueryOrderRequest{
		Order:      model.Order{ID: 1, Name: "test", ApplyDate: time.Now()},
		Pagination: model.Pagination{PageNum: 1, PageSize: 1, TotalSize: 1},
	}
)

func Test_Order(t *testing.T) {
	Convey("test QueryOrder", t, func() {

		var workOrders *model.QueryOrderResponse
		workOrders, _ = s.QueryOrder(&testOrder)
		So(workOrders, ShouldBeEmpty)
	})

	Convey("test DeleteOrder", t, func() {
		order := model.Order{ID: testOrder.ID}
		err := s.DelOrder(order.ID)
		So(err, ShouldBeNil)
	})

	Convey("test UpdateOrder", t, func() {
		order := model.Order{ID: testOrder.ID, ApplyDate: testOrder.ApplyDate}
		err := s.UpdateOrder(&order)
		So(err, ShouldBeNil)
	})

}
