package oplog

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"math"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"

	"go-common/app/admin/main/dm/model/oplog"
	"go-common/library/conf/env"
	"go-common/library/log"

	xhttp "net/http"

	pkgerr "github.com/pkg/errors"
)

const (
	_singleQueryDMLogHbase = `{"startRow": "%s","stopRow": "%s","columns": {"family":"%s"}}`
	_userAgent             = "User-Agent"
)

var (
	signParams = []string{"appKey", "timestamp", "version"}
)

// QueryOpLogs 查找弹幕操作日志，前方高能，这是一段极其恶心的代码（1. 数据平台的key和secret是根据个人用户生成目前是我的账号（madou） 2.sign算法是根据appkey，timestamp，version生成并大小写敏感）
func (d *Dao) QueryOpLogs(c context.Context, dmid int64) (infos []*oplog.InfocResult, err error) {
	v := make(url.Values, 8)
	v.Set("appKey", d.key)
	v.Set("signMethod", "md5")
	v.Set("timestamp", time.Now().Format("2006-01-02 15:04:05"))
	v.Set("version", "1.0")
	//默认只查询三个月
	startRow, stopRow := d.makeRowKeyScope(dmid, -3)
	query := fmt.Sprintf(_singleQueryDMLogHbase, startRow, stopRow, "i")
	v.Set("query", query)
	var res struct {
		Code   int                  `json:"code"`
		Result []*oplog.InfocResult `json:"result"`
	}
	if err = d.doHTTPRequest(c, d.infocQueryURL, "", v, &res); err != nil {
		log.Error("berserker url(%v), err(%v)", d.infocQueryURL+"?"+v.Encode(), err)
		return
	}
	if res.Code == 200 && len(res.Result) > 0 {
		infos = res.Result
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
	req.Header.Set(_userAgent, "haoguanwei@bilibili.com  "+env.AppID)
	if err != nil {
		return
	}
	return d.httpCli.Do(c, req, res)
}

// Sign calc appkey and appsecret sign.
func (d *Dao) sign(params url.Values) (query string, err error) {
	tmp := params.Encode()
	signTmp := d.encode(params)
	if strings.IndexByte(tmp, '+') > -1 {
		tmp = strings.Replace(tmp, "+", "%20", -1)
	}
	var b bytes.Buffer
	b.WriteString(d.secret)
	b.WriteString(signTmp)
	b.WriteString(d.secret)
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

// rowkey存储方式： [dmid倒序补零20位][(Long.Max_Value - timestamp) 的结果后10位]
func (d *Dao) makeRowKeyScope(dmid int64, months int) (startRow, endRow string) {
	endTime := time.Now()
	startTime := endTime.AddDate(0, months, 0)
	startTmp := strconv.FormatInt(math.MaxInt64-startTime.Unix(), 10)
	endTmp := strconv.FormatInt(math.MaxInt64-endTime.Unix(), 10)
	endRow = d.reverse(fmt.Sprintf("%020d", dmid)) + startTmp[len(startTmp)-10:]
	startRow = d.reverse(fmt.Sprintf("%020d", dmid)) + endTmp[len(endTmp)-10:]
	return
}

func (d *Dao) reverse(s string) string {
	n := len(s)
	runes := make([]rune, n)
	for _, rune := range s {
		n--
		runes[n] = rune
	}
	return string(runes[n:])
}
