package dao

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"go-common/app/interface/openplatform/article/model"
	"go-common/library/log"
	"go-common/library/xstr"
)

var _queryStr = `{"select":[{"name":"tid"},{"name":"oid"},{"name":"log_date"}],"where":{"tid":{"in":[%s]}},"page":{"limit":10,"skip":0}}`

// BerserkerTagArts .
func (d *Dao) BerserkerTagArts(c context.Context, tags []int64) (aids []int64, err error) {
	var (
		query string
		res   struct {
			Code   int
			Msg    string
			Result []struct {
				Tid     int64  `json:"tid"`
				Oid     string `json:"oid"`
				LogDate string `json:"log_date"`
			}
		}
		tmps = make(map[int64]bool)
		aid  int64
		date time.Time
		now  = time.Now()
	)
	query = fmt.Sprintf(_queryStr, xstr.JoinInts(tags))
	if err = d.berserkerQuery(c, query, &res); err != nil {
		return
	}
	if res.Code != 200 {
		log.Error("s.BerserkerTagArts.query code(%d) msg(%s)", res.Code, res.Msg)
		return
	}
	for _, v := range res.Result {
		if date, err = time.Parse("20060102", v.LogDate); err != nil {
			log.Error("s.BerserkerTagArts.time.Parse(%s) error(%+v)", v.LogDate, err)
			return
		}
		if now.Sub(date) > time.Hour*60 {
			continue
		}
		ids := strings.Split(v.Oid, "，")
		var ts []int64
		for _, id := range ids {
			if aid, err = strconv.ParseInt(id, 10, 64); err != nil {
				log.Error("s.BerserkerTagArts.ParseInt(%s) error(%+v)", id, err)
				return
			}
			if !tmps[aid] {
				aids = append(aids, aid)
				tmps[aid] = true
			}
			ts = append(ts, aid)
		}
		d.AddCacheAidsByTag(c, v.Tid, &model.TagArts{Tid: v.Tid, Aids: ts})
	}
	return
}

func (d *Dao) berserkerQuery(c context.Context, query string, res interface{}) (err error) {
	var (
		params = url.Values{}
		now    = time.Now().Format("2006-01-02 15:04:05")
		sign   string
		req    *http.Request
	)
	sign = d.sign(now)
	params.Set("appKey", d.c.Berserker.AppKey)
	params.Set("signMethod", "md5")
	params.Set("timestamp", now)
	params.Set("version", "1.0")
	params.Set("query", query)
	params.Set("sign", sign)

	req, err = http.NewRequest(http.MethodGet, d.c.Berserker.URL+"?"+params.Encode(), nil)
	if err != nil {
		log.Error("d.berserkerQuery.NewRequest error(%+v)", err)
		return
	}
	return d.httpClient.Do(c, req, res)
}

// Sign calc appkey and appsecret sign.
func (d *Dao) sign(ts string) string {
	str := d.c.Berserker.AppSecret + "appKey" + d.c.Berserker.AppKey + "timestamp" + ts + "version1.0" + d.c.Berserker.AppSecret
	mh := md5.Sum([]byte(str))
	return strings.ToUpper(hex.EncodeToString(mh[:]))
}
