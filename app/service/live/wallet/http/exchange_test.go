package http

import (
	"context"
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"go-common/app/service/live/wallet/model"
	"net/url"
	"testing"
	"time"
)

func queryExchange(t *testing.T, form *model.ExchangeForm, platform string) *RechargeRes {
	params := url.Values{}
	params.Set("uid", fmt.Sprintf("%d", form.Uid))
	params.Set("src_coin_type", form.SrcCoinType)
	params.Set("src_coin_num", fmt.Sprintf("%d", form.SrcCoinNum))
	params.Set("dest_coin_type", form.DestCoinType)
	params.Set("dest_coin_num", fmt.Sprintf("%d", form.DestCoinNum))
	params.Set("extend_tid", form.ExtendTid)
	params.Set("timestamp", fmt.Sprintf("%d", form.Timestamp))
	params.Set("transaction_id", form.TransactionId)

	req, _ := client.NewRequest("POST", _exchangeURL, "127.0.0.1", params)
	req.Header.Set("platform", platform)

	var res RechargeRes

	err := client.Do(context.TODO(), req, &res)
	if err != nil {
		t.Errorf("client.Do() error(%v)", err)
		t.FailNow()
	}
	return &res
}

func getTestExchangeForm(t *testing.T, uid int64, srcCoinType string, srcCoinNum int64, destCoinType string, destCoinNum int64, tid interface{}) *model.ExchangeForm {

	if tid == nil {
		res := queryGetTid(t, int32(model.EXCHANGETYPE), getTestParamsJson())
		if res.Code != 0 {
			t.Errorf("get tid failed code : %d", res.Code)
			t.FailNow()
		}
		tid = res.Resp.TransactionId
	}
	return &model.ExchangeForm{
		Uid:           uid,
		SrcCoinType:   srcCoinType,
		SrcCoinNum:    srcCoinNum,
		ExtendTid:     getTestExtendTid(),
		Timestamp:     time.Now().Unix(),
		TransactionId: tid.(string),
		DestCoinNum:   destCoinNum,
		DestCoinType:  destCoinType,
	}
}

func TestExchange(t *testing.T) {
	once.Do(startHTTP)

	Convey("exchange normal 先调用get接口　再调用exchange 再调用get接口　比较用户钱包数据", t, func() {
		platforms := []string{"pc", "android", "ios"}
		var num int64 = 1000
		var exchangeNum int64 = 100
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

			res = queryExchange(t, getTestExchangeForm(t, uid, "gold", exchangeNum, "silver", exchangeNum, nil), platform)
			So(res.Code, ShouldEqual, 0)

			afterExchangeWallet := getTestWallet(t, uid, platform)

			So(getIntCoinForTest(afterExchangeWallet.Gold)-getIntCoinForTest(afterWallet.Gold), ShouldEqual, exchangeNum*-1)
			So(getIntCoinForTest(afterExchangeWallet.Silver)-getIntCoinForTest(afterWallet.Silver), ShouldEqual, exchangeNum)
		}
	})
}
