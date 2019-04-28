package dao

import (
	"context"
	"encoding/json"
	"fmt"

	"go-common/app/interface/main/web/model"
	artmdl "go-common/app/interface/openplatform/article/model"
	"go-common/library/cache/redis"
	"go-common/library/log"
)

const (
	_keyArtList = "art_%d_%d"
	_keyArtUp   = "art_u"
)

func keyArtList(rid int64, sort int) string {
	return fmt.Sprintf(_keyArtList, rid, sort)
}

// ArticleListCache get article list cache
func (d *Dao) ArticleListCache(c context.Context, rid int64, sort int) (res []*artmdl.Meta, err error) {
	var (
		value []byte
		key   = keyArtList(rid, sort)
		conn  = d.redis.Get(c)
	)
	defer conn.Close()
	if value, err = redis.Bytes(conn.Do("GET", key)); err != nil {
		if err == redis.ErrNil {
			err = nil
		} else {
			log.Error("conn.Do(GET, %s) error(%v)", key, err)
		}
		return
	}
	res = []*artmdl.Meta{}
	if err = json.Unmarshal(value, &res); err != nil {
		log.Error("json.Unmarshal(%v) error(%v)", value, err)
	}
	return
}

// SetArticleListCache set article list cache
func (d *Dao) SetArticleListCache(c context.Context, rid int64, sort int, list []*artmdl.Meta) (err error) {
	var (
		bs   []byte
		key  = keyArtList(rid, sort)
		conn = d.redis.Get(c)
	)
	defer conn.Close()
	if bs, err = json.Marshal(list); err != nil {
		log.Error("json.Marshal(%v) error (%v)", list, err)
		return
	}
	if err = conn.Send("SET", key, bs); err != nil {
		log.Error("conn.Send(SET, %s, %s) error(%v)", key, string(bs), err)
		return
	}
	if err = conn.Send("EXPIRE", key, d.redisArtBakExpire); err != nil {
		log.Error("conn.Send(Expire, %s, %d) error(%v)", key, d.redisArtBakExpire, err)
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

// ArticleUpListCache get article up list cache.
func (d *Dao) ArticleUpListCache(c context.Context) (res []*model.Info, err error) {
	var (
		value []byte
		conn  = d.redis.Get(c)
	)
	defer conn.Close()
	if value, err = redis.Bytes(conn.Do("GET", _keyArtUp)); err != nil {
		if err == redis.ErrNil {
			err = nil
		} else {
			log.Error("conn.Do(GET, %s) error(%v)", _keyArtUp, err)
		}
		return
	}
	res = []*model.Info{}
	if err = json.Unmarshal(value, &res); err != nil {
		log.Error("json.Unmarshal(%v) error(%v)", value, err)
	}
	return
}

// SetArticleUpListCache set article up list cache.
func (d *Dao) SetArticleUpListCache(c context.Context, list []*model.Info) (err error) {
	var (
		bs   []byte
		key  = _keyArtUp
		conn = d.redis.Get(c)
	)
	defer conn.Close()
	if bs, err = json.Marshal(list); err != nil {
		log.Error("json.Marshal(%v) error (%v)", list, err)
		return
	}
	if err = conn.Send("SET", key, bs); err != nil {
		log.Error("conn.Send(SET, %s, %s) error(%v)", key, string(bs), err)
		return
	}
	if err = conn.Send("EXPIRE", key, d.redisArtBakExpire); err != nil {
		log.Error("conn.Send(Expire, %s, %d) error(%v)", key, d.redisArtBakExpire, err)
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
