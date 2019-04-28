package service

import (
	"testing"

	"go-common/app/job/main/vip/model"

	. "github.com/smartystreets/goconvey/convey"
)

//go test  -test.v -test.run  TestCouponNotify
func TestCouponNotify(t *testing.T) {
	Convey("TestCouponNotify ", t, func() {
		o := &model.VipPayOrderNewMsg{
			OrderNo: "1807211806450011799",
			Status:  model.SUCCESS,
			Mid:     1,
		}
		err := s.CouponNotify(c, o)
		So(err, ShouldBeNil)
	})
}
