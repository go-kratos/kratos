package pay

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"net/url"
	"strconv"
	"time"

	"go-common/library/log"

	"github.com/pkg/errors"
)

// Pay is.
type Pay struct {
	ID              string
	Token           string
	OrderTTL        int
	NotifyURL       string
	RefundNotifyURL string
}

// Create 返回订单创建param
func (p *Pay) Create(orderID string, productID int64, price int64, deviceType int64, serviceType int, mid int64, title string) (params url.Values) {
	params = make(url.Values)
	params.Set("customerId", p.ID)
	params.Set("serviceType", strconv.Itoa(serviceType))
	params.Set("orderId", orderID)
	params.Set("orderCreateTime", strconv.FormatInt(time.Now().Unix()*1000, 10))
	params.Set("orderExpire", strconv.Itoa(p.OrderTTL))
	params.Set("payAmount", strconv.FormatInt(price, 10))
	params.Set("originalAmount", strconv.FormatInt(price, 10))
	params.Set("deviceType", strconv.FormatInt(deviceType, 10))
	params.Set("notifyUrl", p.NotifyURL)
	params.Set("productId", strconv.FormatInt(productID, 10))
	params.Set("showTitle", title)
	params.Set("traceId", p.TraceID())
	params.Set("timestamp", strconv.FormatInt(time.Now().Unix()*1000, 10))
	params.Set("version", "1.0")
	params.Set("uid", strconv.FormatInt(mid, 10))
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

// TraceID 生成traceID
func (p *Pay) TraceID() string {
	return strconv.FormatInt(time.Now().UnixNano(), 10)
}

// Refund 原路返回退款params
func (p *Pay) Refund(txID string, refundFee int64, refundDesc string) (params url.Values) {
	params = make(url.Values)
	params.Set("customerId", p.ID)
	params.Set("txId", txID)
	params.Set("totalAmount", strconv.FormatInt(refundFee, 10))
	params.Set("refundAmount", strconv.FormatInt(refundFee, 10))
	params.Set("refundDesc", refundDesc)
	params.Set("notifyUrl", p.RefundNotifyURL)
	params.Set("version", "1.0")
	params.Set("traceId", p.TraceID())
	params.Set("timestamp", strconv.FormatInt(time.Now().Unix()*1000, 10))
	return
}

// Cancel 返回订单取消param
func (p *Pay) Cancel(orderID string) (params url.Values) {
	params = make(url.Values)
	params.Set("customerId", p.ID)
	params.Set("orderId", orderID)
	params.Set("timestamp", strconv.FormatInt(time.Now().Unix()*1000, 10))
	params.Set("traceId", p.TraceID())
	params.Set("version", "1.0")
	return
}

// ToJSON 将param转换为支付平台请求的body:JSON数据
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

// DeviceType 通过platform获得支付平台DeviceType
func (p *Pay) DeviceType(platform string) (t int64) {
	// 支付设备渠道类型，  1 pc 2 webapp 3 app 4jsapi 5 server 6小程序支付 7聚合二维码支付
	switch platform {
	case "ios", "android":
		return 3
	default:
		return 1
	}
}

// ServiceType 通过platform获得支付平台ServiceType
func (p *Pay) ServiceType(platform string) (t int) {
	/*
		业务方业务类型，用于业务方定制支付渠道，不同的serviceType可以配置成不同的支付渠道列表
		每个渠道列表可以自定义顺序,以下值具有特殊含义:

		1.  7 签约代扣类，如微信代扣，支付宝代扣 （5.25 支持）
		2. 100 表示IAP支付 （5.24 支持）
		3. 100 IAP代扣也传100，根据subscribeType区分是不是代扣
		4. 99 表示 客户端B币快捷支付 （5.26 支持）
		5 98 表示 业务方B币快捷支付（1.1 B币支付）
	*/
	switch platform {
	case "ios", "android":
		return 0
	default:
		return 99
	}
}

// Sign 对param进行签名
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

// Verify 对param进行签名验证
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
