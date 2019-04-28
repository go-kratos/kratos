package history

import (
	"context"
	"encoding/json"
	"strconv"

	"go-common/app/interface/main/history/model"
	"go-common/library/cache/redis"
	"go-common/library/log"
)

const (
	_keyIndex   = "i_" // mid -> score:time member:aid
	_keyHistory = "h_" // mid -> hash(aid,progress)
	_keySwitch  = "s_" // mid -> bit(value)
	_bucket     = 1000 // bit bucket
)

// keyHistory return history key.
func keyHistory(mid int64) string {
	return _keyHistory + strconv.FormatInt(mid, 10)
}

// keyIndex return history index key.
func keyIndex(mid int64) string {
	return _keyIndex + strconv.FormatInt(mid, 10)
}

// keySwitch return Switch key.
func keySwitch(mid int64) string {
	return _keySwitch + strconv.FormatInt(mid/_bucket, 10)
}

// ExpireIndexCache expire index cache.
func (d *Dao) ExpireIndexCache(c context.Context, mid int64) (bool, error) {
	return d.expireCache(c, keyIndex(mid))
}

// ExpireCache expire the user State redis.
func (d *Dao) ExpireCache(c context.Context, mid int64) (bool, error) {
	return d.expireCache(c, keyHistory(mid))
}

// ExpireToView expire toview by mid.
func (d *Dao) expireCache(c context.Context, key string) (ok bool, err error) {
	conn := d.redis.Get(c)
	if ok, err = redis.Bool(conn.Do("EXPIRE", key, d.expire)); err != nil {
		log.Error("conn.Do(EXPIRE %s) error(%v)", key, err)
	}
	conn.Close()
	return
}

// IndexCache get history from redis.
func (d *Dao) IndexCache(c context.Context, mid int64, start, end int) (aids []int64, err error) {
	conn := d.redis.Get(c)
	defer conn.Close()
	values, err := redis.Values(conn.Do("ZREVRANGE", keyIndex(mid), start, end, "WITHSCORES"))
	if err != nil {
		log.Error("conn.Do(ZREVRANGE %v) error(%v)", keyIndex(mid), err)
		return
	}
	if len(values) == 0 {
		return
	}
	var aid, unix int64
	for len(values) > 0 {
		if values, err = redis.Scan(values, &aid, &unix); err != nil {
			log.Error("redis.Scan(%v) error(%v)", values, err)
			return
		}
		aids = append(aids, aid)
	}
	return
}

// IndexCacheByTime get aids from redis where score.
func (d *Dao) IndexCacheByTime(c context.Context, mid int64, start int64) (aids []int64, err error) {
	conn := d.redis.Get(c)
	defer conn.Close()
	values, err := redis.Values(conn.Do("ZRANGEBYSCORE", keyIndex(mid), start, "INF", "WITHSCORES"))
	if err != nil {
		log.Error("conn.Do(ZRANGEBYSCORE %v) error(%v)", keyIndex(mid), err)
		return
	}
	if len(values) == 0 {
		return
	}
	var aid, unix int64
	for len(values) > 0 {
		if values, err = redis.Scan(values, &aid, &unix); err != nil {
			log.Error("redis.Scan(%v) error(%v)", values, err)
			return
		}
		aids = append(aids, aid)
	}
	return
}

// SetShadowCache set the user switch to redis.
func (d *Dao) SetShadowCache(c context.Context, mid, value int64) (err error) {
	key := keySwitch(mid)
	conn := d.redis.Get(c)
	if _, err = conn.Do("HSET", key, mid%_bucket, value); err != nil {
		log.Error("conn.Do(HSET %s,%d) error(%v)", key, value, err)
	}
	conn.Close()
	return
}

// ShadowCache return user switch redis.
func (d *Dao) ShadowCache(c context.Context, mid int64) (value int64, err error) {
	key := keySwitch(mid)
	conn := d.redis.Get(c)
	defer conn.Close()
	if value, err = redis.Int64(conn.Do("HGET", key, mid%_bucket)); err != nil {
		if err == redis.ErrNil {
			return model.ShadowUnknown, nil
		}
		log.Error("conn.Do(HGET %s) error(%v)", key, err)
	}
	return
}

// CacheMap return the user State from redis.
func (d *Dao) CacheMap(c context.Context, mid int64) (amap map[int64]*model.History, err error) {
	var (
		values map[string]string
		key    = keyHistory(mid)
	)
	conn := d.redis.Get(c)
	defer conn.Close()
	if values, err = redis.StringMap(conn.Do("HGETALL", key)); err != nil {
		log.Error("conn.Do(HGETALL %v) error(%v)", key, err)
		if err == redis.ErrNil {
			return nil, nil
		}
		return
	}
	amap = make(map[int64]*model.History)
	for k, v := range values {
		if v == "" {
			continue
		}
		h := &model.History{}
		if err = json.Unmarshal([]byte(v), h); err != nil {
			log.Error("json.Unmarshal(%s) error(%v)", v, err)
			err = nil
			continue
		}
		if h.TP == model.TypeOffline {
			h.TP = model.TypeUGC
		}
		h.Mid = mid
		h.Aid, _ = strconv.ParseInt(k, 10, 0)
		h.FillBusiness()
		amap[h.Aid] = h
	}
	return
}

// Cache return the user State from redis.
func (d *Dao) Cache(c context.Context, mid int64, aids []int64) (amap map[int64]*model.History, miss []int64, err error) {
	var (
		values []string
		aid    int64
		k      int
		key    = keyHistory(mid)
		args   = []interface{}{key}
	)
	for _, aid = range aids {
		args = append(args, aid)
	}
	conn := d.redis.Get(c)
	defer conn.Close()
	if values, err = redis.Strings(conn.Do("HMGET", args...)); err != nil {
		log.Error("conn.Do(HMGET %v) error(%v)", args, err)
		if err == redis.ErrNil {
			return nil, nil, nil
		}
		return
	}
	amap = make(map[int64]*model.History, len(aids))
	for k, aid = range aids {
		if values[k] == "" {
			miss = append(miss, aid)
			continue
		}
		h := &model.History{}
		if err = json.Unmarshal([]byte(values[k]), h); err != nil {
			log.Error("json.Unmarshal(%s) error(%v)", values[k], err)
			err = nil
			continue
		}
		if h.TP == model.TypeOffline {
			h.TP = model.TypeUGC
		}
		h.Mid = mid
		h.Aid = aid
		h.FillBusiness()
		amap[aid] = h
	}
	return
}

// ClearCache clear the user State redis.
func (d *Dao) ClearCache(c context.Context, mid int64) (err error) {
	var (
		idxKey = keyIndex(mid)
		key    = keyHistory(mid)
		conn   = d.redis.Get(c)
	)
	defer conn.Close()
	if err = conn.Send("DEL", idxKey); err != nil {
		log.Error("conn.Send(DEL %s) error(%v)", idxKey, err)
		return
	}
	if err = conn.Send("DEL", key); err != nil {
		log.Error("conn.Send(DEL %s) error(%v)", key, err)
		return
	}
	if err = conn.Flush(); err != nil {
		log.Error("conn.Flush() error(%v)", err)
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

// DelCache delete the history redis.
func (d *Dao) DelCache(c context.Context, mid int64, aids []int64) (err error) {
	var (
		key1  = keyIndex(mid)
		key2  = keyHistory(mid)
		args1 = []interface{}{key1}
		args2 = []interface{}{key2}
	)
	for _, aid := range aids {
		args1 = append(args1, aid)
		args2 = append(args2, aid)
	}
	conn := d.redis.Get(c)
	defer conn.Close()
	if err = conn.Send("ZREM", args1...); err != nil {
		log.Error("conn.Send(ZREM %s,%v) error(%v)", key1, aids, err)
		return
	}
	if err = conn.Send("HDEL", args2...); err != nil {
		log.Error("conn.Send(HDEL %s,%v) error(%v)", key2, aids, err)
		return
	}
	if err = conn.Flush(); err != nil {
		log.Error("conn.Flush() error(%v)", err)
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

// AddCache add the user State to redis.
func (d *Dao) AddCache(c context.Context, mid int64, h *model.History) (err error) {
	var (
		b           []byte
		idxKey, key = keyIndex(mid), keyHistory(mid)
	)
	if b, err = json.Marshal(h); err != nil {
		return
	}
	conn := d.redis.Get(c)
	defer conn.Close()
	if err = conn.Send("ZADD", idxKey, h.Unix, h.Aid); err != nil {
		log.Error("conn.Send(ZADD %s,%d) error(%v)", key, h.Aid, err)
		return
	}
	if err = conn.Send("HSET", key, h.Aid, string(b)); err != nil {
		log.Error("conn.Send(HSET %s,%d) error(%v)", key, h.Aid, err)
		return
	}
	if err = conn.Send("EXPIRE", idxKey, d.expire); err != nil {
		log.Error("conn.Send(EXPIRE) error(%v)", err)
		return
	}
	if err = conn.Send("EXPIRE", key, d.expire); err != nil {
		log.Error("conn.Send(EXPIRE) error(%v)", err)
		return
	}
	if err = conn.Flush(); err != nil {
		log.Error("conn.Flush() error(%v)", err)
		return
	}
	for i := 0; i < 2+2; i++ {
		if _, err = conn.Receive(); err != nil {
			log.Error("conn.Receive() error(%v)", err)
			return
		}
	}
	return
}

// AddCacheMap add the user State to redis.
func (d *Dao) AddCacheMap(c context.Context, mid int64, hm map[int64]*model.History) (err error) {
	var (
		b           []byte
		idxKey, key = keyIndex(mid), keyHistory(mid)
	)
	conn := d.redis.Get(c)
	defer conn.Close()
	for _, h := range hm {
		if b, err = json.Marshal(h); err != nil {
			continue
		}
		if err = conn.Send("ZADD", idxKey, h.Unix, h.Aid); err != nil {
			log.Error("conn.Send(ZADD %s,%d) error(%v)", key, h.Aid, err)
			continue
		}
		if err = conn.Send("HSET", key, h.Aid, string(b)); err != nil {
			log.Error("conn.Send(HSET %s,%d) error(%v)", key, h.Aid, err)
			continue
		}
	}
	if err = conn.Send("EXPIRE", idxKey, d.expire); err != nil {
		log.Error("conn.Send(EXPIRE) error(%v)", err)
		return
	}
	if err = conn.Send("EXPIRE", key, d.expire); err != nil {
		log.Error("conn.Send(EXPIRE) error(%v)", err)
		return
	}
	if err = conn.Flush(); err != nil {
		log.Error("conn.Flush() error(%v)", err)
		return
	}
	for i := 0; i < len(hm)*2+2; i++ {
		if _, err = conn.Receive(); err != nil {
			log.Error("conn.Receive() error(%v)", err)
			return
		}
	}
	return
}

// TrimCache trim history.
func (d *Dao) TrimCache(c context.Context, mid int64, limit int) (err error) {
	conn := d.redis.Get(c)
	defer conn.Close()
	aids, err := redis.Int64s(conn.Do("ZRANGE", keyIndex(mid), 0, -limit-1))
	if err != nil {
		log.Error("conn.Do(ZRANGE %v) error(%v)", keyIndex(mid), err)
		return
	}
	if len(aids) == 0 {
		return
	}
	return d.DelCache(c, mid, aids)
}

// PingRedis check redis connection
func (d *Dao) PingRedis(c context.Context) (err error) {
	return
}
