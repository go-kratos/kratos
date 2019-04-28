package telecom

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"go-common/app/interface/main/app-wall/model/telecom"
	"go-common/library/ecode"
	"go-common/library/log"

	"github.com/xxtea/xxtea-go/xxtea"
)

const (
	_telecomAppsecret  = "FD9B667503E74DDDBcF28E2327F88EEA"
	_telecomAppID      = "1000000053"
	_telecomFormat     = "json"
	_telecomClientType = "3"
	_flowPackageID     = 279
	_payInfo           = "/api/v1/0/getPayInfo.do"
	_cancelRepeatOrder = "/api/v1/0/cancelRepeatOrder.do"
	_sucOrderList      = "/api/v1/0/getSucOrderList.do"
	_phoneArea         = "/api/v1/0/queryOperatorAndProvince.do"
	_orderState        = "/api/v1/0/queryOrderStatus.do"
)

// PayInfo
func (d *Dao) PayInfo(c context.Context, requestNo int64, phone, isRepeatOrder, payChannel, payAction int, orderID int64, ipStr string,
	beginTime, firstOrderEndtime time.Time) (data *telecom.Pay, err error, msg string) {
	var payChannelStr string
	switch payChannel {
	case 1:
		payChannelStr = "31"
		ipStr = ""
	case 2:
		payChannelStr = "29"
	}
	params := url.Values{}
	params.Set("requestNo", strconv.FormatInt(requestNo, 10))
	params.Set("flowPackageId", strconv.Itoa(_flowPackageID))
	params.Set("contractId", "100174")
	params.Set("activityId", "101043")
	params.Set("phoneId", strconv.Itoa(phone))
	params.Set("bindApps", "tv.danmaku.bilianime|tv.danmaku.bili")
	params.Set("bindAppNames", "哔哩哔哩|哔哩哔哩")
	params.Set("isRepeatOrder", strconv.Itoa(isRepeatOrder))
	params.Set("payChannel", payChannelStr)
	if ipStr != "" {
		params.Set("userIp", ipStr)
	}
	params.Set("payPageType", "1")
	if d.telecomReturnURL != "" {
		params.Set("returnUrl", d.telecomReturnURL)
	}
	if d.telecomCancelPayURL != "" {
		params.Set("cancelPayUrl", d.telecomCancelPayURL)
	}
	params.Set("payAction", strconv.Itoa(payAction))
	// if startTime := beginTime.Format("20060102"); startTime != "19700101" && !beginTime.IsZero() {
	// 	params.Set("beginTime", startTime)
	// }
	// if endTime := firstOrderEndtime.Format("20060102"); endTime != "19700101" && !firstOrderEndtime.IsZero() {
	// 	params.Set("firstOrderEndtime", endTime)
	// }
	if orderID > 0 {
		params.Set("orderId", strconv.FormatInt(orderID, 10))
	}
	var res struct {
		Code   int `json:"resCode"`
		Detail struct {
			OrderID int64 `json:"orderId"`
			PayInfo struct {
				PayURL string `json:"payUrl"`
			} `json:"payInfo"`
		} `json:"detail"`
		Msg string `json:"resMsg"`
	}
	if err = d.wallHTTPPost(c, d.payInfoURL, params, &res); err != nil {
		log.Error("telecom_payInfoURL url(%s) error(%v)", d.payInfoURL+"?"+params.Encode(), err)
		return
	}
	if res.Code != 10000 {
		err = ecode.Int(res.Code)
		log.Error("telecom_url(%s) res code(%d) or res.data(%v)", d.payInfoURL+"?"+params.Encode(), res.Code, res.Detail)
		msg = res.Msg
		return
	}
	data = &telecom.Pay{
		RequestNo: requestNo,
		OrderID:   res.Detail.OrderID,
		PayURL:    res.Detail.PayInfo.PayURL,
	}
	return
}

// wallHTTPPost
func (d *Dao) wallHTTPPost(c context.Context, urlStr string, params url.Values, res interface{}) (err error) {
	newParams := url.Values{}
	encryptData := xxtea.Encrypt([]byte(params.Encode()), []byte(_telecomAppsecret))
	hexStr := hex.EncodeToString(encryptData)
	newParams.Set("paras", hexStr)
	mh := md5.Sum([]byte(_telecomAppID + _telecomClientType + _telecomFormat + hexStr + _telecomAppsecret))
	newParams.Set("sign", hex.EncodeToString(mh[:]))
	newParams.Set("appId", _telecomAppID)
	newParams.Set("clientType", _telecomClientType)
	newParams.Set("format", _telecomFormat)
	req, err := http.NewRequest("POST", urlStr, strings.NewReader(newParams.Encode()))
	if err != nil {
		log.Error("telecom_http.NewRequest url(%s) error(%v)", urlStr+"?"+newParams.Encode(), err)
		return
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("X-BACKEND-BILI-REAL-IP", "")
	return d.client.Do(c, req, &res)
}

// CancelRepeatOrder
func (d *Dao) CancelRepeatOrder(c context.Context, phone int, signNo string) (msg string, err error) {
	params := url.Values{}
	params.Set("phoneId", strconv.Itoa(phone))
	params.Set("signNo", signNo)
	var res struct {
		Code int    `json:"resCode"`
		Msg  string `json:"resMsg"`
	}
	if err = d.wallHTTPPost(c, d.cancelRepeatOrderURL, params, &res); err != nil {
		log.Error("telecom_payInfoURL url(%s) error(%v)", d.cancelRepeatOrderURL+"?"+params.Encode(), err)
		return
	}
	if res.Code != 10000 {
		err = ecode.Int(res.Code)
		log.Error("telecom_url(%s) res code(%d)", d.cancelRepeatOrderURL+"?"+params.Encode(), res.Code)
		msg = res.Msg
		return
	}
	msg = res.Msg
	return
}

// SucOrderList user order list
func (d *Dao) SucOrderList(c context.Context, phone int) (res *telecom.SucOrder, err error, msg string) {
	params := url.Values{}
	params.Set("phoneId", strconv.Itoa(phone))
	var resData struct {
		Code   int `json:"resCode"`
		Detail struct {
			AccessToken string              `json:"accessToken"`
			Orders      []*telecom.SucOrder `json:"orders"`
		} `json:"detail"`
		Msg string `json:"resMsg"`
	}
	if err = d.wallHTTPPost(c, d.sucOrderListURL, params, &resData); err != nil {
		log.Error("telecom_sucOrderListURL url(%s) error(%v)", d.sucOrderListURL+"?"+params.Encode(), err)
		return
	}
	if resData.Code != 10000 {
		err = ecode.Int(resData.Code)
		log.Error("telecom_url(%s) res code(%d)", d.sucOrderListURL+"?"+params.Encode(), resData.Code)
		msg = resData.Msg
		return
	}
	if len(resData.Detail.Orders) == 0 {
		err = ecode.NothingFound
		msg = "订单不存在"
		log.Error("telecom_order list phone(%v) len 0", phone)
		return
	}
	for _, r := range resData.Detail.Orders {
		if r.FlowPackageID == strconv.Itoa(_flowPackageID) {
			r.OrderID, _ = strconv.ParseInt(r.OrderIDStr, 10, 64)
			r.OrderIDStr = ""
			r.PortInt, _ = strconv.Atoi(r.Port)
			r.Port = ""
			res = r
			res.AccessToken = resData.Detail.AccessToken
			break
		}
	}
	if res == nil {
		log.Error("telecom_order bili phone(%v) is null", phone)
		msg = "订单不存在"
		err = ecode.NothingFound
		return
	}
	msg = resData.Msg
	return
}

// PhoneArea phone by area
func (d *Dao) PhoneArea(c context.Context, phone int) (area string, err error, msg string) {
	params := url.Values{}
	params.Set("phoneId", strconv.Itoa(phone))
	var resData struct {
		Code   int `json:"resCode"`
		Detail struct {
			RegionCode string `json:"regionCode"`
			AreaName   string `json:"areaName"`
		} `json:"detail"`
		Msg string `json:"resMsg"`
	}
	if err = d.wallHTTPPost(c, d.phoneAreaURL, params, &resData); err != nil {
		log.Error("telecom_phoneAreaURL url(%s) error(%v)", d.phoneAreaURL+"?"+params.Encode(), err)
		return
	}
	if resData.Code != 10000 {
		err = ecode.Int(resData.Code)
		log.Error("telecom_url(%s) res code(%d)", d.phoneAreaURL+"?"+params.Encode(), resData.Code)
		msg = resData.Msg
		return
	}
	area = resData.Detail.RegionCode
	return
}

// OrderState
func (d *Dao) OrderState(c context.Context, orderid int64) (res *telecom.OrderPhoneState, err error) {
	params := url.Values{}
	params.Set("orderId", strconv.FormatInt(orderid, 10))
	var resData struct {
		Code   int                      `json:"resCode"`
		Detail *telecom.OrderPhoneState `json:"detail"`
		Msg    string                   `json:"resMsg"`
	}
	if err = d.wallHTTPPost(c, d.orderStateURL, params, &resData); err != nil {
		log.Error("telecom_orderStateURL url(%s) error(%v)", d.orderStateURL+"?"+params.Encode(), err)
		return
	}
	if resData.Code != 10000 && resData.Code != 10013 {
		err = ecode.Int(resData.Code)
		log.Error("telecom_url(%s) res code(%d)", d.orderStateURL+"?"+params.Encode(), resData.Code)
		return
	}
	if resData.Code == 10013 {
		res = &telecom.OrderPhoneState{
			OrderState: 6,
		}
		return
	}
	res = resData.Detail
	if res.FlowPackageID != _flowPackageID {
		res.OrderState = 5
		return
	}
	switch resData.Detail.OrderState {
	case 5:
		res.OrderState = 6
	}
	return
}
