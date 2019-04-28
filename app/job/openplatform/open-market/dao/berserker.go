package dao

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"net/url"
	"sort"
	"strings"
	"time"

	"go-common/app/job/openplatform/open-market/model"
	"go-common/library/conf/env"
	"go-common/library/log"

	xhttp "net/http"

	pkgerr "github.com/pkg/errors"
)

const (
	_userAgent = "User-Agent"
	_queryJSON = `{"select":[{"name":"item_id"},{"name":"pv"},{"name":"uv"},{"name":"days_before"}],"where":{"item_id":{"in":[%d]}},"page":{"limit":30,"skip":0}}`
)

var (
	signParams = []string{"appKey", "timestamp", "version"}
)

// QueryPUVCount get puv info from berserker
func (d *Dao) QueryPUVCount(c context.Context, itemID int32) (pv map[int32]int64, uv map[int32]int64, err error) {
	pv = make(map[int32]int64)
	uv = make(map[int32]int64)
	v := make(url.Values, 8)
	v.Set("appKey", d.c.Berserker.Appkey)
	v.Set("signMethod", "md5")
	v.Set("timestamp", time.Now().Format("2006-01-02 15:04:05"))
	v.Set("version", "1.0")
	query := fmt.Sprintf(_queryJSON, itemID)
	v.Set("query", query)
	var res struct {
		Code   int               `json:"code"`
		Result []model.PUVResult `json:"result"`
	}

	if err = d.doHTTPRequest(c, d.c.Berserker.URL, "", v, &res); err != nil {
		log.Error(d.c.Berserker.URL+"?"+v.Encode(), err)
		return
	}

	if res.Code == 200 && len(res.Result) > 0 {
		for _, v := range res.Result {
			pv[v.DaysBefore] = v.PV
			uv[v.DaysBefore] = v.UV
		}
		return
	}
	return
}

// doHttpRequest make a http request for data platform api
func (d *Dao) doHTTPRequest(c context.Context, uri, ip string, params url.Values, res interface{}) (err error) {
	enc, err := d.sign(params)
	if err != nil {
		err = pkgerr.Wrapf(err, "uri:%s,params:%v", uri, params)
		return
	}
	if enc != "" {
		uri = uri + "?" + enc
	}

	req, err := xhttp.NewRequest(xhttp.MethodGet, uri, nil)
	if err != nil {
		err = pkgerr.Wrapf(err, "method:%s,uri:%s", xhttp.MethodGet, uri)
		return
	}
	req.Header.Set(_userAgent, "changxuanran@bilibili.com  "+env.AppID)
	if err != nil {
		return
	}
	return d.client.Do(c, req, res)
}

// Sign calc appkey and appsecret sign.
func (d *Dao) sign(params url.Values) (query string, err error) {
	tmp := params.Encode()
	signTmp := d.encode(params)
	if strings.IndexByte(tmp, '+') > -1 {
		tmp = strings.Replace(tmp, "+", "%20", -1)
	}
	var b bytes.Buffer
	b.WriteString(d.c.Berserker.Secret)
	b.WriteString(signTmp)
	b.WriteString(d.c.Berserker.Secret)
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

//func (d *Dao) reverse(s string) string {
//	n := len(s)
//	runes := make([]rune, n)
//	for _, rune := range s {
//		n--
//		runes[n] = rune
//	}
//	return string(runes[n:])
//}
