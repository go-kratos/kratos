package service

import (
	"fmt"
	"testing"

	coupinv1 "go-common/app/service/main/coupon/api"
	"go-common/app/service/main/coupon/model"

	"github.com/smartystreets/goconvey/convey"
	. "github.com/smartystreets/goconvey/convey"
)

// go test  -test.v -test.run TestAddAllowanceCoupon
func TestAddAllowanceCoupon(t *testing.T) {
	convey.Convey("TestAddAllowanceCoupon ", t, func() {
		var (
			err        error
			batchToken = "allowance_test1"
			bi         *model.CouponBatchInfo
			_mid       int64 = 233
			count            = 4
		)
		bi, err = s.dao.BatchInfo(c, batchToken)
		convey.So(err, convey.ShouldBeNil)
		convey.So(bi, convey.ShouldNotBeNil)
		err = s.AddAllowanceCoupon(c, bi, _mid, count, int64(1), 0)
		convey.So(err, convey.ShouldBeNil)
	})
}

// go test  -test.v -test.run TestAllowanceCoupon
func TestAllowanceCoupon(t *testing.T) {
	convey.Convey("TestAllowanceCoupon ", t, func() {
		var (
			err  error
			_mid int64 = 1
			res  []*model.CouponAllowanceInfo
		)
		res, err = s.AllowanceCoupon(c, &model.ArgAllowanceCoupons{
			Mid:   _mid,
			State: model.NotUsed,
		})
		t.Logf("count(%d)", len(res))
		for _, v := range res {
			t.Logf("v(%v)", v)
		}
		convey.So(err, convey.ShouldBeNil)
	})
}

// go test  -test.v -test.run TestUseAllowanceCoupon
func TestUseAllowanceCoupon(t *testing.T) {
	convey.Convey("TestUseAllowanceCoupon ", t, func() {
		var (
			err error
			arg = &model.ArgUseAllowance{
				Mid:         int64(2233),
				CouponToken: "275685366320181010160429",
				Remark:      "test1",
				Price:       float64(120),
				OrderNO:     "100101",
				Platform:    "pc",
			}
		)
		err = s.UseAllowanceCoupon(c, arg)
		convey.So(err, convey.ShouldBeNil)
	})
}

func TestUsableAllowanceCoupons(t *testing.T) {
	convey.Convey("TestUsableAllowanceCoupons ", t, func() {
		var (
			err    error
			_mid   = int64(1)
			_price = float64(120)
			us     []*model.CouponAllowancePanelInfo
			ds     []*model.CouponAllowancePanelInfo
			ui     []*model.CouponAllowancePanelInfo
		)
		us, ds, ui, err = s.UsableAllowanceCoupons(c, _mid, _price, []*model.CouponAllowanceInfo{}, 2, 1, 1)
		t.Logf("count(%d)", len(us))
		for _, v := range us {
			t.Logf("v(%v)", v)
		}
		t.Logf("count(%d)", len(ds))
		for _, v := range ds {
			t.Logf("v(%v)", v)
		}
		for _, v := range ui {
			t.Logf("v(%+v)", v)
		}
		convey.So(err, convey.ShouldBeNil)
	})
}

// go test  -test.v -test.run TestAllowancePanelCoupons
func TestAllowancePanelCoupons(t *testing.T) {
	convey.Convey("TestAllowancePanelCoupons ", t, func() {
		var (
			err    error
			_mid   = int64(233)
			_price = float64(100)
			us     []*model.CouponAllowancePanelInfo
			ds     []*model.CouponAllowancePanelInfo
			ui     []*model.CouponAllowancePanelInfo
		)
		us, ds, ui, err = s.AllowancePanelCoupons(c, _mid, _price, 3, 1, 1)
		t.Logf("count(%d)", len(us))
		for _, v := range us {
			t.Logf("v(%+v)", v)
		}
		t.Logf("count(%d)", len(ds))
		for _, v := range ds {
			t.Logf("v(%+v)", v)
		}
		for _, v := range ui {
			t.Logf("v(%+v)", v)
		}
		convey.So(err, convey.ShouldBeNil)
	})
}

// go test  -test.v -test.run TestMultiUsableAllowanceCoupon
func TestMultiUsableAllowanceCoupon(t *testing.T) {
	convey.Convey("TestMultiUsableAllowanceCoupon ", t, func() {
		var (
			err    error
			_mid   = int64(233)
			_price = append([]float64{}, float64(120))
			res    map[float64]*model.CouponAllowancePanelInfo
		)
		res, err = s.MultiUsableAllowanceCoupon(c, _mid, _price, 3, 1, 1)
		t.Logf("count(%d)", len(res))
		for _, v := range res {
			t.Logf("v(%+v)", v)
		}
		convey.So(err, convey.ShouldBeNil)
	})
}

// go test  -test.v -test.run TestCancelUseCoupon
func TestCancelUseCoupon(t *testing.T) {
	convey.Convey("TestCancelUseCoupon ", t, func() {
		var (
			err    error
			_mid   = int64(1)
			_token = "443385168420180705155534"
		)
		err = s.CancelUseCoupon(c, _mid, _token)
		convey.So(err, convey.ShouldBeNil)
	})
}

// go test  -test.v -test.run TestAllowanceList
func TestAllowanceList(t *testing.T) {
	convey.Convey("TestAllowanceList ", t, func() {
		var (
			err  error
			_mid = int64(1)
			res  []*model.CouponAllowancePanelInfo
		)
		res, err = s.AllowanceList(c, _mid, model.NotUsed)
		t.Logf("count(%d)", len(res))
		for _, v := range res {
			t.Logf("v(%+v)", v)
		}
		convey.So(err, convey.ShouldBeNil)
	})
}

// go test  -test.v -test.run TestScopeExplainFmt
func TestScopeExplainFmt(t *testing.T) {
	convey.Convey("TestScopeExplainFmt ", t, func() {
		c := &model.CouponAllowancePanelInfo{}
		c.ScopeExplainFmt("3", 1, 1, map[string]string{"1": "iPhone"})
		t.Logf("v(%+v)", c.ScopeExplain)
		convey.So(c.ScopeExplain != "", convey.ShouldBeTrue)
	})
}

// go test  -test.v -test.run TestUsableAllowanceCoupon
func TestUsableAllowanceCoupon(t *testing.T) {
	convey.Convey("TestUsableAllowanceCoupon ", t, func() {
		var (
			err    error
			_mid   = int64(233)
			_price = float64(120)
			res    *model.CouponAllowancePanelInfo
		)
		res, err = s.UsableAllowanceCoupon(c, _mid, _price, 2, 1, 1)
		convey.So(err, convey.ShouldBeNil)
		convey.So(res, convey.ShouldNotBeNil)
	})
}

// go test  -test.v -test.run TestCouponNotify
func TestCouponNotify(t *testing.T) {
	convey.Convey("TestCouponNotify ", t, func() {
		err := s.CouponNotify(c, _mid, "1807202041103", 1)
		convey.So(err, convey.ShouldBeNil)
	})
}

//go test -v -run TestReceiveAllowance
func TestReceiveAllowance(t *testing.T) {
	convey.Convey("test receive allowance", t, func() {
		arg := new(model.ArgReceiveAllowance)
		arg.Mid = 210
		arg.OrderNo = "1806141156245717948"
		arg.Appkey = "6a29f8ed87407c11"
		arg.BatchToken = "595017790420180808150514"
		couponToken, err := s.ReceiveAllowance(c, arg)
		t.Logf("%+v", couponToken)
		convey.So(err, convey.ShouldBeNil)
	})
}

// go test  -test.v -test.run TestUseNotify
func TestUseNotify(t *testing.T) {
	convey.Convey("TestUseNotify ", t, func() {
		_, err := s.UseNotify(c, &model.ArgAllowanceCheck{Mid: 1, OrderNo: "1"})
		convey.So(err, convey.ShouldBeNil)
	})
}

// go test  -test.v -test.run TestconvertCoupon
func TestConvertCoupon(t *testing.T) {
	convey.Convey("TestconvertCoupon ", t, func(convCtx convey.C) {
		r := s.convertCoupon(&model.CouponAllowanceInfo{}, []string{"test"}, float64(0), model.AllowanceDisables)
		convey.So(r, convey.ShouldNotBeNil)
	})
}

// go test  -test.v -test.run TestUsableAllowanceCouponV2
func TestUsableAllowanceCouponV2(t *testing.T) {
	Convey("TestUsableAllowanceCouponV2 ", t, func() {
		res, err := s.UsableAllowanceCouponV2(c, &coupinv1.UsableAllowanceCouponV2Req{
			Mid: 1,
			PriceInfo: []*coupinv1.ModelPriceInfo{
				{
					Price:          25,
					Plat:           1,
					ProdLimMonth:   1,
					ProdLimRenewal: 1,
				},
			},
		})
		fmt.Println("res", res, res.CouponInfo)
		So(err, ShouldBeNil)
		res, err = s.UsableAllowanceCouponV2(c, &coupinv1.UsableAllowanceCouponV2Req{
			Mid: 1,
			PriceInfo: []*coupinv1.ModelPriceInfo{
				{
					Price:          25,
					Plat:           1,
					ProdLimMonth:   1,
					ProdLimRenewal: 1,
				},
				{
					Price:          148,
					Plat:           1,
					ProdLimMonth:   12,
					ProdLimRenewal: 2,
				},
			},
		})
		fmt.Println("res", res, res.CouponInfo)
		So(err, ShouldBeNil)
		res, err = s.UsableAllowanceCouponV2(c, &coupinv1.UsableAllowanceCouponV2Req{
			Mid: 1,
			PriceInfo: []*coupinv1.ModelPriceInfo{
				{
					Price:          148,
					Plat:           1,
					ProdLimMonth:   12,
					ProdLimRenewal: 2,
				},
				{
					Price:          68,
					Plat:           1,
					ProdLimMonth:   3,
					ProdLimRenewal: 2,
				},
			},
		})
		fmt.Println("res", res, res.CouponInfo)
		So(err, ShouldBeNil)
	})
}
