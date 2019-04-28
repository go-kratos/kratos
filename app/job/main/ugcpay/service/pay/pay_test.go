package pay

import (
	"encoding/json"
	"flag"
	"net/url"
	"os"
	"testing"

	"go-common/app/service/main/ugcpay/conf"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	p *Pay
)

func TestMain(m *testing.M) {
	flag.Set("conf", "../../cmd/test.toml")
	if err := conf.Init(); err != nil {
		panic(err)
	}
	p = &Pay{
		ID:                     conf.Conf.Biz.Pay.ID,
		Token:                  conf.Conf.Biz.Pay.Token,
		RechargeShellNotifyURL: "http://api.bilibili.co/x/internal/ugcpay/trade/recharge/callback",
	}

	m.Run()
	os.Exit(0)
}

func TestCheckOrder(t *testing.T) {
	Convey("", t, func() {
		param := p.CheckOrder("3059753508505497600")
		p.Sign(param)
		t.Log(p.ToJSON(param))
	})
}

func TestCheckRefundOrder(t *testing.T) {
	Convey("", t, func() {
		param := p.CheckRefundOrder("3059753508505497600")
		p.Sign(param)
		t.Log(p.ToJSON(param))
	})
}

func TestRechargeShell(t *testing.T) {
	var (
		orderID = "123"
		mid     = int64(46333)
		assetBP = int64(1)
		shell   = int64(1)
	)
	Convey("", t, func() {
		_, json, err := p.RechargeShell(orderID, mid, assetBP, shell)
		So(err, ShouldBeNil)
		t.Log(json)
	})
}

func TestSign(t *testing.T) {
	Convey("", t, func() {
		var (
			param = url.Values{
				"customerId":      []string{"10017"},
				"deviceType":      []string{"3"},
				"notifyUrl":       []string{"http://api.bilibili.co/x/internal/ugcpay/trade/pay/callback"},
				"orderCreateTime": []string{"1539935981000"},
				"orderExpire":     []string{"1800"},
				"orderId":         []string{"224"},
				"originalAmount":  []string{"2000"},
				"payAmount":       []string{"2000"},
				"productId":       []string{"10110688"},
				"serviceType":     []string{"99"},
				"showTitle":       []string{"传点什么好呢？"},
				"timestamp":       []string{"1539935981000"},
				"traceId":         []string{"1539935981967342977"},
				"uid":             []string{"27515244"},
				"version":         []string{"1.0"},
				"feeType":         []string{"CNY"},
			}
		)
		err := p.Sign(param)
		So(err, ShouldBeNil)

		pmap := make(map[string]string)
		var payBytes []byte
		for k, v := range param {
			if len(v) > 0 {
				pmap[k] = v[0]
			}
		}
		if payBytes, err = json.Marshal(pmap); err != nil {
			return
		}
		t.Log(string(payBytes))
	})
}

func TestSignVerify(t *testing.T) {
	Convey("", t, func() {
		var (
			param = url.Values{
				"customerId":      []string{"10017"},
				"deviceType":      []string{"3"},
				"notifyUrl":       []string{"http://api.bilibili.co/x/internal/ugcpay/trade/pay/callback"},
				"orderCreateTime": []string{"1539935981000"},
				"orderExpire":     []string{"1800"},
				"orderId":         []string{"15"},
				"originalAmount":  []string{"2000"},
				"payAmount":       []string{"2000"},
				"productId":       []string{"10110688"},
				"serviceType":     []string{"99"},
				"showTitle":       []string{"传点什么好呢？"},
				"timestamp":       []string{"1539935981000"},
				"traceId":         []string{"1539935981967342977"},
				"uid":             []string{"27515244"},
				"version":         []string{"1.0"},
				"feeType":         []string{"CNY"},
			}
		)
		err := p.Sign(param)
		So(err, ShouldBeNil)

		ok := p.Verify(param)
		So(ok, ShouldBeTrue)
	})
}
