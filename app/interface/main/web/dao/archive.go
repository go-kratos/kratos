package dao

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"

	"go-common/library/cache/redis"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/net/metadata"
	"go-common/library/xstr"

	"github.com/pkg/errors"
)

const (
	_keyArcAppeal = "arc_appeal_%d_%d"
	_testerGroup  = "20"
)

// ArcReport add archive report
func (d *Dao) ArcReport(c context.Context, mid, aid, tp int64, reason, pics string) (err error) {
	var (
		params = url.Values{}
		ip     = metadata.String(c, metadata.RemoteIP)
	)
	params.Set("mid", strconv.FormatInt(mid, 10))
	params.Set("aid", strconv.FormatInt(aid, 10))
	params.Set("type", strconv.FormatInt(tp, 10))
	params.Set("reason", reason)
	params.Set("pics", pics)
	var res struct {
		Code int `json:"code"`
	}
	if err = d.httpW.Post(c, d.arcReportURL, ip, params, &res); err != nil {
		return
	}
	if res.Code != ecode.OK.Code() {
		log.Error("archive report(%s) param(%v) ecode err(%d)", d.arcReportURL, params, res.Code)
		err = ecode.Int(res.Code)
	}
	return
}

// ArcAppeal add archive appeal.
func (d *Dao) ArcAppeal(c context.Context, mid int64, data map[string]string, business int) (err error) {
	var (
		params = url.Values{}
		ip     = metadata.String(c, metadata.RemoteIP)
	)
	for name, value := range data {
		params.Set(name, value)
	}
	params.Set("mid", strconv.FormatInt(mid, 10))
	params.Set("business", strconv.Itoa(business))
	if v, ok := data["attach"]; ok && v != "" {
		params.Set("attachments", v)
	}
	var res struct {
		Code int `json:"code"`
	}
	if err = d.httpW.Post(c, d.arcAppealURL, ip, params, &res); err != nil {
		return
	}
	if res.Code != ecode.OK.Code() {
		log.Error("archive report(%s) ecode err(%d)", d.arcAppealURL+"?"+params.Encode(), res.Code)
		err = ecode.Int(res.Code)
	}
	return
}

// AppealTags get appeal tags.
func (d *Dao) AppealTags(c context.Context, business int) (rs json.RawMessage, err error) {
	var (
		params = url.Values{}
		ip     = metadata.String(c, metadata.RemoteIP)
	)
	params.Set("business", strconv.Itoa(business))
	var res struct {
		Code int             `json:"code"`
		Data json.RawMessage `json:"data"`
	}
	if err = d.httpR.Get(c, d.appealTagsURL, ip, params, &res); err != nil {
		return
	}
	if res.Code != ecode.OK.Code() {
		log.Error("archive report(%s) param(%v) ecode err(%d)", d.arcReportURL, params, res.Code)
		err = ecode.Int(res.Code)
	}
	rs = res.Data
	return
}

// RelatedAids get related aids from bigdata
func (d *Dao) RelatedAids(c context.Context, aid int64) (aids []int64, err error) {
	var (
		params = url.Values{}
		ip     = metadata.String(c, metadata.RemoteIP)
	)
	params.Set("key", strconv.FormatInt(aid, 10))
	var res struct {
		Code int `json:"code"`
		Data []*struct {
			Value string `json:"value"`
		} `json:"data"`
	}
	if err = d.httpR.Get(c, d.relatedURL, ip, params, &res); err != nil {
		log.Error("realte url(%s) error(%v) ", d.relatedURL+"?"+params.Encode(), err)
		return
	}
	if res.Code != 0 {
		log.Error("url(%s) res code(%d) or res.result(%v)", d.relatedURL+"?"+params.Encode(), res.Code, res.Data)
		err = ecode.Int(res.Code)
		return
	}
	if res.Data == nil {
		err = nil
		return
	}
	if len(res.Data) > 0 {
		if aids, err = xstr.SplitInts(res.Data[0].Value); err != nil {
			log.Error("realte aids url(%s) error(%v)", d.relatedURL+"?"+params.Encode(), err)
		}
	}
	return
}

func keyArcAppealLimit(mid, aid int64) string {
	return fmt.Sprintf(_keyArcAppeal, mid, aid)
}

// SetArcAppealCache set arc appeal cache.
func (d *Dao) SetArcAppealCache(c context.Context, mid, aid int64) (err error) {
	key := keyArcAppealLimit(mid, aid)
	conn := d.redisBak.Get(c)
	defer conn.Close()
	if err = conn.Send("SET", key, "1"); err != nil {
		log.Error("SetArcAppealCache conn.Send(SET, %s) error(%v)", key, err)
		return
	}
	if err = conn.Send("EXPIRE", key, d.redisAppealLimitExpire); err != nil {
		log.Error("SetArcAppealCache conn.Send(Expire, %s, %d) error(%v)", key, d.redisAppealLimitExpire, err)
		return
	}
	if err = conn.Flush(); err != nil {
		log.Error("SetArcAppealCache conn.Flush error(%v)", err)
		return
	}
	for i := 0; i < 2; i++ {
		if _, err = conn.Receive(); err != nil {
			log.Error("SetArcAppealCache conn.Receive() error(%v)", err)
			return
		}
	}
	return
}

// ArcAppealCache get arc appeal cache.
func (d *Dao) ArcAppealCache(c context.Context, mid, aid int64) (err error) {
	key := keyArcAppealLimit(mid, aid)
	conn := d.redisBak.Get(c)
	defer conn.Close()
	if _, err = redis.Bytes(conn.Do("GET", key)); err != nil {
		if err == redis.ErrNil {
			err = nil
			return
		}
		log.Error("ArcAppealCache conn.Do(GET, %s) error(%v)", key, err)
	}
	err = ecode.ArcAppealLimit
	return
}

// Special manager special mid.
func (d *Dao) Special(c context.Context) (midsM map[int64]struct{}, err error) {
	params := url.Values{}
	params.Set("group_id", _testerGroup)
	var res struct {
		Code int `json:"code"`
		Data []struct {
			Mid int64 `json:"mid"`
		} `json:"data"`
	}
	if err = d.httpR.Get(c, d.special, "", params, &res); err != nil {
		err = errors.Wrap(err, d.special+"?"+params.Encode())
		return
	}
	if res.Code != ecode.OK.Code() {
		err = errors.Wrap(ecode.Int(res.Code), d.special+"?"+params.Encode())
		return
	}
	midsM = make(map[int64]struct{}, len(res.Data))
	for _, l := range res.Data {
		midsM[l.Mid] = struct{}{}
	}
	return
}
