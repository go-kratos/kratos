package service

import (
	"testing"

	"go-common/app/job/main/vip/model"

	. "github.com/smartystreets/goconvey/convey"
)

//go test  -test.v -test.run  TestAutoRenewPay
func TestAutoRenewPay(t *testing.T) {
	Convey("TestAutoRenewPay ", t, func() {
		err := s.autoRenewPay(c, &model.VipUserInfoMsg{
			IsAutoRenew:  model.AutoRenew,
			PayChannelID: 100,
			Status:       0,
			Type:         1,
		}, &model.VipUserInfoMsg{
			Type:         1,
			Status:       1,
			IsAutoRenew:  model.AutoRenew,
			PayChannelID: 1,
		})
		So(err, ShouldBeNil)
		err = s.autoRenewPay(c, &model.VipUserInfoMsg{
			IsAutoRenew:  model.AutoRenew,
			PayChannelID: 1,
			Status:       0,
			Type:         1,
		}, &model.VipUserInfoMsg{
			Type:         1,
			Status:       1,
			IsAutoRenew:  model.AutoRenew,
			PayChannelID: 1,
		})
		So(err, ShouldBeNil)
	})
}
