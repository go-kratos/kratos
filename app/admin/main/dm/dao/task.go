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

	"go-common/library/ecode"
	"go-common/library/log"
)

const (
	_selectQuery = "SELECT index.dmid,index.oid,index.mid,index.state,content.msg,b_long2ip(content.ip),content.ctime FROM ods.ods_dm_index AS index, ods.ods_dm_content AS content WHERE index.dmid=content.dmid AND index.state in (0,2,6) %s limit 1000000"
)

var (
	signParams = []string{"appKey", "timestamp", "version"}
)

// SendTask send task to BI
func (d *Dao) SendTask(c context.Context, taskSQL []string) (statusURL string, err error) {
	var (
		sql string
		res struct {
			Code      int64  `json:"code"`
			StatusURL string `json:"jobStatusUrl"`
			Message   string `json:"msg"`
		}
		params = url.Values{}
	)
	if len(taskSQL) > 0 {
		sql = fmt.Sprintf(" AND %s", strings.Join(taskSQL, " AND "))
	} else {
		err = ecode.RequestErr
		return
	}
	log.Warn("send task sql(%s)", fmt.Sprintf(_selectQuery, sql))
	params.Set("appKey", "672bc22888af701529e8b3052fd2c4a7")
	params.Set("query", fmt.Sprintf(_selectQuery, sql))
	params.Set("timestamp", time.Now().Format("2006-01-02 15:04:05"))
	params.Set("version", "1.0")
	params.Set("signMethod", "md5")
	uri := d.berserkerURI + "?" + sign(params)
	log.Warn("send task uri(%s)", uri)
	req, err := http.NewRequest(http.MethodGet, uri, nil)
	if err != nil {
		log.Error("http.NewRequest(%s) error(%v)", uri, err)
		return
	}
	for i := 0; i < 3; i++ {
		if err = d.httpCli.Do(c, req, &res); err != nil {
			log.Error("d.httpCli.Do error:%v", err)
			time.Sleep(100 * time.Millisecond)
			continue
		}
		if res.Code != 200 {
			err = fmt.Errorf("uri:%s,code:%d", uri, res.Code)
			log.Error("d.res.Code error:%v", err)
			time.Sleep(100 * time.Millisecond)
			continue
		}
		break
	}
	if err != nil {
		log.Error("d.SendTask(uri:%s) error(%v)", uri, err)
	}
	return res.StatusURL, err
}

// Sign calculate appkey and appsecret sign.
func sign(params url.Values) (query string) {
	tmp := params.Encode()
	signTmp := encode(params)
	if strings.IndexByte(tmp, '+') > -1 {
		tmp = strings.Replace(tmp, "+", "%20", -1)
	}
	var b bytes.Buffer
	b.WriteString("bee5e4b744a22a59abbaecc7ade5de9c")
	b.WriteString(signTmp)
	b.WriteString("bee5e4b744a22a59abbaecc7ade5de9c")
	mh := md5.Sum(b.Bytes())
	// fmt.Println(b.String())
	// query
	var qb bytes.Buffer
	qb.WriteString(tmp)
	qb.WriteString("&sign=")
	qb.WriteString(strings.ToUpper(hex.EncodeToString(mh[:])))
	query = qb.String()
	return
}

// Encode encodes the values into ``sign encoded'' form
// ("barbazfooquux") sorted by key.
func encode(v url.Values) string {
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
