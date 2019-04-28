package pay

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"time"

	"github.com/pkg/errors"
	"go-common/library/log"
)

// Pay is.
type Pay struct {
	ID                     string
	Token                  string
	RechargeShellNotifyURL string
}

// TraceID .
func (p *Pay) TraceID() string {
	return strconv.FormatInt(time.Now().UnixNano(), 10)
}

// RechargeShellReq .
type RechargeShellReq struct {
	CustomerID  string              `json:"customerId"`
	ProductName string              `json:"productName"`
	Rate        string              `json:"rate"`
	NotifyURL   string              `json:"notifyUrl"`
	Timestamp   int64               `json:"timestamp"`
	SignType    string              `json:"signType"`
	Sign        string              `json:"sign"`
	Data        []RechargeShellData `json:"data"`
}

// RechargeShellData .
type RechargeShellData struct {
	ThirdOrderNo string `json:"thirdOrderNo"`
	MID          int64  `json:"mid"`
	ThirdCoin    string `json:"thirdCoin"`
	Brokerage    string `json:"brokerage"`
	ThirdCtime   int64  `json:"thirdCtime"`
}

// RechargeShell 转入贝壳
func (p *Pay) RechargeShell(orderID string, mid int64, assetBP int64, shell int64) (params url.Values, jsonData string, err error) {
	var (
		productName = "UGC付费"
		rate        = "1.00"
		timestamp   = time.Now().Unix() * 1000
		thirdCoin   = float64(assetBP) / 100
		brokerage   = float64(shell) / 100
	)

	params = make(url.Values)
	params.Set("customerId", p.ID)
	params.Set("productName", productName)
	params.Set("rate", rate)
	params.Set("notifyUrl", p.RechargeShellNotifyURL)
	params.Set("timestamp", strconv.FormatInt(timestamp, 10))
	params.Set("data", fmt.Sprintf("[{brokerage=%.2f&mid=%d&thirdCoin=%.2f&thirdCtime=%d&thirdOrderNo=%s}]", brokerage, mid, thirdCoin, timestamp, orderID))
	p.Sign(params)

	data := RechargeShellData{
		ThirdOrderNo: orderID,
		MID:          mid,
		ThirdCoin:    fmt.Sprintf("%.2f", thirdCoin),
		Brokerage:    fmt.Sprintf("%.2f", brokerage),
		ThirdCtime:   timestamp,
	}
	req := RechargeShellReq{
		CustomerID:  p.ID,
		ProductName: productName,
		Rate:        rate,
		NotifyURL:   p.RechargeShellNotifyURL,
		Timestamp:   timestamp,
		SignType:    params.Get("signType"),
		Sign:        params.Get("sign"),
		Data:        []RechargeShellData{data},
	}

	payBytes, err := json.Marshal(req)
	if err != nil {
		err = errors.Wrapf(err, "pay.RechargeShell.ToJSON : %s", params.Encode())
		return
	}
	jsonData = string(payBytes)
	return
}

// CheckOrder 对账param
func (p Pay) CheckOrder(txID string) (params url.Values) {
	params = make(url.Values)
	params.Set("customerId", p.ID)
	params.Set("txIds", txID)
	return
}

// CheckRefundOrder 退款对账param
func (p Pay) CheckRefundOrder(txID string) (params url.Values) {
	params = make(url.Values)
	params.Set("customerId", p.ID)
	params.Set("txIds", txID)
	return
}

// Query 返回订单查询param
func (p *Pay) Query(orderID string) (params url.Values) {
	params = make(url.Values)
	params.Set("customerId", p.ID)
	params.Set("orderIds", orderID)
	params.Set("timestamp", strconv.FormatInt(time.Now().Unix()*1000, 10))
	params.Set("traceId", p.TraceID())
	params.Set("version", "1.0")
	return
}

// ToJSON param to json
func (p *Pay) ToJSON(params url.Values) (j string, err error) {
	var (
		payBytes []byte
		pmap     = make(map[string]string)
	)
	for k, v := range params {
		if len(v) > 0 {
			pmap[k] = v[0]
		}
	}
	if payBytes, err = json.Marshal(pmap); err != nil {
		err = errors.Wrapf(err, "pay.ToJSON : %s", params.Encode())
		return
	}
	j = string(payBytes)
	return
}

// DeviceType 支付平台DeviceType
func (p *Pay) DeviceType(platform string) (t int64) {
	// 支付设备渠道类型，  1 pc 2 webapp 3 app 4jsapi 5 server 6小程序支付 7聚合二维码支付
	switch platform {
	case "ios", "android":
		return 3
	default:
		return 1
	}
}

// Sign 支付平台接口签名
func (p *Pay) Sign(params url.Values) (err error) {
	params.Set("signType", "MD5")
	sortedStr := params.Encode()
	if sortedStr, err = url.QueryUnescape(sortedStr); err != nil {
		return
	}
	b := bytes.Buffer{}
	b.WriteString(sortedStr)
	b.WriteString("&token=" + p.Token)
	signMD5 := md5.Sum(b.Bytes())
	sign := hex.EncodeToString(signMD5[:])
	params.Set("sign", sign)
	return
}

// Verify 支付平台返回param校验
func (p *Pay) Verify(params url.Values) (ok bool) {
	var (
		rs = params.Get("sign")
		s  string
	)
	ok = false
	defer func() {
		if !ok {
			params.Set("sign", rs)
			log.Error("Verify pay sign error, expect : %s, actual : %s, params : %s", s, rs, params.Encode())
		}
	}()
	if rs == "" {
		return
	}
	params.Del("sign")
	if err := p.Sign(params); err != nil {
		log.Error("Verify pay sign error : %+v", err)
		return
	}
	s = params.Get("sign")
	if rs == s {
		ok = true
		return
	}
	return
}
