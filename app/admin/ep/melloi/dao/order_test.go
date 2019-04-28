package dao

import (
	"testing"
	"time"

	"go-common/app/admin/ep/melloi/model"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	testOrder = model.Order{
		ApplyDate: time.Now(),
		ID:        1,
	}
	pageNum  int32 = 1
	pageSize int32 = 2
)

func Test_Order(t *testing.T) {
	Convey("test QueryOrder", t, func() {

		var workOrders *model.QueryOrderResponse
		order := model.Order{ID: testOrder.ID}
		workOrders, _ = d.QueryOrder(&order, pageNum, pageSize)
		So(workOrders, ShouldBeEmpty)
	})

	Convey("test DeleteOrder", t, func() {
		order := model.Order{ID: testOrder.ID}
		err := d.DelOrder(order.ID)
		So(err, ShouldBeNil)
	})

	Convey("test UpdateOrder", t, func() {
		order := model.Order{ID: testOrder.ID, ApplyDate: testOrder.ApplyDate}
		err := d.UpdateOrder(&order)
		So(err, ShouldBeNil)
	})

}
