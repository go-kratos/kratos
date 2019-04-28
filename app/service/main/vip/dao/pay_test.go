package dao

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/url"
	"testing"

	"go-common/app/service/main/vip/model"
	"go-common/library/ecode"

	"github.com/smartystreets/goconvey/convey"
	"gopkg.in/h2non/gock.v1"
)

func TestDaoproductID(t *testing.T) {
	var (
		months = int16(1)
	)
	convey.Convey("productID", t, func(ctx convey.C) {
		id := d.productID(months)
		ctx.Convey("id should not be nil", func(ctx convey.C) {
			ctx.So(id, convey.ShouldEqual, _productMonthID)
		})
	})
	convey.Convey("productID == _productQuarterID", t, func(ctx convey.C) {
		months = 3
		id := d.productID(months)
		ctx.Convey("productID == 61", func(ctx convey.C) {
			ctx.So(id, convey.ShouldEqual, 61)
		})
	})
	convey.Convey("productID == _productYearID", t, func(ctx convey.C) {
		months = _yearMonths
		id := d.productID(months)
		ctx.Convey("productID == 61", func(ctx convey.C) {
			ctx.So(id, convey.ShouldEqual, _productYearID)
		})
	})
}

func TestDaoPayWallet(t *testing.T) {
	var (
		c    = context.TODO()
		mid  = int64(20606508)
		ip   = ""
		data = &model.PayAccountResp{}
	)
	convey.Convey("PayWallet", t, func(ctx convey.C) {
		defer gock.OffAll()
		httpMock("POST", _payWallet).Reply(200).JSON(`{"code":0}`)
		err := d.PayWallet(c, mid, ip, data)
		ctx.Convey("Error should be nil", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoPayBanks(t *testing.T) {
	var (
		c    = context.TODO()
		ip   = ""
		data []*model.PayBankResp
	)
	convey.Convey("PayBanks", t, func(ctx convey.C) {
		defer gock.OffAll()
		httpMock("GET", _payBanks).Reply(200).JSON(`{"code":0}`)
		err := d.PayBanks(c, ip, data)
		ctx.Convey("Error should be nil", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoAddPayOrder(t *testing.T) {
	var (
		c  = context.TODO()
		ip = ""
		o  = &model.PayOrder{
			OrderNo: "201808141642435698853709",
			Money:   148,
		}
		data = &model.AddPayOrderResp{}
	)
	convey.Convey("AddPayOrder", t, func(ctx convey.C) {
		defer gock.OffAll()
		httpMock("POST", _addPayOrder).Reply(200).JSON(`{"code":0}`)
		err := d.AddPayOrder(c, ip, o, data)
		ctx.Convey("Error should be nil", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoPaySDK(t *testing.T) {
	var (
		c       = context.TODO()
		ip      = ""
		o       = &model.PayOrder{}
		data    = &model.APIPayOrderResp{}
		payCode = ""
	)
	convey.Convey("PaySDK", t, func(ctx convey.C) {
		defer gock.OffAll()
		httpMock("POST", _paySDK).Reply(200).JSON(`{"code":0}`)
		err := d.PaySDK(c, ip, o, data, payCode)
		ctx.Convey("Error should be nil", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoPayQrcode(t *testing.T) {
	var (
		c  = context.TODO()
		ip = "120.11.188.230"
		o  = &model.PayOrder{
			OrderNo: "2014100120411262876516",
			Mid:     4908640,
			Money:   148.00,
		}
		data    = &model.APIPayOrderResp{}
		payCode = "alipay"
	)
	convey.Convey("PayQrcode", t, func(ctx convey.C) {
		defer gock.OffAll()
		httpMock("POST", _payQrcode).Reply(200).JSON(`{"code":0}`)
		err := d.PayQrcode(c, ip, o, data, payCode)
		ctx.Convey("Error should be nil", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoQuickPayToken(t *testing.T) {
	var (
		c         = context.TODO()
		ip        = ""
		accessKey = ""
		cookie    []*http.Cookie
		data      = &model.QucikPayResp{}
	)
	convey.Convey("QuickPayToken", t, func(ctx convey.C) {
		defer gock.OffAll()
		httpMock("POST", d.c.Property.PayURL+_quickPayToken).Reply(200).JSON(`{"code":0}`)
		err := d.QuickPayToken(c, ip, accessKey, cookie, data)
		ctx.Convey("Error should be nil", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoDoQuickPay(t *testing.T) {
	var (
		c            = context.TODO()
		ip           = ""
		token        = ""
		thirdTradeNo = ""
		data         = &model.PayRetResp{}
	)
	convey.Convey("DoQuickPay", t, func(ctx convey.C) {
		defer gock.OffAll()
		httpMock("POST", _quickPayDo).Reply(200).JSON(`{"code":0}`)
		err := d.DoQuickPay(c, ip, token, thirdTradeNo, data)
		ctx.Convey("Error should be nil", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoPayCashier(t *testing.T) {
	var (
		c        = context.TODO()
		ip       = ""
		o        = &model.PayOrder{}
		data     = &model.APIPayOrderResp{}
		payCode  = ""
		bankCode = ""
	)
	convey.Convey("PayCashier", t, func(ctx convey.C) {
		defer gock.OffAll()
		httpMock("POST", _payCashier).Reply(200).JSON(`{"code":0}`)
		err := d.PayCashier(c, ip, o, data, payCode, bankCode)
		ctx.Convey("Error should be nil", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoPayIapAccess(t *testing.T) {
	var (
		c  = context.TODO()
		ip = ""
		o  = &model.PayOrder{
			OrderNo:      "123",
			AppID:        1,
			Platform:     1,
			OrderType:    1,
			AppSubID:     "456",
			Mid:          20606508,
			ToMid:        20606509,
			BuyMonths:    1,
			Money:        3.0,
			Status:       1,
			PayType:      1,
			RechargeBp:   1.0,
			ThirdTradeNo: "123",
			Ver:          1,
		}
		data      = &model.APIPayOrderResp{}
		productID = "123"
	)
	convey.Convey("PayIapAccess", t, func(ctx convey.C) {
		defer gock.OffAll()
		httpMock("POST", _payIapAccess).Reply(200).JSON(`{"code":0}`)
		err := d.PayIapAccess(c, ip, o, data, productID)
		ctx.Convey("Error should be nil", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoPayClose(t *testing.T) {
	var (
		c       = context.TODO()
		orderNO = "123456"
		ip      = "127.0.0.1"
	)
	convey.Convey("PayClose", t, func(ctx convey.C) {
		defer gock.OffAll()
		httpMock("POST", d.payCloseURL).Reply(200).JSON(`{"code":0}`)
		_, err := d.PayClose(c, orderNO, ip)
		ctx.Convey("Error should be nil", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaodopay(t *testing.T) {
	var (
		c       = context.TODO()
		urlPath = _payWallet
		ip      = "172.18.35.12"
		params  = url.Values{}
		data    = &model.PayAccountResp{}
		mid     = int64(20606508)
	)
	convey.Convey("dopay", t, func(ctx convey.C) {
		params.Add("mid", fmt.Sprintf("%d", mid))
		defer gock.OffAll()
		httpMock("POST", _payWallet).Reply(200).JSON(`{"code":0}`)
		err := d.dopay(c, urlPath, ip, params, data, d.client.Post)
		ctx.Convey("Error should be nil", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoPayRecission(t *testing.T) {
	var (
		c       = context.TODO()
		params  map[string]interface{}
		clietIP = ""
	)
	convey.Convey("PayRecission", t, func(ctx convey.C) {
		defer gock.OffAll()
		httpMock("POST", _payWallet).Reply(200).JSON(`{"code":0}`)
		err := d.PayRecission(c, params, clietIP)
		ctx.Convey("Error should be nil", func(ctx convey.C) {
			ctx.So(err, convey.ShouldEqual, ecode.VipRescissionErr)
		})
	})
}

func TestDaoPaySign(t *testing.T) {
	var (
		params map[string]string
		token  = ""
	)
	convey.Convey("PaySign", t, func(ctx convey.C) {
		sign := d.PaySign(params, token)
		ctx.Convey("sign should not be nil", func(ctx convey.C) {
			ctx.So(sign, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoPaySignNotDel(t *testing.T) {
	var (
		params map[string]string
		token  = ""
	)
	convey.Convey("PaySignNotDel", t, func(ctx convey.C) {
		sign := d.PaySignNotDel(params, token)
		ctx.Convey("sign should not be nil", func(ctx convey.C) {
			ctx.So(sign, convey.ShouldNotBeNil)
		})
	})
}

func TestDaosortParamsKey(t *testing.T) {
	var (
		v map[string]string
	)
	convey.Convey("sortParamsKey", t, func(ctx convey.C) {
		p1 := d.sortParamsKey(v)
		ctx.Convey("p1 should not be nil", func(ctx convey.C) {
			ctx.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func TestDaodoPaySend(t *testing.T) {
	var (
		c        = context.TODO()
		basePath = d.c.Property.PayURL
		path     = _quickPayToken
		IP       = ""
		cookie   []*http.Cookie
		header   map[string]string
		method   = http.MethodPost
		params   = map[string]string{"name": "vip"}
		data     = new(struct {
			Code int64 `json:"errno"`
			Data int64 `json:"msg"`
		})
	)
	convey.Convey("doPaySend", t, func(ctx convey.C) {
		defer gock.OffAll()
		httpMock("POST", basePath+path).Reply(200).JSON(`{"code":0}`)
		err := d.doPaySend(c, basePath, path, IP, cookie, header, method, params, data, d.PaySign)
		ctx.Convey("Error should be nil", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoreadAll(t *testing.T) {
	var (
		r        = bytes.NewReader([]byte{})
		capacity = int64(0)
	)
	convey.Convey("readAll", t, func(ctx convey.C) {

		b, err := readAll(r, capacity)
		ctx.Convey("Error should be nil", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
		ctx.Convey("b should not be nil", func(ctx convey.C) {
			ctx.So(b, convey.ShouldNotBeNil)
		})
	})
}

// go test  -test.v -test.run TestDaoPayClose
func TestDaoPayQrCode(t *testing.T) {
	convey.Convey("TestDaoPayQrCode", t, func() {
		res, err := d.PayQrCode(context.TODO(), 1, "1", make(map[string]interface{}))
		convey.So(err, convey.ShouldBeNil)
		convey.So(res, convey.ShouldNotBeNil)
	})
}
