package dataplatform

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/hex"
	xhttp "net/http"
	"net/url"
	"sort"
	"strings"
	"time"

	"go-common/app/job/main/growup/model"
	income "go-common/app/job/main/growup/model/income"
	"go-common/library/conf/env"
	"go-common/library/log"
)

var (
	signParams = []string{"appKey", "timestamp", "version"}
	_userAgent = "User-Agent"
)

func (d *Dao) setParams() url.Values {
	params := url.Values{}
	params.Set("appKey", d.c.DPClient.Key)
	params.Set("timestamp", time.Now().Format("2006-01-02 15:04:05"))
	params.Set("version", "1.0")
	params.Set("signMethod", "md5")
	return params
}

// GetArchiveByMID get archive id by mid
func (d *Dao) GetArchiveByMID(c context.Context, query string) (ids []int64, err error) {
	ids = make([]int64, 0)
	params := d.setParams()
	params.Set("query", query)
	var res struct {
		Code   int                `json:"code"`
		Result []*model.ArchiveID `json:"result"`
	}
	if err = d.NewRequest(c, d.url, "", params, &res); err != nil {
		log.Error("dataplatform.send NewRequest url(%s) error(%v)", d.url+"?"+params.Encode(), err)
		return
	}
	if res.Code != 200 {
		log.Error("dateplatform.send NewRequest error code:%d ; url(%s) ", res.Code, d.url+"?"+params.Encode())
		return
	}
	if res.Code == 200 && len(res.Result) > 0 {
		for _, archive := range res.Result {
			ids = append(ids, archive.ID)
		}
	}
	return
}

// Send ...
func (d *Dao) Send(c context.Context, query string) (infos []*model.ArchiveInfo, err error) {
	log.Info("dateplatform Send start")
	params := url.Values{}
	params.Set("appKey", d.c.DPClient.Key)
	params.Set("timestamp", time.Now().Format("2006-01-02 15:04:05"))
	params.Set("version", "1.0")
	params.Set("signMethod", "md5")
	params.Set("query", query)
	var res struct {
		Code   int                  `json:"code"`
		Result []*model.ArchiveInfo `json:"result"`
	}
	if err = d.NewRequest(c, d.url, "", params, &res); err != nil {
		log.Error("dataplatform.send NewRequest url(%s) error(%v)", d.url+"?"+params.Encode(), err)
		return
	}
	if res.Code != 200 {
		log.Error("dateplatform.send NewRequest error code:%d ; url(%s) ", res.Code, d.url+"?"+params.Encode())
		return
	}
	if res.Code == 200 && len(res.Result) > 0 {
		infos = res.Result
	}
	return
}

// SendSpyRequest send.
func (d *Dao) SendSpyRequest(c context.Context, query string) (infos []*model.Spy, err error) {
	log.Info("dateplatform Send start")
	params := url.Values{}
	params.Set("appKey", d.c.DPClient.Key)
	params.Set("timestamp", time.Now().Format("2006-01-02 15:04:05"))
	params.Set("version", "1.0")
	params.Set("signMethod", "md5")
	params.Set("query", query)
	var res struct {
		Code   int          `json:"code"`
		Result []*model.Spy `json:"result"`
	}
	if err = d.NewRequest(c, d.spyURL, "", params, &res); err != nil {
		log.Error("dataplatform.SendSpyRequest NewRequest url(%s) error(%v)", d.spyURL+"?"+params.Encode(), err)
		return
	}
	if res.Code != 200 {
		log.Error("dateplatform.SendSpyRequest NewRequest error code:%d ; url(%s) ", res.Code, d.spyURL+"?"+params.Encode())
		return
	}
	if res.Code == 200 && len(res.Result) > 0 {
		infos = res.Result
	}
	return
}

// SendBGMRequest get bgm infos
func (d *Dao) SendBGMRequest(c context.Context, query string) (infos []*income.BGM, err error) {
	params := url.Values{}
	params.Set("appKey", d.c.DPClient.Key)
	params.Set("timestamp", time.Now().Format("2006-01-02 15:04:05"))
	params.Set("version", "1.0")
	params.Set("signMethod", "md5")
	params.Set("query", query)
	var res struct {
		Code   int           `json:"code"`
		Result []*income.BGM `json:"result"`
	}
	if err = d.NewRequest(c, d.bgmURL, "", params, &res); err != nil {
		log.Error("dataplatform.SendBGMRequest NewRequest url(%s) error(%v)", d.bgmURL+"?"+params.Encode(), err)
		return
	}
	if res.Code != 200 {
		log.Error("dateplatform.SendBGMRequest NewRequest error code:%d ; url(%s) ", res.Code, d.bgmURL+"?"+params.Encode())
		return
	}
	if res.Code == 200 && len(res.Result) > 0 {
		infos = res.Result
	}
	return
}

// SendBasicDataRequest get basic data status
func (d *Dao) SendBasicDataRequest(c context.Context, query string) (ok bool, err error) {
	params := url.Values{}
	params.Set("appKey", d.c.DPClient.Key)
	params.Set("timestamp", time.Now().Format("2006-01-02 15:04:05"))
	params.Set("version", "1.0")
	params.Set("signMethod", "md5")
	params.Set("query", query)
	var res struct {
		Code   int `json:"code"`
		Result []*struct {
			Stat int `json:"stat"`
		} `json:"result"`
	}
	url := d.basicURL
	if err = d.NewRequest(c, url, "", params, &res); err != nil {
		log.Error("SendBasicDataRequest NewRequest url(%s) error(%v)", url+"?"+params.Encode(), err)
		return
	}
	if res.Code != 200 {
		log.Error("SendBasicDataRequest NewRequest error code:%d ; url(%s) ", res.Code, url+"?"+params.Encode())
		return
	}
	if res.Code == 200 && len(res.Result) > 0 {
		info := res.Result[0]
		if info.Stat == 1 {
			ok = true
		}
	}
	return
}

// NewRequest new http request with method, url, ip, values and headers.
func (d *Dao) NewRequest(c context.Context, url, realIP string, params url.Values, res interface{}) (err error) {
	enc, err := d.sign(params)
	if err != nil {
		log.Error("url:%s,params:%v", url, params)
		return
	}
	if enc != "" {
		url = url + "?" + enc
	}
	req, err := xhttp.NewRequest(xhttp.MethodGet, url, nil)
	if err != nil {
		log.Error("method:%s,url:%s", xhttp.MethodGet, url)
		return
	}
	req.Header.Set(_userAgent, "haoguanwei@bilibili.com  "+env.AppID)
	if err != nil {
		return
	}
	return d.client.Do(c, req, res)
}

// sign calc appkey and appsecret sign.
func (d *Dao) sign(params url.Values) (query string, err error) {
	tmp := params.Encode()
	signTmp := d.encode(params)
	if strings.IndexByte(tmp, '+') > -1 {
		tmp = strings.Replace(tmp, "+", "%20", -1)
	}
	var b bytes.Buffer
	b.WriteString(d.c.DPClient.Secret)
	b.WriteString(signTmp)
	b.WriteString(d.c.DPClient.Secret)
	mh := md5.Sum(b.Bytes())
	// query
	var qb bytes.Buffer
	qb.WriteString(tmp)
	qb.WriteString("&sign=")
	qb.WriteString(strings.ToUpper(hex.EncodeToString(mh[:])))
	query = qb.String()
	return
}

// Encode encodes the values into ``URL encoded'' form
// ("bar=baz&foo=quux") sorted by key.
func (d *Dao) encode(v url.Values) string {
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
		found := false
		for _, p := range signParams {
			if p == k {
				found = true
				break
			}
		}
		if !found {
			continue
		}
		vs := v[k]
		prefix := k
		for _, v := range vs {
			buf.WriteString(prefix)
			buf.WriteString(v)
		}
	}
	return buf.String()
}
