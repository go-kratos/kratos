package dao

import (
	"context"
	"encoding/json"
	"net/url"
	"strconv"
	"time"

	"go-common/app/service/main/member/model"
	"go-common/library/cache/memcache"
	"go-common/library/ecode"
	"go-common/library/log"
)

const (
	_expLogID = 11
)

// Exp get user exp from cache,if miss get from db.
func (d *Dao) Exp(c context.Context, mid int64) (exp int64, err error) {
	if exp, err = d.expCache(c, mid); err == nil {
		return
	}
	if err != nil && err != memcache.ErrNotFound {
		return
	}
	if exp, err = d.ExpDB(c, mid); err != nil {
		return
	}
	d.SetExpCache(c, mid, exp)
	return
}

// Exps get exps by mids.
func (d *Dao) Exps(c context.Context, mids []int64) (exps map[int64]int64, err error) {
	exps, miss, err := d.expsCache(c, mids)
	if err != nil {
		return
	}
	if len(miss) == 0 {
		return
	}
	for _, mid := range miss {
		mid := mid
		exp, err := d.ExpDB(c, mid)
		if err != nil {
			log.Error("exp mid %d err %v", mid, err)
			err = nil
			continue
		}
		exps[mid] = exp
		d.cache.Do(c, func(ctx context.Context) {
			d.SetExpCache(ctx, mid, exp)
		})
	}
	return
}

// ExpLog is
func (d *Dao) ExpLog(ctx context.Context, mid int64, ip string) ([]*model.UserLog, error) {
	t := time.Now().Add(-time.Hour * 24 * 7)
	params := url.Values{}
	params.Set("mid", strconv.FormatInt(mid, 10))
	params.Set("appid", "log_user_action")
	params.Set("business", strconv.FormatInt(_expLogID, 10))
	params.Set("pn", "1")
	params.Set("ps", "1000")
	params.Set("ctime_from", t.Format("2006-01-02 00:00:00"))
	params.Set("sort", "desc")
	params.Set("order", "ctime")

	res := &model.SearchResult{}
	if err := d.client.Get(ctx, d.c.Host.Search+_searchLogURI, ip, params, res); err != nil {
		return nil, err
	}
	if res.Code != 0 {
		return nil, ecode.Int(res.Code)
	}
	logs := asExpLog(res)
	return logs, nil
}

func asExpLog(res *model.SearchResult) []*model.UserLog {
	logs := make([]*model.UserLog, 0, len(res.Data.Result))
	for _, r := range res.Data.Result {
		ts, err := time.ParseInLocation("2006-01-02 15:04:05", r.Ctime, time.Local)
		if err != nil {
			log.Warn("Failed to parse log ctime: ctime: %s: %+v", r.Ctime, err)
			continue
		}
		content := map[string]string{
			"from_exp": "",
			"operater": "",
			"reason":   "",
			"to_exp":   "",
			"log_id":   "",
		}
		if err := json.Unmarshal([]byte(r.ExtraData), &content); err != nil {
			log.Warn("Failed to parse extra data in exp log: mid: %d, extra_data: %s: %+v", r.Mid, r.ExtraData, err)
			continue
		}
		l := &model.UserLog{
			Mid:     r.Mid,
			IP:      r.IP,
			TS:      ts.Unix(),
			LogID:   content["log_id"],
			Content: content,
		}
		logs = append(logs, l)
	}
	return logs
}
