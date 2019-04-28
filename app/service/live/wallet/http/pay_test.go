package http

import (
	"context"
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"go-common/app/service/live/wallet/model"
	"net/url"
	"testing"
)

func queryPay(t *testing.T, form *model.RechargeOrPayForm, platform string) *RechargeRes {
	params := url.Values{}
	params.Set("uid", fmt.Sprintf("%d", form.Uid))
	params.Set("coin_type", form.CoinType)
	params.Set("coin_num", fmt.Sprintf("%d", form.CoinNum))
	params.Set("extend_tid", form.ExtendTid)
	params.Set("timestamp", fmt.Sprintf("%d", form.Timestamp))
	params.Set("transaction_id", form.TransactionId)

	req, _ := client.NewRequest("POST", _payURL, "127.0.0.1", params)
	req.Header.Set("platform", platform)

	var res RechargeRes

	err := client.Do(context.TODO(), req, &res)
	if err != nil {
		t.Errorf("client.Do() error(%v)", err)
		t.FailNow()
	}
	return &res
}

func queryPayWithReason(t *testing.T, form *model.RechargeOrPayForm, platform string, reason string) *RechargeRes {
	params := url.Values{}
	params.Set("uid", fmt.Sprintf("%d", form.Uid))
	params.Set("coin_type", form.CoinType)
	params.Set("coin_num", fmt.Sprintf("%d", form.CoinNum))
	params.Set("extend_tid", form.ExtendTid)
	params.Set("timestamp", fmt.Sprintf("%d", form.Timestamp))
	params.Set("transaction_id", form.TransactionId)
	params.Set("reason", reason)

	req, _ := client.NewRequest("POST", _payURL, "127.0.0.1", params)
	req.Header.Set("platform", platform)

	var res RechargeRes

	err := client.Do(context.TODO(), req, &res)
	if err != nil {
		t.Errorf("client.Do() error(%v)", err)
		t.FailNow()
	}
	return &res
}

func TestPay(t *testing.T) {
	once.Do(startHTTP)

	Convey("pay normal 先调用get接口　再调用pay 再调用get接口　比较用户钱包数据", t, func() {
		platforms := []string{"pc", "android", "ios"}
		var num int64 = 1000
		var payNum int64 = 100
		uid := getTestRandUid()
		for _, platform := range platforms {

			beforeWallet := getTestWallet(t, uid, platform)

			res := queryRecharge(t, getTestRechargeOrPayForm(t, int32(model.RECHARGETYPE), uid, "gold", num, nil), platform)

			So(res.Code, ShouldEqual, 0)
			So(getIntCoinForTest(res.Resp.Gold)-getIntCoinForTest(beforeWallet.Gold), ShouldEqual, num)

			res = queryRecharge(t, getTestRechargeOrPayForm(t, int32(model.RECHARGETYPE), uid, "silver", num, nil), platform)

			So(res.Code, ShouldEqual, 0)
			So(getIntCoinForTest(res.Resp.Silver)-getIntCoinForTest(beforeWallet.Silver), ShouldEqual, num)

			afterWallet := getTestWallet(t, uid, platform)

			So(getIntCoinForTest(afterWallet.Gold)-getIntCoinForTest(beforeWallet.Gold), ShouldEqual, num)
			So(getIntCoinForTest(afterWallet.Silver)-getIntCoinForTest(beforeWallet.Silver), ShouldEqual, num)

			f1 := getTestRechargeOrPayForm(t, int32(model.PAYTYPE), uid, "gold", payNum, nil)
			res = queryPay(t, f1, platform)
			So(res.Code, ShouldEqual, 0)
			So(getIntCoinForTest(res.Resp.Gold)-getIntCoinForTest(afterWallet.Gold), ShouldEqual, -1*payNum)

			sr := queryStatus(t, uid, f1.TransactionId)
			So(sr.Code, ShouldEqual, 0)
			So(sr.Resp.Status, ShouldEqual, 0)

			res = queryPay(t, getTestRechargeOrPayForm(t, int32(model.PAYTYPE), uid, "silver", payNum, nil), platform)
			So(res.Code, ShouldEqual, 0)
			So(getIntCoinForTest(res.Resp.Silver)-getIntCoinForTest(afterWallet.Silver), ShouldEqual, -1*payNum)

			payWallet := getTestWallet(t, uid, platform)

			So(getIntCoinForTest(payWallet.Gold)-getIntCoinForTest(afterWallet.Gold), ShouldEqual, -1*payNum)
			So(getIntCoinForTest(payWallet.Silver)-getIntCoinForTest(afterWallet.Silver), ShouldEqual, -1*payNum)

		}
	})
}

func TestPayMetal(t *testing.T) {
	once.Do(startHTTP)

	Convey("pay metal", t, func() {
		var uid int64 = 1
		platform := "pc"
		f1 := getTestRechargeOrPayForm(t, int32(model.PAYTYPE), uid, "metal", 1, nil)
		res := queryPayWithReason(t, f1, platform, "ut")
		So(res.Code == 0 || res.Code == 1000000, ShouldBeTrue)
	})
}
