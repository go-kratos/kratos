package service

import (
	"testing"

	"go-common/app/service/main/vip/model"

	. "github.com/smartystreets/goconvey/convey"
)

// go test  -test.v -test.run TestServiceCreateOrder
func TestServiceCreateOrder(t *testing.T) {
	Convey("shold return  true where err == nil", t, func() {
		arg := &model.ArgCreateOrder{
			Mid:       1,
			AppID:     1,
			AppSubID:  "20",
			Months:    3,
			OrderType: 0,
			DType:     0,
			Platform:  "android",
		}
		pp, err := s.CreateOrder(c, arg, "")
		t.Logf("data(+%v)", pp)
		So(err, ShouldBeNil)
	})
}

func TestServicePrice(t *testing.T) {
	Convey(" price test ", t, func() {
		price, err := s.Price(c, 1, 1, 0, 0)
		t.Logf("data(+%v)", price)
		So(err, ShouldBeNil)
	})
}

// go test  -test.v -test.run TestServiceOrderList
func TestServiceOrderList(t *testing.T) {
	Convey(" price test ", t, func() {
		data, count, err := s.OrderList(c, int64(1684013), 1, 10)
		t.Logf("data(+%v)  count(%d)", data, count)
		So(err, ShouldBeNil)
	})
}

func TestServiceOrderInfo(t *testing.T) {
	Convey(" OrderInfo test ", t, func() {
		o, err := s.OrderInfo(c, "1")
		t.Logf("data(+%v)", o)
		So(err, ShouldBeNil)
	})
}

// go test  -test.v -test.run TestServiceOrderMng
func TestServiceOrderMng(t *testing.T) {
	Convey(" OrderMng test ", t, func() {
		mng, err := s.OrderMng(c, int64(iapAutoRenewOverdueMid))
		t.Logf("data(+%v)", mng)
		So(err, ShouldBeNil)
		So(mng.IsAutoRenew == 0, ShouldBeTrue)

		mng, err = s.OrderMng(c, int64(wechatAutoRenewOverdueMid))
		t.Logf("data(+%v)", mng)
		So(err, ShouldBeNil)
		So(mng.IsAutoRenew == 0, ShouldBeTrue)

		mng, err = s.OrderMng(c, int64(wechatAutoRenewNotOverdueMid))
		t.Logf("data(+%v)", mng)
		So(err, ShouldBeNil)
		So(mng.IsAutoRenew == 1, ShouldBeTrue)

		mng, err = s.OrderMng(c, int64(iapAutoRenewTodayOverdueMid))
		t.Logf("data(+%v)", mng)
		So(err, ShouldBeNil)
		So(mng.IsAutoRenew == 1, ShouldBeTrue)
		So(mng.PayType == "苹果应用内购买", ShouldBeTrue)

		mng, err = s.OrderMng(c, int64(iapAutoRenewNotOverdueMid))
		t.Logf("data(+%v)", mng)
		So(err, ShouldBeNil)
		So(mng.IsAutoRenew == 1, ShouldBeTrue)
		So(mng.PayType == "苹果应用内购买", ShouldBeTrue)

	})
}

// go test  -test.v -test.run TestServiceCreateOldOrder
func TestServiceCreateOldOrder(t *testing.T) {
	Convey("old order shold return true where err == nil", t, func() {
		arg := &model.ArgOldPayOrder{
			OrderNo:      "22kkk2222",
			AppID:        1,
			Platform:     1,
			OrderType:    0,
			AppSubID:     "1",
			Mid:          1,
			BuyMonths:    3,
			Money:        1,
			Status:       1,
			PayType:      0,
			RechargeBp:   1,
			ThirdTradeNo: "11",
		}
		err := s.CreateOldOrder(c, arg)
		So(err, ShouldBeNil)
	})
}

// go test  -test.v -test.run TestServiceMemberInfo
func TestServiceMemberInfo(t *testing.T) {
	Convey(" MemberInfo test ", t, func() {
		v, err := s.memInfoRetry(c, int64(1))
		t.Logf("data(+%v)", v)
		So(err, ShouldBeNil)
	})
}
