package dao

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"

	"go-common/app/job/main/push/model"
	"go-common/library/ecode"
	"go-common/library/log"
)

const (
	// 提交查询的接口
	_dpSubmitQueryURL = "http://berserker.bilibili.co/avenger/api/74/query"
)

var (
	dpSignParams = []string{"appKey", "timestamp", "version"}
)

// DpSubmitQuery .
func (d *Dao) DpSubmitQuery(ctx context.Context, query string) (statusRUL string, err error) {
	params := d.params()
	params.Set("query", query)
	var res struct {
		Code      int    `json:"code"`
		Msg       string `json:"msg"`
		StatusURL string `json:"jobStatusUrl"`
	}
	if err = d.newRequest(ctx, _dpSubmitQueryURL, params, &res); err != nil {
		log.Error("d.DpSubmitQuery newRequest url(%s) error(%v)", _dpSubmitQueryURL+"?"+params.Encode(), err)
		return
	}
	if res.Code != http.StatusOK {
		log.Error("d.DpSubmitQuery newRequest error code:%d ; url(%s) ", res.Code, _dpSubmitQueryURL+"?"+params.Encode())
		err = ecode.Int(res.Code)
		err = fmt.Errorf("code(%d) msg(%s)", res.Code, res.Msg)
		return
	}
	statusRUL = res.StatusURL
	return
}

// DpCheckJob .
func (d *Dao) DpCheckJob(ctx context.Context, url string) (res *model.DpCheckJobResult, err error) {
	params := d.params()
	if err = d.newRequest(ctx, url, params, &res); err != nil {
		log.Error("d.DpCheckJob newRequest error(%v)", err)
		return
	}
	if res.Code != http.StatusOK {
		log.Error("d.DpCheckJob newRequest error code:%d ; url(%s) ", res.Code, url+"?"+params.Encode())
		err = fmt.Errorf("code(%d) msg(%s)", res.Code, res.Msg)
	}
	return
}

// DpDownloadFile .
func (d *Dao) DpDownloadFile(ctx context.Context, url string) (content []byte, err error) {
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	if content, err = d.dpClient.Raw(ctx, req); err != nil {
		log.Error("d.dpClient.Raw(%s) error(%v)", url, err)
	}
	return
}

func (d *Dao) params() url.Values {
	params := url.Values{}
	params.Set("appKey", d.c.DpClient.Key)
	params.Set("timestamp", time.Now().Format("2006-01-02 15:04:05"))
	params.Set("version", "1.0")
	params.Set("signMethod", "md5")
	return params
}

// newRequest new http request with method, url, ip, values and headers.
func (d *Dao) newRequest(c context.Context, url string, params url.Values, res interface{}) (err error) {
	enc, err := d.dpSign(params)
	if err != nil {
		log.Error("url:%s,params:%v", url, params)
		return
	}
	if enc != "" {
		url = url + "?" + enc
	}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Error("method:%s,url:%s", http.MethodGet, url)
		return
	}
	return d.httpClient.Do(c, req, res)
}

// dpSign calc appkey and appsecret sign.
func (d *Dao) dpSign(params url.Values) (query string, err error) {
	tmp := params.Encode()
	signTmp := d.encode(params)
	if strings.IndexByte(tmp, '+') > -1 {
		tmp = strings.Replace(tmp, "+", "%20", -1)
	}
	var b bytes.Buffer
	b.WriteString(d.c.DpClient.Secret)
	b.WriteString(signTmp)
	b.WriteString(d.c.DpClient.Secret)
	mh := md5.Sum(b.Bytes())
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
		for _, p := range dpSignParams {
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
