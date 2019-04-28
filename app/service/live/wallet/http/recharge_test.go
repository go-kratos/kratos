package http

import (
	"fmt"
	"sync"
	"testing"

	"context"
	. "github.com/smartystreets/goconvey/convey"
	"go-common/app/service/live/wallet/conf"
	"go-common/app/service/live/wallet/dao"
	"go-common/app/service/live/wallet/model"
	"go-common/app/service/live/wallet/service"
	httpx "go-common/library/net/http/blademaster"
	"math/rand"
	"net/url"
	"strconv"
	"time"
)

const (
	_getURL      = "http://localhost:9901/x/internal/livewallet/wallet/get"
	_delCacheURL = "http://localhost:9901/x/internal/livewallet/wallet/delCache"
	_getAllURL   = "http://localhost:9901/x/internal/livewallet/wallet/getAll"
	_getTidURL   = "http://localhost:9901/x/internal/livewallet/wallet/getTid"
	_rechargeURL = "http://localhost:9901/x/internal/livewallet/wallet/recharge"
	_payURL      = "http://localhost:9901/x/internal/livewallet/wallet/pay"
	_exchangeURL = "http://localhost:9901/x/internal/livewallet/wallet/exchange"
	_queryURL    = "http://localhost:9901/x/internal/livewallet/wallet/query"
)

var (
	once   sync.Once
	client *httpx.Client
	r      *rand.Rand
)

type RechargeRes struct {
	Code int                  `json:"code"`
	Resp *model.MelonseedResp `json:"data"`
}

func getTestRandUid() int64 {
	return r.Int63n(10000000)
}

func getTestExtendTid() string {
	return fmt.Sprintf("test:ex:%d", r.Int31n(1000000))
}

func getTestRechargeOrPayForm(t *testing.T, serviceType int32, uid int64, coinType string, coinNum int64, tid interface{}) *model.RechargeOrPayForm {
	if tid == nil {
		res := queryGetTid(t, serviceType, getTestParamsJson())
		if res.Code != 0 {
			t.Errorf("get tid failed code : %d", res.Code)
			t.FailNow()
		}
		tid = res.Resp.TransactionId
	}
	return &model.RechargeOrPayForm{
		Uid:           uid,
		CoinType:      coinType,
		CoinNum:       coinNum,
		ExtendTid:     getTestExtendTid(),
		Timestamp:     time.Now().Unix(),
		TransactionId: tid.(string),
	}
}

func queryRecharge(t *testing.T, form *model.RechargeOrPayForm, platform string) *RechargeRes {
	params := url.Values{}
	params.Set("uid", fmt.Sprintf("%d", form.Uid))
	params.Set("coin_type", form.CoinType)
	params.Set("coin_num", fmt.Sprintf("%d", form.CoinNum))
	params.Set("extend_tid", form.ExtendTid)
	params.Set("timestamp", fmt.Sprintf("%d", form.Timestamp))
	params.Set("transaction_id", form.TransactionId)
	params.Set("biz_reason", "2")
	params.Set("version", "1")

	req, _ := client.NewRequest("POST", _rechargeURL, "127.0.0.1", params)
	req.Header.Set("platform", platform)

	var res RechargeRes

	err := client.Do(context.TODO(), req, &res)
	if err != nil {
		t.Errorf("client.Do() error(%v)", err)
		t.FailNow()
	}
	return &res
}

func startHTTP() {
	if err := conf.Init(); err != nil {
		panic(fmt.Errorf("conf.Init() error(%v)", err))
	}
	svr := service.New(conf.Conf)
	client = httpx.NewClient(conf.Conf.HTTPClient)
	Init(conf.Conf, svr)
	r = rand.New(rand.NewSource(time.Now().UnixNano()))

}

func getIntCoinForTest(coinStr string) int64 {
	coin, _ := strconv.Atoi(coinStr)
	return int64(coin)
}

func TestRecharge(t *testing.T) {
	Convey("recharge normal 先调用get接口　再调用recharge 再调用get接口　比较用户钱包数据", t, testWith(func() {
		platforms := []string{"pc", "android", "ios"}
		var num int64 = 1000
		uid := getTestRandUid()
		d := dao.New(conf.Conf)
		for _, platform := range platforms {

			beforeWallet := getTestWallet(t, uid, platform)

			resTid := queryGetTid(t, int32(model.RECHARGETYPE), getTestParamsJson())
			if resTid.Code != 0 {
				t.Errorf("get tid failed code : %d", resTid.Code)
				t.FailNow()
			}
			tid := resTid.Resp.TransactionId

			res := queryRecharge(t, getTestRechargeOrPayForm(t, int32(model.RECHARGETYPE), uid, "gold", num, tid), platform)

			So(res.Code, ShouldEqual, 0)
			So(getIntCoinForTest(res.Resp.Gold)-getIntCoinForTest(beforeWallet.Gold), ShouldEqual, num)

			record, err := d.GetCoinStreamByTidAndOffset(context.TODO(), tid, 0)
			So(err, ShouldBeNil)
			So(record.Reserved1, ShouldEqual, 2)
			So(record.Version, ShouldEqual, 1)

			res = queryRecharge(t, getTestRechargeOrPayForm(t, int32(model.RECHARGETYPE), uid, "silver", num, nil), platform)

			So(res.Code, ShouldEqual, 0)
			So(getIntCoinForTest(res.Resp.Silver)-getIntCoinForTest(beforeWallet.Silver), ShouldEqual, num)

			afterWallet := getTestWallet(t, uid, platform)

			So(getIntCoinForTest(afterWallet.Gold)-getIntCoinForTest(beforeWallet.Gold), ShouldEqual, num)
			So(getIntCoinForTest(afterWallet.Silver)-getIntCoinForTest(beforeWallet.Silver), ShouldEqual, num)

		}

	}))

}
