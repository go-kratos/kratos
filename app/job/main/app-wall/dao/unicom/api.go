package unicom

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"

	"go-common/app/job/main/app-wall/model"
	"go-common/app/job/main/app-wall/model/unicom"
	"go-common/library/ecode"
	"go-common/library/log"
	httpx "go-common/library/net/http/blademaster"
)

const (
	// unicom
	_unicomUser                = "000012"
	_unicomPass                = "1pibH5e1BN4V"
	_unicomAppKey              = "com.aop.app.bilibili"
	_unicomMethodQryFlowChange = "com.ssp.method.outqryflowchange"
	_unicomSecurity            = "DVniSMVU6Z3cCIG3vbOn4Fqbof+QJ/6etD+lpa4M4clgj/Dv6XT0syTR8Xgu5nVzKuzro8eiTUzHy/QAzGjp+A=="
	// url
	_unicomFlowExchangeURL = "/openservlet"
	_unicomIPURL           = "/web/statistics/subsystem_2/query_ip.php"
)

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
	if err = d.unicomHTTPPost(c, d.unicomIPURL, params, &res); err != nil {
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

// unicomHTTPGet http get
func (d *Dao) unicomHTTPGet(c context.Context, urlStr string, params url.Values, res interface{}) (err error) {
	return d.wallHTTP(c, d.uclient, http.MethodGet, urlStr, params, res)
}

// unicomHTTPPost http post
func (d *Dao) unicomHTTPPost(c context.Context, urlStr string, params url.Values, res interface{}) (err error) {
	return d.wallHTTP(c, d.uclient, http.MethodPost, urlStr, params, res)
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
	return client.Do(c, req, &res)
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
