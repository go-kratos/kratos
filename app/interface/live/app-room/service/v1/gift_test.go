package v1

import (
	"context"
	"go-common/app/interface/live/app-room/api/http/v1"
	"go-common/app/interface/live/app-room/model"
	"go-common/library/net/metadata"
	"strconv"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"

	v1pb "go-common/app/interface/live/app-room/api/http/v1"

	"github.com/smartystreets/goconvey/convey"
)

func TestV1DailyBag(t *testing.T) {
	convey.Convey("DailyBag", t, func(c convey.C) {
		var (
			ctx = metadata.NewContext(context.Background(), metadata.MD{
				metadata.Mid: int64(88895029),
			})
			req = &v1pb.DailyBagReq{}
		)
		c.Convey("When everything gose positive", func(c convey.C) {
			resp, err := s.DailyBag(ctx, req)
			c.Convey("Then err should be nil.resp should not be nil.", func(c convey.C) {
				c.So(err, convey.ShouldBeNil)
				c.So(resp, convey.ShouldNotBeNil)
			})
		})
	})
}

func Test_RealGoldRecharge(t *testing.T) {
	Convey("normal", t, func() {
		So(1000, ShouldEqual, realRechargeGold(1000))
		So(2000, ShouldEqual, realRechargeGold(1300))
		So(2000, ShouldEqual, realRechargeGold(1700))
		So(1000, ShouldEqual, realRechargeGold(1))
		So(1000, ShouldEqual, realRechargeGold(233))
		So(2000, ShouldEqual, realRechargeGold(1398))

		ts, _ := time.Parse("2006-01-02", "2018-12-01")
		So(1, ShouldEqual, day(ts))
		So(201812, ShouldEqual, yearMonthNum(ts))

		ts, _ = time.Parse("2006-01-02", "2018-11-30")
		So(30, ShouldEqual, day(ts))
		So(201811, ShouldEqual, yearMonthNum(ts))

	})
}

func getContextWithMid(mid int64) context.Context {
	md := metadata.MD{
		metadata.Mid: mid,
	}
	return metadata.NewContext(context.Background(), md)
}

func TestGiftService_NeedTipRecharge_Silver(t *testing.T) {
	Convey("silver", t, func() {
		Convey("day empty", func() {
			ctx := getContextWithMid(1)
			testGiftService.conf.Gift.RechargeTip.SilverTipDays = []int{}
			resp, err := testGiftService.NeedTipRecharge(ctx, &v1.NeedTipRechargeReq{
				From:     v1.From_Silver,
				NeedGold: 0,
				Platform: "android",
			})
			So(err, ShouldBeNil)
			So(resp.Show, ShouldEqual, 0)
		})
		Convey("day hit and has coupon", func() {
			mid := int64(2)
			yearMonth, _ := strconv.ParseInt(time.Now().Format("200601"), 10, 64)
			testGiftService.dao.DelUserConf(context.Background(), mid, model.SilverTarget, yearMonth)
			testGiftService.dao.DelUserConf(context.Background(), mid, model.SilverTarget, StopPush)
			w, _ := testGiftService.dao.PayCenterWallet(context.Background(), mid, "android")
			if w.CouponBalance < 1 {
				t.Logf("own not enouth: %v", w.CouponBalance)
				return
			}
			testGiftService.conf.Gift.RechargeTip.SilverTipDays = append(testGiftService.conf.Gift.RechargeTip.SilverTipDays, day(time.Now()))
			t.Logf("%v", testGiftService.conf.Gift.RechargeTip.SilverTipDays)
			ctx := getContextWithMid(mid)
			resp, err := testGiftService.NeedTipRecharge(ctx, &v1.NeedTipRechargeReq{
				From:     v1.From_Silver,
				NeedGold: 0,
				Platform: "android",
			})
			So(err, ShouldBeNil)
			t.Logf("resp:%+v", resp)
			So(resp.Show, ShouldEqual, 1)
		})

		Convey("day hit and has coupon and set stop push", func() {
			mid := int64(2)
			yearMonth, _ := strconv.ParseInt(time.Now().Format("200601"), 10, 64)
			testGiftService.dao.DelUserConf(context.Background(), mid, model.SilverTarget, yearMonth)
			testGiftService.dao.DelUserConf(context.Background(), mid, model.SilverTarget, StopPush)
			w, _ := testGiftService.dao.PayCenterWallet(context.Background(), mid, "android")
			if w.CouponBalance < 1 {
				t.Logf("own not enouth: %v", w.CouponBalance)
				return
			}
			testGiftService.conf.Gift.RechargeTip.SilverTipDays = append(testGiftService.conf.Gift.RechargeTip.SilverTipDays, day(time.Now()))
			t.Logf("%v", testGiftService.conf.Gift.RechargeTip.SilverTipDays)
			ctx := getContextWithMid(mid)

			_, err := testGiftService.TipRechargeAction(ctx, &v1.TipRechargeActionReq{
				From:   v1.From_Silver,
				Action: v1.UserAction_StopPush,
			})
			So(err, ShouldBeNil)

			resp, err := testGiftService.NeedTipRecharge(ctx, &v1.NeedTipRechargeReq{
				From:     v1.From_Silver,
				NeedGold: 0,
				Platform: "android",
			})
			So(err, ShouldBeNil)
			t.Logf("resp:%+v", resp)
			So(resp.Show, ShouldEqual, 0)
		})

	})
}

func TestGiftService_NeedTipRecharge_Gold(t *testing.T) {
	Convey("normal", t, func() {
		mid := int64(2)
		ctx := getContextWithMid(mid)
		testGiftService.dao.DelUserConf(context.Background(), mid, model.GoldTarget, HasShow)
		w, _ := testGiftService.dao.PayCenterWallet(context.Background(), mid, "android")
		bpGold := int64(w.BcoinBalance * 1000)
		t.Logf("own : %d", bpGold)
		if bpGold < 3*1000 {
			t.Logf("own not enouth: %d", bpGold)
			return
		}
		lw, _ := testGiftService.dao.LiveWallet(context.Background(), mid, "android")
		if lw.GoldPayCnt > 0 {
			t.Logf("goldPayCnt enouth: %d", lw.GoldPayCnt)
			return
		}
		req := &v1.NeedTipRechargeReq{
			From:     v1.From_Gold,
			NeedGold: 666,
			Platform: "android",
		}
		resp, err := testGiftService.NeedTipRecharge(ctx, req)
		So(err, ShouldBeNil)
		t.Logf("resp:%+v", resp)
		So(resp.Show, ShouldEqual, 1)
		So(resp.Bp, ShouldEqual, w.BcoinBalance)
		So(resp.BpCoupon, ShouldEqual, w.CouponBalance)
		So(resp.RechargeGold, ShouldEqual, 1000)
	})
}

func TestGiftService_NeedTipRecharge(t *testing.T) {

	Convey("silver condition", t, func() {
		Convey("day empty", func() {
			testGiftService.dao.DelUserConf(context.Background(), 1, model.SilverTarget, StopPush)
			testGiftService.conf.Gift.RechargeTip.SilverTipDays = []int{}
			resp, err := testGiftService.silverNeedTipRecharge(context.Background(), 1, &v1.NeedTipRechargeReq{
				From:     v1.From_Silver,
				NeedGold: 0,
				Platform: "android",
			})
			So(err, ShouldBeNil)
			So(resp.Show, ShouldEqual, 0)
		})

		Convey("day hit and has coupon", func() {
			mid := int64(3)
			yearMonth, _ := strconv.ParseInt(time.Now().Format("200601"), 10, 64)
			testGiftService.dao.DelUserConf(context.Background(), mid, model.SilverTarget, yearMonth)
			testGiftService.dao.DelUserConf(context.Background(), mid, model.SilverTarget, StopPush)
			time.Sleep(time.Millisecond * 10) // 因为是异步设置所以需要一点时间延迟
			w, _ := testGiftService.dao.PayCenterWallet(context.Background(), mid, "android")
			if w.CouponBalance < 1 {
				t.Logf("own not enouth: %v", w.CouponBalance)
				return
			}
			testGiftService.conf.Gift.RechargeTip.SilverTipDays = append(testGiftService.conf.Gift.RechargeTip.SilverTipDays, day(time.Now()))
			t.Logf("%v", testGiftService.conf.Gift.RechargeTip.SilverTipDays)
			resp, err := testGiftService.silverNeedTipRecharge(context.Background(), mid, &v1.NeedTipRechargeReq{
				From:     v1.From_Silver,
				NeedGold: 0,
				Platform: "android",
			})
			So(err, ShouldBeNil)
			t.Logf("resp:%+v", resp)
			So(resp.Show, ShouldEqual, 1)
			time.Sleep(time.Millisecond * 20) // 因为是异步设置所以需要一点时间延迟
			resp, err = testGiftService.silverNeedTipRecharge(context.Background(), mid, &v1.NeedTipRechargeReq{
				From:     v1.From_Silver,
				NeedGold: 0,
				Platform: "android",
			})
			So(err, ShouldBeNil)
			t.Logf("resp:%+v", resp)
			So(resp.Show, ShouldEqual, 0)
		})

		Convey("day hit and has no coupon", func() {
			mid := int64(20)
			yearMonth, _ := strconv.ParseInt(time.Now().Format("200601"), 10, 64)
			testGiftService.dao.DelUserConf(context.Background(), mid, model.SilverTarget, yearMonth)
			testGiftService.dao.DelUserConf(context.Background(), mid, model.SilverTarget, StopPush)
			w, _ := testGiftService.dao.PayCenterWallet(context.Background(), mid, "android")
			if w.CouponBalance > 1 {
				t.Logf("own enouth: %v", w.CouponBalance)
				return
			}
			testGiftService.conf.Gift.RechargeTip.SilverTipDays = append(testGiftService.conf.Gift.RechargeTip.SilverTipDays, day(time.Now()))
			t.Logf("%v", testGiftService.conf.Gift.RechargeTip.SilverTipDays)
			resp, err := testGiftService.silverNeedTipRecharge(context.Background(), mid, &v1.NeedTipRechargeReq{
				From:     v1.From_Silver,
				NeedGold: 0,
				Platform: "android",
			})
			So(err, ShouldBeNil)
			t.Logf("resp:%+v", resp)
			So(resp.Show, ShouldEqual, 0)
		})

	})

	Convey("gold condition", t, func() {
		mid := int64(3)
		testGiftService.dao.DelUserConf(context.Background(), mid, model.GoldTarget, HasShow)
		w, _ := testGiftService.dao.PayCenterWallet(context.Background(), mid, "android")
		bpGold := int64(w.BcoinBalance * 1000)
		t.Logf("own : %d", bpGold)
		if bpGold < 3*1000 {
			t.Logf("own not enouth: %d", bpGold)
			return
		}
		lw, _ := testGiftService.dao.LiveWallet(context.Background(), mid, "android")
		if lw.GoldPayCnt > 0 {
			t.Logf("goldPayCnt enouth: %d", lw.GoldPayCnt)
			return
		}
		req := &v1.NeedTipRechargeReq{
			From:     v1.From_Gold,
			NeedGold: 1,
			Platform: "android",
		}
		ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*300)
		resp, err := testGiftService.goldNeedTipRecharge(ctx, mid, req)
		So(err, ShouldBeNil)
		t.Logf("resp:%+v", resp)
		So(resp.Show, ShouldEqual, 1)
		So(resp.Bp, ShouldEqual, w.BcoinBalance)
		So(resp.BpCoupon, ShouldEqual, w.CouponBalance)
		So(resp.RechargeGold, ShouldEqual, 1000)
		cancel()

		time.Sleep(time.Millisecond * 20) // 因为是异步设置所以需要一点时间延迟
		testGiftService.dao.DelUserConf(context.Background(), mid, model.GoldTarget, HasShow)

		req.NeedGold = 1230
		ctx, cancel = context.WithTimeout(context.Background(), time.Millisecond*300)
		resp, err = testGiftService.goldNeedTipRecharge(ctx, mid, req)
		cancel()
		So(err, ShouldBeNil)
		t.Logf("resp:%+v", resp)
		So(resp.Show, ShouldEqual, 1)
		So(resp.Bp, ShouldEqual, w.BcoinBalance)
		So(resp.BpCoupon, ShouldEqual, w.CouponBalance)
		So(resp.RechargeGold, ShouldEqual, 2000)
		time.Sleep(time.Millisecond * 10) // 因为是异步设置所以需要一点时间延迟
		testGiftService.dao.DelUserConf(context.Background(), mid, model.GoldTarget, HasShow)

		req.NeedGold = bpGold + 1
		ctx, cancel = context.WithTimeout(context.Background(), time.Millisecond*300)
		resp, err = testGiftService.goldNeedTipRecharge(ctx, mid, req)
		cancel()
		So(err, ShouldBeNil)
		t.Logf("resp:%+v", resp)
		So(resp.Show, ShouldEqual, 0)
	})
}
