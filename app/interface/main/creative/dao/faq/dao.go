package faq

import (
	"context"
	"encoding/json"
	"fmt"
	"go-common/app/interface/main/creative/conf"
	"go-common/app/interface/main/creative/model/faq"
	"go-common/library/cache/redis"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/metadata"
	xtime "go-common/library/time"
	"net/url"
	"strconv"
)

const (
	_notRobot  = -1
	_rsOk      = "000000"
	_searchURL = "/kb/searchInerDocListBilibili/4"
	_hdKey     = "faq_%s_%d_%d_%d"
)

// Dao  define
type Dao struct {
	c           *conf.Config
	client      *bm.Client
	searchURL   string
	redis       *redis.Pool
	redisExpire int32
}

// New init dao
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		c:           c,
		client:      bm.NewClient(c.HTTPClient.Normal),
		searchURL:   c.Host.HelpAPI + _searchURL,
		redis:       redis.NewPool(c.Redis.Cover.Config),
		redisExpire: int32(120),
	}
	return
}

func keyHd(qTypeID string, keyFlag, pn, ps int) string {
	return fmt.Sprintf(_hdKey, qTypeID, keyFlag, pn, ps)
}

// SetDetailCache  set help detail  to cache.
func (d *Dao) SetDetailCache(c context.Context, qTypeID string, keyFlag, pn, ps, total int, data []*faq.Detail) (err error) {
	conn := d.redis.Get(c)
	defer conn.Close()
	count := 0
	key := keyHd(qTypeID, keyFlag, pn, ps)
	if err = conn.Send("DEL", key); err != nil {
		log.Error("conn.Send(DEL, %s) error(%v)", key, err)
		return
	}
	count++
	var bs []byte
	for _, detail := range data {
		if bs, err = json.Marshal(detail); err != nil {
			log.Error("json.Marshal(%v) error (%v)", detail, err)
			return
		}
		if err = conn.Send("ZADD", key, combineHd(detail.UpdateTime, total), bs); err != nil {
			log.Error("conn.Send(ZADD, %s, %s) error(%v)", key, string(bs), err)
			return
		}
		count++
	}
	if err = conn.Send("EXPIRE", key, d.redisExpire); err != nil {
		log.Error("conn.Send(Expire, %s, %d) error(%v)", key, d.redisExpire, err)
		return
	}
	count++
	if err = conn.Flush(); err != nil {
		log.Error("conn.Flush error(%v)", err)
		return
	}
	for i := 0; i < count; i++ {
		if _, err = conn.Receive(); err != nil {
			log.Error("conn.Receive() error(%v)", err)
			return
		}
	}
	return
}

// DetailCache  get help detail  to cache.
func (d *Dao) DetailCache(c context.Context, qTypeID string, keyFlag, pn, ps int) (res []*faq.Detail, count int, err error) {
	conn := d.redis.Get(c)
	defer conn.Close()
	key := keyHd(qTypeID, keyFlag, pn, ps)
	values, err := redis.Values(conn.Do("ZREVRANGE", key, 0, -1, "WITHSCORES"))
	if err != nil {
		log.Error("conn.Do(ZREVRANGE, %s) error(%v)", key, err)
		return
	}
	if len(values) == 0 {
		return
	}
	if err = conn.Flush(); err != nil {
		log.Error("conn.Flush() err(%v)", err)
		return
	}
	var num int64
	for len(values) > 0 {
		bs := []byte{}
		if values, err = redis.Scan(values, &bs, &num); err != nil {
			log.Error("redis.Scan(%v) error(%v)", values, err)
			return
		}
		detail := &faq.Detail{}
		if err = json.Unmarshal(bs, detail); err != nil {
			log.Error("json.Unmarshal(%v) error(%v)", bs, err)
			return
		}
		res = append(res, detail)
	}
	count = fromHd(num)
	return
}

func fromHd(i int64) int {
	return int(i & 0xffff)
}

func combineHd(create xtime.Time, count int) int64 {
	return create.Time().Unix()<<16 | int64(count)
}

// Detail fn
func (d *Dao) Detail(c context.Context, qTypeID string, keyFlag, pn, ps int) (data []*faq.Detail, total int, err error) {
	ip := metadata.String(c, metadata.RemoteIP)
	params := url.Values{}
	params.Set("questionTypeId", qTypeID)
	params.Set("keyFlag", strconv.Itoa(keyFlag))
	params.Set("keyWords", "")
	params.Set("pageNo", strconv.Itoa(pn))
	params.Set("pageSize", strconv.Itoa(ps))
	params.Set("robotFlag", strconv.Itoa(_notRobot))
	var res struct {
		Code  string        `json:"retCode"`
		Data  []*faq.Detail `json:"items"`
		Total int           `json:"totalCount"`
	}
	if err = d.client.Get(c, d.searchURL, ip, params, &res); err != nil {
		log.Error("Detail d.searchURL url(%s)|(%+v)", d.searchURL+"?"+params.Encode(), res)
		err = ecode.HelpDetailError
		return
	}
	log.Info("Detail d.searchURL url(%s)|(%+v)", d.searchURL+"?"+params.Encode(), res.Code)
	if res.Code != _rsOk {
		log.Error("Detail d.searchURL url(%s)|(%+v)", d.searchURL+"?"+params.Encode(), res.Code)
		err = ecode.HelpDetailError
		return
	}
	total = res.Total
	data = res.Data
	return
}
