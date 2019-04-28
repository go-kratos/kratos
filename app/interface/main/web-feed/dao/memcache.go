package dao

import (
	"context"
	"encoding/json"
	"strconv"

	"go-common/app/interface/main/web-feed/model"
	"go-common/library/cache/memcache"
	"go-common/library/log"
)

const (
	_prefixFeed = "f_"
)

func feedKey(mid int64) string {
	return _prefixFeed + strconv.FormatInt(mid, 10)
}

// pingMc ping memcache
func (d *Dao) pingMC(c context.Context) (err error) {
	conn := d.mc.Get(c)
	defer conn.Close()
	item := memcache.Item{Key: "ping", Value: []byte{1}, Expiration: 0}
	err = conn.Set(&item)
	return
}

// FeedCache gets feed cache
func (d *Dao) FeedCache(c context.Context, mid int64) (feeds []*model.Feed, err error) {
	PromInfo("backup-cache")
	var (
		reply *memcache.Item
		conn  = d.mc.Get(c)
		key   = feedKey(mid)
	)
	defer conn.Close()
	if reply, err = conn.Get(key); err != nil {
		if err == memcache.ErrNotFound {
			err = nil
		} else {
			PromError("mc:获取feed缓存")
			log.Error("conn.Get(%s) error(%v)", key, err)
		}
		return
	}
	err = conn.Scan(reply, &feeds)
	return
}

// SetFeedCache sets feed cache
func (d *Dao) SetFeedCache(c context.Context, mid int64, feeds []*model.Feed) (err error) {
	var (
		bs   []byte
		key  = feedKey(mid)
		conn = d.mc.Get(c)
	)
	defer conn.Close()
	if bs, err = json.Marshal(feeds); err != nil {
		PromError("mc:json feed缓存")
		log.Error("json.Marshal(%v) error(%v)", feeds, err)
		return
	}
	item := &memcache.Item{Key: key, Value: bs, Expiration: d.mcFeedExpire}
	if err = conn.Set(item); err != nil {
		PromError("mc:设置feed缓存")
		log.Error("conn.Set(%+v) error(%v)", item, err)
	}
	return
}
