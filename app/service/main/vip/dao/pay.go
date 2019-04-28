package dao

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"

	"go-common/app/service/main/vip/model"
	"go-common/library/ecode"
	"go-common/library/log"

	"github.com/pkg/errors"
)

const (
	//pay url
	_addPayOrder     = "/api/add.pay.order"
	_paySDK          = "/api/pay.sdk"
	_quickPayToken   = "/api/quick.pay.token"
	_quickPayDo      = "/api/quick.pay.do"
	_payQrcode       = "/api/pay.qrcode"
	_payCashier      = "/api/pay.cashier"
	_payIapAccess    = "/api/add.iap.access"
	_payBanks        = "/api/pay.banks"
	_payWallet       = "/api/v2/user/account"
	_payRescission   = "/payplatform/pay/rescission"
	_payClose        = "/payplatform/pay/cancel"
	_createPayQrCode = "/payplatform/qrcode/createPayQrCode"

	_payRemark        = "购买大会员服务%d个月"
	_merchantID       = 17
	_productMonthID   = 60
	_productQuarterID = 61
	_productYearID    = 62
	_quarterMonths    = 3
	_yearMonths       = 12
	_iapQuantity      = "1"
	_minRead          = 1024 * 64

	_retry       = 3
	_defversion  = "1.0"
	_defsigntype = "MD5"
)

//BasicResp pay response.
type BasicResp struct {
	Code    int          `json:"code"`
	Message string       `json:"messge"`
	Data    *interface{} `json:"data"`
}

// merchantProduct get Product id by months
func (d *Dao) productID(months int16) (id int8) {
	id = _productMonthID
	if months >= _quarterMonths && months < _yearMonths {
		id = _productQuarterID
	} else if months >= _yearMonths {
		id = _productYearID
	}
	return
}

//PayWallet .
func (d *Dao) PayWallet(c context.Context, mid int64, ip string, data *model.PayAccountResp) (err error) {
	val := url.Values{}
	val.Add("mid", fmt.Sprintf("%d", mid))
	return d.dopay(c, _payWallet, ip, val, data, d.client.Post)
}

//PayBanks .
func (d *Dao) PayBanks(c context.Context, ip string, data []*model.PayBankResp) (err error) {
	val := url.Values{}
	return d.dopay(c, _payBanks, ip, val, data, d.client.Get)
}

//AddPayOrder add pay order.
func (d *Dao) AddPayOrder(c context.Context, ip string, o *model.PayOrder, data *model.AddPayOrderResp) (err error) {
	val := url.Values{}
	val.Add("mid", fmt.Sprintf("%d", o.Mid))
	val.Add("to_mid", fmt.Sprintf("%d", o.ToMid))
	val.Add("money", fmt.Sprintf("%f", o.Money))
	val.Add("remark", fmt.Sprintf(_payRemark, o.BuyMonths))
	val.Add("subject", fmt.Sprintf(_payRemark, o.BuyMonths))
	val.Add("out_trade_no", o.OrderNo)
	val.Add("notify_url", d.c.Property.NotifyURL)
	val.Add("merchant_id", fmt.Sprintf("%d", _merchantID))
	val.Add("merchant_product_id", fmt.Sprintf("%d", d.productID(o.BuyMonths)))
	var platform string
	switch o.Platform {
	case model.DeviceIOS:
		platform = "1"
	case model.DeviceANDROID:
		platform = "2"
	default:
		platform = "3"
	}
	val.Add("platform_type", platform)
	return d.dopay(c, _addPayOrder, ip, val, data, d.client.Post)
}

// PaySDK moblie pay sdk.
func (d *Dao) PaySDK(c context.Context, ip string, o *model.PayOrder, data *model.APIPayOrderResp, payCode string) (err error) {
	val := url.Values{}
	val.Add("money", fmt.Sprintf("%f", o.Money))
	val.Add("pay_order_no", o.OrderNo)
	val.Add("pay_type", payCode)
	return d.dopay(c, _paySDK, ip, val, data, d.client.Post)
}

// PayQrcode pay qrcode.
func (d *Dao) PayQrcode(c context.Context, ip string, o *model.PayOrder, data *model.APIPayOrderResp, payCode string) (err error) {
	val := url.Values{}
	val.Add("money", fmt.Sprintf("%f", o.Money))
	val.Add("pay_order_no", o.OrderNo)
	val.Add("pay_type", payCode)

	return d.dopay(c, _payQrcode, ip, val, data, d.client.Post)
}

// QuickPayToken quick pay token.
func (d *Dao) QuickPayToken(c context.Context, ip string, accessKey string, cookie []*http.Cookie, data *model.QucikPayResp) (err error) {
	params := make(map[string]string)
	params["access_key"] = accessKey

	//return d.dopay(c, _quickPayToken, ip, val, cookie, data, d.client.Post)
	return d.doPaySend(c, d.c.Property.PayURL, _quickPayToken, ip, cookie, nil, http.MethodPost, params, data, d.PaySign)
}

// DoQuickPay do quick pay.
func (d *Dao) DoQuickPay(c context.Context, ip string, token string, thirdTradeNo string, data *model.PayRetResp) (err error) {
	val := url.Values{}
	val.Add("token", token)
	val.Add("pay_order_no", thirdTradeNo)
	return d.dopay(c, _quickPayDo, ip, val, data, d.client.Post)
}

// PayCashier pay cashier.
func (d *Dao) PayCashier(c context.Context, ip string, o *model.PayOrder, data *model.APIPayOrderResp, payCode string, bankCode string) (err error) {
	val := url.Values{}
	val.Add("pay_order_no", o.ThirdTradeNo)
	val.Add("money", fmt.Sprintf("%f", o.Money))
	val.Add("pay_type", payCode)
	val.Add("bank_code", bankCode)
	return d.dopay(c, _payCashier, ip, val, data, d.client.Post)
}

// PayIapAccess pay iap access.
func (d *Dao) PayIapAccess(c context.Context, ip string, o *model.PayOrder, data *model.APIPayOrderResp, productID string) (err error) {
	val := url.Values{}
	val.Add("mid", fmt.Sprintf("%d", o.Mid))
	val.Add("product_id", productID)
	val.Add("quantity", _iapQuantity)
	val.Add("money", fmt.Sprintf("%f", o.Money))
	val.Add("remark", fmt.Sprintf(_payRemark, o.BuyMonths))
	val.Add("pay_order_no", o.ThirdTradeNo)
	val.Add("merchant_id", fmt.Sprintf("%d", _merchantID))
	val.Add("merchant_product_id", fmt.Sprintf("%d", d.productID(o.BuyMonths)))
	return d.dopay(c, _payIapAccess, ip, val, data, d.client.Post)
}

//PayClose pay close.
func (d *Dao) PayClose(c context.Context, orderNO string, ip string) (data *model.APIPayCancelResp, err error) {
	params := make(map[string]string)
	params["customerId"] = strconv.FormatInt(d.c.PayConf.CustomerID, 10)
	params["orderId"] = orderNO
	params["timestamp"] = fmt.Sprintf("%d", time.Now().Unix()*1000)
	params["traceId"] = model.UUID4()
	params["version"] = _defversion
	params["signType"] = _defsigntype
	sign := d.PaySign(params, d.c.PayConf.Token)
	params["sign"] = sign
	resp := new(
		struct {
			Code int64                   `json:"errno"`
			Data *model.APIPayCancelResp `json:"data"`
		})
	header := make(map[string]string)
	header["Content-Type"] = "application/json"
	marshal, _ := json.Marshal(params)
	if err = d.doSend(c, d.payCloseURL, "127.0.0.1", header, marshal, resp); err != nil {
		err = errors.Wrapf(err, "Call pay service(%s)", d.payCloseURL+"?"+string(marshal))
		return
	}
	if resp.Code != 0 {
		err = fmt.Errorf("Call pay service(%s) error, response code is not 0, resp:%v", d.payCloseURL+"?"+string(marshal), resp)
		return
	}
	data = resp.Data
	log.Info("Call pay service(%s) successful, resp(%v)", d.payCloseURL+"?"+string(marshal), data)
	return
}

//PayQrCode pay qr code.
func (d *Dao) PayQrCode(c context.Context, mid int64, orderID string, req map[string]interface{}) (data *model.PayQrCode, err error) {
	payParam, _ := json.Marshal(req)
	params := make(map[string]string)
	params["customerId"] = strconv.FormatInt(d.c.PayConf.CustomerID, 10)
	params["uid"] = fmt.Sprintf("%d", mid)
	params["orderId"] = orderID
	params["payParam"] = string(payParam)
	params["timestamp"] = fmt.Sprintf("%d", time.Now().Unix()*1000)
	params["traceId"] = model.UUID4()
	params["version"] = d.c.PayConf.Version
	params["signType"] = d.c.PayConf.SignType
	sign := d.PaySign(params, d.c.PayConf.Token)
	params["sign"] = sign
	resp := new(
		struct {
			Code int64            `json:"errno"`
			Data *model.PayQrCode `json:"data"`
		})
	url := d.c.Property.PayCoURL + _createPayQrCode
	header := make(map[string]string)
	header["Content-Type"] = "application/json"
	marshal, _ := json.Marshal(params)
	if err = d.doSend(c, url, "127.0.0.1", header, marshal, resp); err != nil {
		err = errors.Wrapf(err, "Call pay service(%s)", url+"?"+string(marshal))
		return
	}
	if resp.Code != 0 {
		err = fmt.Errorf("Call pay service(%s) error, response code is not 0, resp:%v", url+"?"+string(marshal), resp)
		return
	}
	data = resp.Data
	log.Info("Call pay service(%s) successful, resp(%v)", url+"?"+string(marshal), data)
	return
}

func (d *Dao) doSend(c context.Context, url, IP string, header map[string]string, marshal []byte, data interface{}) (err error) {
	var (
		req    *http.Request
		client = new(http.Client)
		resp   *http.Response
		bs     []byte
	)
	if req, err = http.NewRequest(http.MethodPost, url, strings.NewReader(string(marshal))); err != nil {
		err = errors.WithStack(err)
		return
	}
	for k, v := range header {
		req.Header.Add(k, v)
	}
	if resp, err = client.Do(req); err != nil {
		err = errors.Wrapf(err, "Call pay service(%s)", url+"?"+string(marshal))
		return
	}
	defer resp.Body.Close()

	defer func() {
		log.Info("call url:%v params:(%v) result:(%+v) header:(%+v)", url, string(marshal), data, header)
	}()
	if resp.StatusCode >= http.StatusBadRequest {
		err = errors.Errorf("incorrect http status:%d host:%s, url:%s", resp.StatusCode, req.URL.Host, req.URL.String())
		return
	}
	if bs, err = readAll(resp.Body, _minRead); err != nil {
		err = errors.Wrapf(err, "host:%s, url:%s", req.URL.Host, req.URL.String())
		return
	}
	if err = json.Unmarshal(bs, data); err != nil {
		err = errors.WithStack(err)
		return
	}
	return
}

func (d *Dao) dopay(c context.Context,
	urlPath string, ip string, params url.Values, data interface{},
	fn func(c context.Context, uri string, ip string, params url.Values, r interface{}) error,
) (err error) {
	var (
		resp    = &BasicResp{}
		urlAddr string
	)
	urlAddr = d.c.Property.PayURL + urlPath
	err = fn(c, urlAddr, ip, params, resp)
	if err != nil {
		err = errors.Wrapf(err, "Call pay service(%s)", urlAddr+"?"+params.Encode())
		return
	}
	if resp.Code != 0 {
		err = fmt.Errorf("Call pay service(%s) error, response code is not 0, resp:%v,%v", urlAddr+"?"+params.Encode(), resp, data)
		err = errors.WithStack(err)
		return
	}
	data = resp.Data
	log.Info("Call pay service(%s) successful, resp: %v", urlAddr+"?"+params.Encode(), resp)
	return
}

// PayRecission call pay refund api.
func (d *Dao) PayRecission(c context.Context, params map[string]interface{}, clietIP string) (err error) {
	rel := new(struct {
		Code int64 `json:"errno"`
		Data int64 `json:"msg"`
	})
	paramsObj := make(map[string]string)
	for k, v := range params {
		paramsObj[k] = fmt.Sprintf("%v", v)
	}
	header := make(map[string]string)
	header["Content-Type"] = "application/json"

	defer func() {
		log.Info("url:%v params:%+v return:%+v error(%+v)", d.c.Property.PayURL+_payRescission, paramsObj, rel, err)
	}()

	success := false
	for i := 0; i < _retry; i++ {

		if err = d.doPaySend(c, d.c.Property.PayURL, _payRescission, clietIP, nil, header, http.MethodPost, paramsObj, rel, d.PaySignNotDel); err != nil {
			log.Error(" dopaysend(url:%v,params:%+v) return(%+v) ", d.c.Property.PayURL+_payRescission, paramsObj, rel)
			continue
		}
		if rel.Code == int64(ecode.OK.Code()) {
			success = true
			break
		}
	}
	if !success {
		err = ecode.VipRescissionErr
	}
	return
}

// PaySign pay sign.
func (d *Dao) PaySign(params map[string]string, token string) (sign string) {
	delete(params, "payChannelId")
	delete(params, "payChannel")
	delete(params, "accessKey")
	delete(params, "sdkVersion")
	delete(params, "openId")
	delete(params, "sign")
	delete(params, "device")

	tmp := d.sortParamsKey(params)

	var b bytes.Buffer
	b.WriteString(tmp)
	b.WriteString(fmt.Sprintf("&token=%s", token))
	log.Info("pay sign params:(%s) \n", b.String())
	mh := md5.Sum(b.Bytes())
	// query
	sign = hex.EncodeToString(mh[:])
	log.Info("pay sign (%v)", sign)
	return
}

//PaySignNotDel pay sign not del.
func (d *Dao) PaySignNotDel(params map[string]string, token string) (sign string) {

	tmp := d.sortParamsKey(params)

	var b bytes.Buffer
	b.WriteString(tmp)
	b.WriteString(fmt.Sprintf("&token=%s", token))
	log.Info("pay sign params:(%s) \n", b.String())
	mh := md5.Sum(b.Bytes())
	// query
	sign = hex.EncodeToString(mh[:])
	log.Info("pay sign (%v)", sign)
	return
}

func (d *Dao) sortParamsKey(v map[string]string) string {
	if v == nil {
		return ""
	}

	var buf bytes.Buffer
	keys := make([]string, 0, len(v))
	for k := range v {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		vs := v[k]
		prefix := k + "="
		if buf.Len() > 0 {
			buf.WriteByte('&')
		}
		buf.WriteString(prefix)
		buf.WriteString(vs)
	}
	return buf.String()
}

func (d *Dao) doPaySend(c context.Context, basePath, path, IP string, cookie []*http.Cookie, header map[string]string, method string, params map[string]string, data interface{}, signFn func(params map[string]string, token string) string) (err error) {
	var (
		req    *http.Request
		client = new(http.Client)
		resp   *http.Response
		bs     []byte
	)
	url := basePath + path
	sign := signFn(params, d.c.PayConf.Token)
	params["sign"] = sign
	marshal, _ := json.Marshal(params)
	if req, err = http.NewRequest(method, url, strings.NewReader(string(marshal))); err != nil {
		err = errors.WithStack(err)
		return
	}
	for _, v := range cookie {
		req.AddCookie(v)
	}
	for k, v := range header {
		req.Header.Add(k, v)
	}
	if resp, err = client.Do(req); err != nil {
		log.Error("call url:%v params:(%+v)", basePath+path, params)
		err = errors.WithStack(err)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode >= http.StatusBadRequest {
		err = errors.Errorf("incorrect http status:%d host:%s, url:%s", resp.StatusCode, req.URL.Host, req.URL.String())
		return
	}
	if bs, err = readAll(resp.Body, _minRead); err != nil {
		err = errors.Wrapf(err, "host:%s, url:%s", req.URL.Host, req.URL.String())
		return
	}
	if err = json.Unmarshal(bs, data); err != nil {
		err = errors.WithStack(err)
		return
	}
	log.Info("call url:%v params:%+v result:%+v", url, params, data)
	return
}

func readAll(r io.Reader, capacity int64) (b []byte, err error) {
	buf := bytes.NewBuffer(make([]byte, 0, capacity))
	// If the buffer overflows, we will get bytes.ErrTooLarge.
	// Return that as an error. Any other panic remains.
	defer func() {
		e := recover()
		if e == nil {
			return
		}
		if panicErr, ok := e.(error); ok && panicErr == bytes.ErrTooLarge {
			err = panicErr
		} else {
			panic(e)
		}
	}()
	_, err = buf.ReadFrom(r)
	return buf.Bytes(), err
}
