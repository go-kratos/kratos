package dao

import (
	"context"
	"encoding/json"
	"net/url"

	"go-common/app/interface/main/web/model"
	"go-common/library/cache/redis"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/net/metadata"
)

const (
	_wxHotURI     = "hot-weixin.json"
	_wxCacheKey   = "wx_hot"
	_wxBkCacheKey = _keyBakPrefix + _wxCacheKey
)

// WxHot get wx hot aids.
func (d *Dao) WxHot(c context.Context) (aids []int64, err error) {
	ip := metadata.String(c, metadata.RemoteIP)
	var res struct {
		Code int                 `json:"code"`
		List []*model.NewArchive `json:"list"`
	}
	if err = d.httpBigData.Get(c, d.wxHotURL, ip, url.Values{}, &res); err != nil {
		log.Error("d.httpBigData.Get(%s) error(%v)", d.wxHotURL, err)
		return
	}
	if res.Code != ecode.OK.Code() {
		log.Error("d.httpBigData.Get(%s) error(%v)", d.wxHotURL, err)
		err = ecode.Int(res.Code)
		return
	}
	for _, v := range res.List {
		if v.Aid > 0 {
			aids = append(aids, v.Aid)
		}
	}
	return
}

// WxHotCache get wx hot cache.
func (d *Dao) WxHotCache(c context.Context) (arcs []*model.WxArchive, err error) {
	key := _wxCacheKey
	conn := d.redis.Get(c)
	defer conn.Close()
	arcs, err = wxHotCache(conn, key)
	return
}

// WxHotBakCache get wx hot bak cache.
func (d *Dao) WxHotBakCache(c context.Context) (arcs []*model.WxArchive, err error) {
	key := _wxBkCacheKey
	conn := d.redisBak.Get(c)
	defer conn.Close()
	arcs, err = wxHotCache(conn, key)
	return
}

func wxHotCache(conn redis.Conn, key string) (res []*model.WxArchive, err error) {
	var value []byte
	if value, err = redis.Bytes(conn.Do("GET", key)); err != nil {
		if err == redis.ErrNil {
			err = nil
		} else {
			log.Error("conn.Do(GET, %s) error(%v)", key, err)
		}
		return
	}
	res = []*model.WxArchive{}
	if err = json.Unmarshal(value, &res); err != nil {
		log.Error("json.Unmarshal(%v) error(%v)", value, err)
	}
	return
}

// SetWxHotCache set wx hot to cache.
func (d *Dao) SetWxHotCache(c context.Context, arcs []*model.WxArchive) (err error) {
	key := _wxCacheKey
	conn := d.redis.Get(c)
	defer conn.Close()
	if err = d.setWxHotCache(c, conn, key, d.redisWxHotExpire, arcs); err != nil {
		return
	}
	key = _wxBkCacheKey
	connBak := d.redisBak.Get(c)
	err = d.setWxHotCache(c, connBak, key, d.redisWxHotBakExpire, arcs)
	connBak.Close()
	return
}

func (d *Dao) setWxHotCache(c context.Context, conn redis.Conn, key string, expire int32, arcs []*model.WxArchive) (err error) {
	var bs []byte
	if bs, err = json.Marshal(arcs); err != nil {
		log.Error("json.Marshal(%v) error (%v)", arcs, err)
		return
	}
	if err = conn.Send("SET", key, bs); err != nil {
		log.Error("conn.Send(SET, %s, %s) error(%v)", key, string(bs), err)
		return
	}
	if err = conn.Send("EXPIRE", key, expire); err != nil {
		log.Error("conn.Send(Expire, %s, %d) error(%v)", key, expire, err)
		return
	}
	if err = conn.Flush(); err != nil {
		log.Error("conn.Flush error(%v)", err)
		return
	}
	for i := 0; i < 2; i++ {
		if _, err = conn.Receive(); err != nil {
			log.Error("conn.Receive() error(%v)", err)
			return
		}
	}
	return
}
