package toview

import (
	"context"
	"strconv"

	"go-common/app/interface/main/history/model"
	"go-common/library/cache/redis"
	"go-common/library/log"
)

const _key = "v_" // mid -> score:time member:aid

// keyToView return key string
func key(mid int64) string {
	return _key + strconv.FormatInt(mid, 10)
}

// Expire expire toview by mid.
func (d *Dao) Expire(c context.Context, mid int64) (ok bool, err error) {
	conn := d.redis.Get(c)
	if ok, err = redis.Bool(conn.Do("EXPIRE", key(mid), d.expire)); err != nil {
		log.Error("conn.Do(EXPIRE, %s) error(%v)", key(mid), err)
	}
	conn.Close()
	return
}

// Cache return the user all toview from redis.
func (d *Dao) Cache(c context.Context, mid int64, start, end int) (res []*model.ToView, err error) {
	conn := d.redis.Get(c)
	defer conn.Close()
	values, err := redis.Values(conn.Do("ZREVRANGE", key(mid), start, end, "WITHSCORES"))
	if err != nil {
		log.Error("dao.Do(ZREVRANGE %v) error(%v)", key(mid), err)
		return
	}
	if len(values) == 0 {
		return
	}
	res = make([]*model.ToView, 0, len(values)/2)
	for len(values) > 0 {
		t := &model.ToView{}
		if values, err = redis.Scan(values, &t.Aid, &t.Unix); err != nil {
			log.Error("redis.Scan(%v) error(%v)", values, err)
			return
		}
		res = append(res, t)
	}
	return
}

// CacheMap return the user all toview map from redis.
func (d *Dao) CacheMap(c context.Context, mid int64) (res map[int64]*model.ToView, err error) {
	conn := d.redis.Get(c)
	defer conn.Close()
	values, err := redis.Values(conn.Do("ZREVRANGE", key(mid), 0, -1, "WITHSCORES"))
	if err != nil {
		log.Error("dao.Do(ZREVRANGE %v) error(%v)", key(mid), err)
		return
	}
	if len(values) == 0 {
		return
	}
	res = make(map[int64]*model.ToView, len(values)/2)
	for len(values) > 0 {
		t := &model.ToView{}
		if values, err = redis.Scan(values, &t.Aid, &t.Unix); err != nil {
			log.Error("redis.Scan(%v) error(%v)", values, err)
			return
		}
		res[t.Aid] = t
	}
	return
}

// CntCache return the user toview count from redis.
func (d *Dao) CntCache(c context.Context, mid int64) (count int, err error) {
	conn := d.redis.Get(c)
	if count, err = redis.Int(conn.Do("ZCARD", key(mid))); err != nil {
		log.Error("dao.Do(ZCARD,%s) err(%v)", key(mid), err)
	}
	conn.Close()
	return
}

// ClearCache delete the user toview redis.
func (d *Dao) ClearCache(c context.Context, mid int64) (err error) {
	conn := d.redis.Get(c)
	if _, err = conn.Do("DEL", key(mid)); err != nil {
		log.Error("conn.Do(DEL %s) error(%v)", key(mid), err)
	}
	conn.Close()
	return
}

// DelCaches delete the user toview redis.
func (d *Dao) DelCaches(c context.Context, mid int64, aids []int64) (err error) {
	conn := d.redis.Get(c)
	defer conn.Close()
	for _, aid := range aids {
		if err = conn.Send("ZREM", key(mid), aid); err != nil {
			log.Error("conn.Send(ZREM %s,%d) error(%v)", key(mid), aid, err)
			return
		}
	}
	if err = conn.Flush(); err != nil {
		log.Error("conn.Flush() error(%v)", err)
		return
	}
	for i := 0; i < len(aids); i++ {
		if _, err = conn.Receive(); err != nil {
			log.Error("conn.Receive() error(%v)", err)
			return
		}
	}
	return
}

// AddCache add user toview to redis.
func (d *Dao) AddCache(c context.Context, mid, aid, now int64) error {
	return d.addCache(c, key(mid), []*model.ToView{&model.ToView{Aid: aid, Unix: now}})
}

// AddCacheList add user toview to redis.
func (d *Dao) AddCacheList(c context.Context, mid int64, views []*model.ToView) error {
	return d.addCache(c, key(mid), views)
}

// addCache add user toview to redis.
func (d *Dao) addCache(c context.Context, key string, views []*model.ToView) (err error) {
	conn := d.redis.Get(c)
	defer conn.Close()
	for _, v := range views {
		if err = conn.Send("ZADD", key, v.Unix, v.Aid); err != nil {
			log.Error("conn.Send(ZREM %s,%d) error(%v)", key, v.Aid, err)
			return
		}
	}
	if err = conn.Send("EXPIRE", key, d.expire); err != nil {
		log.Error("conn.Send(EXPIRE) error(%v)", err)
		return
	}
	if err = conn.Flush(); err != nil {
		log.Error("conn.Flush() error(%v)", err)
		return
	}
	for i := 0; i < len(views)+1; i++ {
		if _, err = conn.Receive(); err != nil {
			log.Error("conn.Receive() error(%v)", err)
			return
		}
	}
	return
}

// PingRedis check redis connection
func (d *Dao) PingRedis(c context.Context) (err error) {
	return
}
