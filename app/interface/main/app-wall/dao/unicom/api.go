package unicom

import (
	"bytes"
	"context"
	"crypto/des"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"

	"go-common/app/interface/main/app-wall/model"
	"go-common/app/interface/main/app-wall/model/unicom"
	"go-common/library/ecode"
	"go-common/library/log"
	httpx "go-common/library/net/http/blademaster"
)

const (
	_cpid          = "bilibl"
	_spid          = "979"
	_apptype       = "2"
	_broadbandPass = "9ed226d9"
	// url
	_orderURL       = "/videoif/order.do"
	_ordercancelURL = "/videoif/cancelOrder.do"
	_sendsmscodeURL = "/videoif/sendSmsCode.do"
	_smsNumberURL   = "/videoif/smsNumber.do"
	// unicom
	_unicomIPURL               = "/web/statistics/subsystem_2/query_ip.php"
	_unicomUser                = "000012"
	_unicomPass                = "1pibH5e1BN4V"
	_unicomFlowExchangeURL     = "/openservlet"
	_unicomAppKey              = "com.aop.app.bilibili"
	_unicomSecurity            = "DVniSMVU6Z3cCIG3vbOn4Fqbof+QJ/6etD+lpa4M4clgj/Dv6XT0syTR8Xgu5nVzKuzro8eiTUzHy/QAzGjp+A=="
	_unicomAppMethodFlow       = "com.ssp.method.outflowchange"
	_unicomMethodNumber        = "com.aop.method.checkusernumber"
	_unicomMethodFlowPre       = "com.ssp.method.outflowpre"
	_unicomMethodQryFlowChange = "com.ssp.method.outqryflowchange"
)

// Order unicom order
func (d *Dao) Order(c context.Context, usermob, channel string, ordertype int) (data *unicom.BroadbandOrder, msg string, err error) {
	params := url.Values{}
	params.Set("cpid", _cpid)
	params.Set("spid", _spid)
	params.Set("ordertype", strconv.Itoa(ordertype))
	params.Set("userid", usermob)
	params.Set("apptype", _apptype)
	if channel != "" {
		params.Set("channel", channel)
	}
	var res struct {
		Code string `json:"resultcode"`
		Msg  string `json:"errorinfo"`
		*unicom.BroadbandOrder
	}
	if err = d.broadbandHTTPGet(c, d.orderURL, params, &res); err != nil {
		log.Error("unicom order url(%v) error(%v)", d.orderURL+"?"+params.Encode(), err)
		return
	}
	b, _ := json.Marshal(&res)
	log.Info("unicom order url(%v) response(%s)", d.orderURL+"?"+params.Encode(), b)
	if res.Code != "0" {
		err = ecode.String(res.Code)
		msg = res.Msg
		log.Error("unicom order url(%v) code(%s) Msg(%s)", d.orderURL+"?"+params.Encode(), res.Code, res.Msg)
		return
	}
	data = res.BroadbandOrder
	return
}

// CancelOrder unicom cancel order
func (d *Dao) CancelOrder(c context.Context, usermob string) (data *unicom.BroadbandOrder, msg string, err error) {
	params := url.Values{}
	params.Set("cpid", _cpid)
	params.Set("spid", _spid)
	params.Set("userid", usermob)
	params.Set("apptype", _apptype)
	var res struct {
		Code string `json:"resultcode"`
		Msg  string `json:"errorinfo"`
		*unicom.BroadbandOrder
	}
	if err = d.broadbandHTTPGet(c, d.ordercancelURL, params, &res); err != nil {
		log.Error("unicom cancel order url(%s) error(%v)", d.ordercancelURL+"?"+params.Encode(), err)
		return
	}
	b, _ := json.Marshal(&res)
	log.Info("unicom cancel order url(%s) response(%s)", d.ordercancelURL+"?"+params.Encode(), b)
	if res.Code != "0" {
		err = ecode.String(res.Code)
		msg = res.Msg
		log.Error("unicom cancel order url(%v) code(%s) Msg(%s)", d.orderURL+"?"+params.Encode(), res.Code, res.Msg)
		return
	}
	data = res.BroadbandOrder
	return
}

// UnicomIP unicom ip orders
func (d *Dao) UnicomIP(c context.Context, now time.Time) (unicomIPs []*unicom.UnicomIP, err error) {
	params := url.Values{}
	params.Set("user", _unicomUser)
	tick := strconv.FormatInt(now.Unix(), 10)
	params.Set("tick", tick)
	mh := md5.Sum([]byte(_unicomUser + tick + _unicomPass))
	var (
		key string
	)
	if key = hex.EncodeToString(mh[:]); len(key) > 16 {
		key = key[:16]
	}
	params.Set("key", key)
	var res struct {
		Code           int   `json:"code"`
		LastUpdateTime int64 `json:"last_update_time"`
		Desc           []struct {
			StartIP string `json:"start_ip"`
			Length  string `json:"length"`
		} `json:"desc"`
	}
	if err = d.broadbandHTTPPost(c, d.unicomIPURL, params, &res); err != nil {
		log.Error("unicom ip order url(%s) error(%v)", d.unicomIPURL+"?"+params.Encode(), err)
		return
	}
	if res.Code != 0 {
		err = ecode.Int(res.Code)
		log.Error("unicom ip order url(%s) res code (%d)", d.unicomIPURL+"?"+params.Encode(), res.Code)
		return
	}
	b, _ := json.Marshal(&res)
	log.Info("unicom ip url(%s) response(%s)", d.unicomIPURL+"?"+params.Encode(), b)
	for _, uip := range res.Desc {
		uiplen, _ := strconv.Atoi(uip.Length)
		if uiplen < 1 {
			log.Error("unicom ip length 0")
			continue
		}
		ipEndInt := model.InetAtoN(uip.StartIP) + uint32((uiplen - 1))
		ipEnd := model.InetNtoA(ipEndInt)
		unicomIP := &unicom.UnicomIP{}
		unicomIP.UnicomIPStrToint(uip.StartIP, ipEnd)
		unicomIPs = append(unicomIPs, unicomIP)
	}
	return
}

// SendSmsCode unicom sms code
func (d *Dao) SendSmsCode(c context.Context, phone string) (msg string, err error) {
	var (
		key       = []byte(_broadbandPass)
		phoneByte = []byte(phone)
		userid    string
	)
	userid, err = d.desEncrypt(phoneByte, key)
	if err != nil {
		log.Error("d.desEncrypt error(%v)", err)
		return
	}
	params := url.Values{}
	params.Set("cpid", _cpid)
	params.Set("userid", string(userid))
	params.Set("apptype", _apptype)
	var res struct {
		Code string `json:"resultcode"`
		Msg  string `json:"errorinfo"`
	}
	if err = d.unicomHTTPGet(c, d.sendsmscodeURL, params, &res); err != nil {
		log.Error("unicom sendsmscode url(%v) error(%v)", d.sendsmscodeURL+"?"+params.Encode(), err)
		return
	}
	b, _ := json.Marshal(&res)
	log.Info("unicom sendsmscode url(%v) response(%s)", d.sendsmscodeURL+"?"+params.Encode(), b)
	if res.Code != "0" {
		err = ecode.String(res.Code)
		msg = res.Msg
		log.Error("unicom sendsmscode url(%v) code(%s) Msg(%s)", d.sendsmscodeURL+"?"+params.Encode(), res.Code, res.Msg)
		return
	}
	return
}

// SmsNumber unicom sms usermob
func (d *Dao) SmsNumber(c context.Context, phone string, code int) (usermob string, msg string, err error) {
	var (
		key       = []byte(_broadbandPass)
		phoneByte = []byte(phone)
		userid    string
	)
	userid, err = d.desEncrypt(phoneByte, key)
	if err != nil {
		log.Error("d.desEncrypt error(%v)", err)
		return
	}
	params := url.Values{}
	params.Set("cpid", _cpid)
	params.Set("userid", userid)
	params.Set("vcode", strconv.Itoa(code))
	params.Set("apptype", _apptype)
	var res struct {
		Code    string `json:"resultcode"`
		Usermob string `json:"userid"`
		Msg     string `json:"errorinfo"`
	}
	if err = d.unicomHTTPGet(c, d.smsNumberURL, params, &res); err != nil {
		log.Error("unicom smsNumberURL url(%v) error(%v)", d.smsNumberURL+"?"+params.Encode(), err)
		return
	}
	b, _ := json.Marshal(&res)
	log.Info("unicom sendsmsnumber url(%v) response(%s)", d.smsNumberURL+"?"+params.Encode(), b)
	if res.Code != "0" {
		err = ecode.String(res.Code)
		msg = res.Msg
		log.Error("unicom sendsmsnumber url(%v) code(%s) Msg(%s)", d.smsNumberURL+"?"+params.Encode(), res.Code, res.Msg)
		return
	}
	usermob = res.Usermob
	return
}

// FlowExchange flow exchange
func (d *Dao) FlowExchange(c context.Context, phone int, flowcode string, requestNo int64, ts time.Time) (orderID, outorderID, msg string, err error) {
	outorderIDStr := "bili" + ts.Format("20060102") + strconv.FormatInt(requestNo%10000000, 10)
	if len(outorderIDStr) > 22 {
		outorderIDStr = outorderIDStr[:22]
	}
	param := url.Values{}
	param.Set("appkey", _unicomAppKey)
	param.Set("apptx", strconv.FormatInt(requestNo, 10))
	param.Set("flowexchangecode", flowcode)
	param.Set("method", _unicomAppMethodFlow)
	param.Set("outorderid", outorderIDStr)
	param.Set("timestamp", ts.Format("2006-01-02 15:04:05"))
	param.Set("usernumber", strconv.Itoa(phone))
	urlVal := d.urlParams(param)
	urlVal = urlVal + "&" + d.sign(urlVal)
	var res struct {
		Code       string `json:"respcode"`
		Msg        string `json:"respdesc"`
		OrderID    string `json:"orderid"`
		OutorderID string `json:"outorderid"`
	}
	if err = d.unicomHTTPGet(c, d.unicomFlowExchangeURL+"?"+urlVal, nil, &res); err != nil {
		log.Error("unicom flow change url(%v) error(%v)", d.unicomFlowExchangeURL+"?"+urlVal, err)
		return
	}
	b, _ := json.Marshal(&res)
	log.Info("unicom flow url(%v) response(%s)", d.unicomFlowExchangeURL+"?"+urlVal, b)
	msg = res.Msg
	if res.Code != "0000" {
		err = ecode.String(res.Code)
		log.Error("unicom flow change url(%v) code(%v) msg(%v)", d.unicomFlowExchangeURL+"?"+urlVal, res.Code, res.Msg)
		return
	}
	orderID = res.OrderID
	outorderID = res.OutorderID
	return
}

// PhoneVerification unicom phone verification
func (d *Dao) PhoneVerification(c context.Context, phone string, requestNo int64, ts time.Time) (msg string, err error) {
	param := url.Values{}
	param.Set("appkey", _unicomAppKey)
	param.Set("apptx", strconv.FormatInt(requestNo, 10))
	param.Set("method", _unicomMethodNumber)
	param.Set("timestamp", ts.Format("2006-01-02 15:04:05"))
	param.Set("usernumber", phone)
	urlVal := d.urlParams(param)
	urlVal = urlVal + "&" + d.sign(urlVal)
	var res struct {
		Code string `json:"respcode"`
		Msg  string `json:"respdesc"`
	}
	if err = d.unicomHTTPGet(c, d.unicomFlowExchangeURL+"?"+urlVal, nil, &res); err != nil {
		log.Error("unicom phone url(%v) error(%v)", d.unicomFlowExchangeURL+"?"+urlVal, err)
		return
	}
	b, _ := json.Marshal(&res)
	log.Info("unicom phone url(%v) response(%s)", d.unicomFlowExchangeURL+"?"+urlVal, b)
	msg = res.Msg
	if res.Code != "0000" {
		err = ecode.String(res.Code)
		log.Error("unicom phone url(%v) code(%v) msg(%v)", d.unicomFlowExchangeURL+"?"+urlVal, res.Code, res.Msg)
		return
	}
	return
}

// FlowPre unicom phone flow pre
func (d *Dao) FlowPre(c context.Context, phone int, requestNo int64, ts time.Time) (msg string, err error) {
	param := url.Values{}
	param.Set("appkey", _unicomAppKey)
	param.Set("apptx", strconv.FormatInt(requestNo, 10))
	param.Set("method", _unicomMethodFlowPre)
	param.Set("timestamp", ts.Format("2006-01-02 15:04:05"))
	param.Set("usernumber", strconv.Itoa(phone))
	urlVal := d.urlParams(param)
	urlVal = urlVal + "&" + d.sign(urlVal)
	var res struct {
		Code   string `json:"respcode"`
		Notice string `json:"noticecontent"`
		Msg    string `json:"respdesc"`
	}
	if err = d.unicomHTTPGet(c, d.unicomFlowExchangeURL+"?"+urlVal, nil, &res); err != nil {
		log.Error("unicom flowpre url(%v) error(%v)", d.unicomFlowExchangeURL+"?"+urlVal, err)
		return
	}
	b, _ := json.Marshal(&res)
	log.Info("unicom flowpre url(%v) response(%s)", d.unicomFlowExchangeURL+"?"+urlVal, b)
	msg = res.Msg
	if res.Code != "0000" {
		if res.Code == "0001" {
			err = ecode.String(res.Code)
			msg = res.Notice
		} else {
			err = ecode.String(res.Code)
			log.Error("unicom flowpre url(%v) code(%v) msg(%v)", d.unicomFlowExchangeURL+"?"+urlVal, res.Code, res.Msg)
		}
		return
	}
	return
}

// FlowQry unicom phone qryflowchange
func (d *Dao) FlowQry(c context.Context, phone int, requestNo int64, outorderid, orderid string, ts time.Time) (orderstatus, msg string, err error) {
	param := url.Values{}
	param.Set("appkey", _unicomAppKey)
	param.Set("apptx", strconv.FormatInt(requestNo, 10))
	param.Set("method", _unicomMethodQryFlowChange)
	param.Set("orderid", orderid)
	param.Set("outorderid", outorderid)
	param.Set("timestamp", ts.Format("2006-01-02 15:04:05"))
	param.Set("usernumber", strconv.Itoa(phone))
	urlVal := d.urlParams(param)
	urlVal = urlVal + "&" + d.sign(urlVal)
	var res struct {
		Code        string `json:"respcode"`
		Orderstatus string `json:"orderstatus"`
		Failurtype  string `json:"failurtype"`
		Msg         string `json:"respdesc"`
	}
	if err = d.unicomHTTPGet(c, d.unicomFlowExchangeURL+"?"+urlVal, nil, &res); err != nil {
		log.Error("unicom flowQry url(%v) error(%v)", d.unicomFlowExchangeURL+"?"+urlVal, err)
		return
	}
	b, _ := json.Marshal(&res)
	log.Info("unicom flowQry url(%v) response(%s)", d.unicomFlowExchangeURL+"?"+urlVal, b)
	msg = res.Msg
	if res.Code != "0000" {
		err = ecode.String(res.Code)
		log.Error("unicom flowQry url(%v) code(%v) msg(%v)", d.unicomFlowExchangeURL+"?"+urlVal, res.Code, res.Msg)
		return
	}
	orderstatus = res.Orderstatus
	return
}

// broadbandHTTPGet http get
func (d *Dao) broadbandHTTPGet(c context.Context, urlStr string, params url.Values, res interface{}) (err error) {
	return d.wallHTTP(c, d.client, http.MethodGet, urlStr, params, res)
}

// broadbandHTTPPost http post
func (d *Dao) broadbandHTTPPost(c context.Context, urlStr string, params url.Values, res interface{}) (err error) {
	return d.wallHTTP(c, d.client, http.MethodPost, urlStr, params, res)
}

// unicomHTTPGet http get
func (d *Dao) unicomHTTPGet(c context.Context, urlStr string, params url.Values, res interface{}) (err error) {
	return d.wallHTTP(c, d.uclient, http.MethodGet, urlStr, params, res)
}

// wallHTTP http
func (d *Dao) wallHTTP(c context.Context, client *httpx.Client, method, urlStr string, params url.Values, res interface{}) (err error) {
	var (
		req *http.Request
	)
	ru := urlStr
	if params != nil {
		ru = urlStr + "?" + params.Encode()
	}
	switch method {
	case http.MethodGet:
		req, err = http.NewRequest(http.MethodGet, ru, nil)
	default:
		req, err = http.NewRequest(http.MethodPost, urlStr, strings.NewReader(params.Encode()))
	}
	if err != nil {
		log.Error("unicom_http.NewRequest url(%s) error(%v)", urlStr+"?"+params.Encode(), err)
		return
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("X-BACKEND-BILI-REAL-IP", "")
	return d.client.Do(c, req, &res)
}

func (d *Dao) desEncrypt(src, key []byte) (string, error) {
	block, err := des.NewCipher(key)
	if err != nil {
		return "", err
	}
	bs := block.BlockSize()
	src = d.pkcs5Padding(src, bs)
	if len(src)%bs != 0 {
		return "", errors.New("Need a multiple of the blocksize")
	}
	out := make([]byte, len(src))
	dst := out
	for len(src) > 0 {
		block.Encrypt(dst, src[:bs])
		src = src[bs:]
		dst = dst[bs:]
	}
	encodeString := base64.StdEncoding.EncodeToString(out)
	return encodeString, nil
}

func (d *Dao) pkcs5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func (d *Dao) urlParams(v url.Values) string {
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
		for _, v := range vs {
			if buf.Len() > 0 {
				buf.WriteByte('&')
			}
			buf.WriteString(prefix)
			buf.WriteString(v)
		}
	}
	return buf.String()
}

func (d *Dao) sign(params string) string {
	str := strings.Replace(params, "&", "$", -1)
	str2 := strings.Replace(str, "=", "$", -1)
	mh := md5.Sum([]byte(_unicomSecurity + "$" + str2 + "$" + _unicomSecurity))
	signparam := url.Values{}
	signparam.Set("sign", base64.StdEncoding.EncodeToString(mh[:]))
	return signparam.Encode()
}
