package dao

import (
	"fmt"
	"testing"
	"time"

	"go-common/app/service/openplatform/ticket-sales/api/grpc/type"
	"go-common/app/service/openplatform/ticket-sales/model"

	"net"

	. "github.com/smartystreets/goconvey/convey"
)

func TestUpdate(t *testing.T) {
	Convey("Get Order", t, func() {
		data := &model.OrderDetail{
			OrderID:    10000012001970,
			Buyer:      "TEST",
			Tel:        "13800138000",
			PersonalID: "342921",
			ExpressCO:  "shunfeng",
			ExpressNO:  "000",
			Remark:     "TEST",
			DeviceType: 1,
			IP:         net.ParseIP("::1"),
			DeliverDetail: &_type.OrderDeliver{
				AddrID: 1,
				Name:   "张三",
				Tel:    "13810559189",
				Addr:   "北京市",
			},
			Detail: &_type.OrderExtra{
				AutoRecvTime:   time.Now().Unix(),
				DelayRecvTimes: 1,
			},
		}

		// effID, err := d.UpdateDetail(context.TODO(), data)
		fmt.Println(data)

		// So(err, ShouldBeNil)
	})
}
