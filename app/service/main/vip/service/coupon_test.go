package service

import (
	"fmt"
	"testing"

	col "go-common/app/service/main/coupon/model"
	v1 "go-common/app/service/main/vip/api"
	"go-common/app/service/main/vip/model"

	. "github.com/smartystreets/goconvey/convey"
)

//  go test  -test.v -test.run TestServiceCouponBySuitID
func TestServiceCouponBySuitID(t *testing.T) {
	Convey(" TestServiceCouponBySuitID ", t, func() {
		var (
			mid = int64(1)
			sid = int64(3)
			res *col.CouponAllowancePanelInfo
		)
		res, err := s.CouponBySuitID(c, &model.ArgCouponPanel{Mid: mid, Sid: sid})
		t.Logf("data(%+v)", res)
		So(err, ShouldBeNil)
	})
}

func TestServiceCouponsForPanel(t *testing.T) {
	Convey(" TestServiceCouponsForPanel ", t, func() {
		var (
			mid = int64(1)
			sid = int64(3)
			res *col.CouponAllowancePanelResp
		)
		res, err := s.CouponsForPanel(c, &model.ArgCouponPanel{Mid: mid, Sid: sid})
		t.Logf("data(+%v)", res)
		So(err, ShouldBeNil)
	})
}

// go test  -test.v -test.run TestServiceCancelUseCoupon

func TestServiceCancelUseCoupon(t *testing.T) {
	Convey(" TestServiceCancelUseCoupon ", t, func() {
		var (
			mid   = int64(1)
			token = "772991379820180716121343"
		)
		err := s.CancelUseCoupon(c, mid, token, "127.0.0.1")
		So(err, ShouldBeNil)
	})
}

// go test  -test.v -test.run TestServiceVipUserPanelV4
func TestServiceVipUserPanelV4(t *testing.T) {
	Convey(" TestServiceVipUserPanelV4 ", t, func() {
		var (
			mid = int64(1)
			res *model.VipPirceResp
		)
		res, err := s.VipUserPanelV4(c, &model.ArgPanel{
			Mid:      mid,
			Platform: "android",
			SortTp:   model.PanelMonthDESC,
		})
		So(res, ShouldNotBeNil)
		for _, v := range res.Vps {
			t.Logf("panel info(%+v)", v)
		}
		t.Logf("coupon info(%+v)", res.CouponInfo)
		t.Logf("data(+%v)", res)
		So(err, ShouldBeNil)
	})
}

// go test  -test.v -test.run TestServiceVipPrice
func TestServiceVipPrice(t *testing.T) {
	Convey(" TestServiceVipPrice ", t, func() {
		var (
			mid     = int64(1)
			plat    = int64(3)
			subType int8
			month   = int16(12)
			token   = "992628254320180713122015"
			res     *model.VipPirce
		)
		res, err := s.VipPrice(c, mid, month, plat, subType, token, "pc", 0)
		So(res, ShouldNotBeNil)
		t.Logf("coupon info(%+v)", res.Coupon)
		t.Logf("price info(%+v)", res.Panel)
		t.Logf("data(+%v)", res)
		So(err, ShouldBeNil)
	})
}

// go test  -test.v -test.run TestServiceCancelUseCoupon2

func TestServiceCancelUseCoupon2(t *testing.T) {
	Convey(" TestServiceCancelUseCoupon2 ", t, func() {
		var (
			mid   = int64(1)
			token = "992628254320180713122015"
		)
		err := s.CancelUseCoupon(c, mid, token, "127.0.0.1")
		So(err, ShouldBeNil)
	})
}

func TestServiceFormatRate(t *testing.T) {
	Convey("test formatRate", t, func() {
		config := new(model.VipPriceConfig)
		config.OPrice = 233
		config.DPrice = 148
		t.Logf(config.FormatRate())
	})
}

func TestServiceCouponBySuitIDV2(t *testing.T) {
	Convey("TestServiceCouponBySuitIDV2 ", t, func() {
		res, err := s.CouponBySuitIDV2(c, &v1.CouponBySuitIDReq{
			Sid:       158,
			Mid:       1,
			Platform:  "pc",
			PanelType: "normal",
		})
		fmt.Println("res", res.CouponTip, res.CouponInfo)
		So(err, ShouldBeNil)
	})
}
