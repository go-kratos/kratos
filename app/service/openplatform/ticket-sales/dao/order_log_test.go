package dao

import (
	"context"
	"fmt"
	"testing"

	"go-common/app/service/openplatform/ticket-sales/api/grpc/v1"

	. "github.com/smartystreets/goconvey/convey"
)

func TestListOrderLogs(t *testing.T) {

	Convey("ListOrderLogs", t, func() {

		data, err := d.GetOrderLogList(context.TODO(), 12222, 0, 10, "ctime")
		if err != nil {
			fmt.Println(err)
		}
		So(err, ShouldBeNil)
		So(data, ShouldNotBeNil)
	})
}

func TestAddOrderLogs(t *testing.T) {

	Convey("AddOrderLogs", t, func() {

		oi := &v1.OrderLog{}
		oi.UID = "wlt"
		oi.OpData = "test"
		oi.OID = 12222
		oi.OpName = "name"
		oi.OpObject = "object"
		oi.IP = "127.0.0.1"
		oi.Remark = "remark"

		data, err := d.AddOrderLog(context.TODO(), oi)
		if err != nil {
			fmt.Println(err)
		}

		So(err, ShouldBeNil)
		So(data, ShouldNotBeNil)
	})
}
