package dao

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	pb "go-common/app/service/main/history/api/grpc"
	"go-common/app/service/main/history/model"
	"go-common/library/cache/redis"
	"go-common/library/log"
)

const (
	_keySwitch = "s_" // mid -> bit(value)
	_bucket    = 1000 // bit bucket
)

// keyHistory return history key.
func keyHistory(business string, mid int64) string {
	return fmt.Sprintf("h_%d_%s", mid, business)
}

// keyIndex return history index key.
func keyIndex(business string, mid int64) string {
	return fmt.Sprintf("i_%d_%s", mid, business)
}

// keySwitch return Switch key.
func keySwitch(mid int64) string {
	return _keySwitch + strconv.FormatInt(mid/_bucket, 10)
}

// ListCacheByTime get aids from redis where score.
func (d *Dao) ListCacheByTime(c context.Context, business string, mid int64, start int64) (aids []int64, err error) {
	conn := d.redis.Get(c)
	defer conn.Close()
	values, err := redis.Values(conn.Do("ZRANGEBYSCORE", keyIndex(business, mid), start, "INF", "WITHSCORES"))
	if err != nil {
		log.Error("conn.Do(ZRANGEBYSCORE %v) error(%v)", keyIndex(business, mid), err)
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

// ListsCacheByTime get aids from redis where score.
func (d *Dao) ListsCacheByTime(c context.Context, businesses []string, mid int64, viewAt, ps int64) (res map[string][]int64, err error) {
	conn := d.redis.Get(c)
	defer conn.Close()
	var count int
	for _, business := range businesses {
		if err = conn.Send("ZREVRANGEBYSCORE", keyIndex(business, mid), viewAt, "-INF", "LIMIT", 0, ps); err != nil {
			log.Error("conn.Do(ZRANGEBYSCORE %v) error(%v)", keyIndex(business, mid), err)
			return
		}
		count++
	}
	if err = conn.Flush(); err != nil {
		log.Error("conn.Flush() error(%v)", err)
		return
	}
	for i := 0; i < count; i++ {
		var values []int64
		values, err = redis.Int64s(conn.Receive())
		if err != nil {
			if err == redis.ErrNil {
				err = nil
				continue
			}
			log.Error("receive error(%v)", err)
			return
		}
		if len(values) == 0 {
			continue
		}
		if res == nil {
			res = make(map[string][]int64)
		}
		res[businesses[i]] = values
	}
	return
}

// SetUserHideCache set the user hide to redis.
func (d *Dao) SetUserHideCache(c context.Context, mid, value int64) (err error) {
	key := keySwitch(mid)
	conn := d.redis.Get(c)
	if _, err = conn.Do("HSET", key, mid%_bucket, value); err != nil {
		log.Error("conn.Do(HSET %s,%d) error(%v)", key, value, err)
	}
	conn.Close()
	return
}

// UserHideCache return user hide state from redis.
func (d *Dao) UserHideCache(c context.Context, mid int64) (value int64, err error) {
	key := keySwitch(mid)
	conn := d.redis.Get(c)
	defer conn.Close()
	if value, err = redis.Int64(conn.Do("HGET", key, mid%_bucket)); err != nil {
		if err == redis.ErrNil {
			return model.HideStateNotFound, nil
		}
		log.Error("conn.Do(HGET %s) error(%v)", key, err)
	}
	return
}

// HistoriesCache return the user histories from redis.
func (d *Dao) HistoriesCache(c context.Context, mid int64, hs map[string][]int64) (res map[string]map[int64]*model.History, err error) {
	var (
		values, businesses []string
		aid                int64
		k                  int
	)
	conn := d.redis.Get(c)
	defer conn.Close()
	for business, aids := range hs {
		businesses = append(businesses, business)
		key := keyHistory(business, mid)
		args := []interface{}{key}
		for _, aid := range aids {
			args = append(args, aid)
		}
		if err = conn.Send("HMGET", args...); err != nil {
			log.Error("conn.Send(HMGET %v %v) error(%v)", key, aids, err)
			return
		}
	}
	if err = conn.Flush(); err != nil {
		log.Error("conn.Flush() error(%v)", err)
		return
	}
	for i := 0; i < len(hs); i++ {
		if values, err = redis.Strings(conn.Receive()); err != nil {
			log.Error("conn.Receive error(%v)", err)
			if err == redis.ErrNil {
				continue
			}
			return
		}
		if res == nil {
			res = make(map[string]map[int64]*model.History)
		}
		business := businesses[i]
		for k, aid = range hs[business] {
			if values[k] == "" {
				continue
			}
			h := &model.History{}
			if err = json.Unmarshal([]byte(values[k]), h); err != nil {
				log.Error("json.Unmarshal(%s) error(%v)", values[k], err)
				err = nil
				continue
			}
			h.BusinessID = d.BusinessNames[h.Business].ID
			if res[business] == nil {
				res[business] = make(map[int64]*model.History)
			}
			res[business][aid] = h
		}
	}
	return
}

// ClearHistoryCache clear the user redis.
func (d *Dao) ClearHistoryCache(c context.Context, mid int64, businesses []string) (err error) {
	var conn = d.redis.Get(c)
	var count int
	defer conn.Close()
	for _, business := range businesses {
		idxKey := keyIndex(business, mid)
		key := keyHistory(business, mid)
		if err = conn.Send("DEL", idxKey); err != nil {
			log.Error("conn.Send(DEL %s) error(%v)", idxKey, err)
			return
		}
		count++
		if err = conn.Send("DEL", key); err != nil {
			log.Error("conn.Send(DEL %s) error(%v)", key, err)
			return
		}
		count++
	}
	if err = conn.Flush(); err != nil {
		log.Error("conn.Flush() error(%v)", err)
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

// DelHistoryCache delete the history redis.
func (d *Dao) DelHistoryCache(c context.Context, arg *pb.DelHistoriesReq) (err error) {
	conn := d.redis.Get(c)
	defer conn.Close()
	var count int
	for _, r := range arg.Records {
		var (
			indxKey = keyIndex(r.Business, arg.Mid)
			key     = keyHistory(r.Business, arg.Mid)
		)
		if err = conn.Send("ZREM", indxKey, r.ID); err != nil {
			log.Error("conn.Send(ZREM %s,%v) error(%v)", indxKey, r.ID, err)
			return
		}
		count++
		if err = conn.Send("HDEL", key, r.ID); err != nil {
			log.Error("conn.Send(HDEL %s,%v) error(%v)", key, r.ID, err)
			return
		}
		count++
	}
	if err = conn.Flush(); err != nil {
		log.Error("conn.Flush() error(%v)", err)
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

// AddHistoryCache add the history to redis.
func (d *Dao) AddHistoryCache(c context.Context, h *pb.AddHistoryReq) (err error) {
	var (
		b           []byte
		mid         = h.Mid
		idxKey, key = keyIndex(h.Business, mid), keyHistory(h.Business, mid)
	)
	if b, err = json.Marshal(h); err != nil {
		return
	}
	conn := d.redis.Get(c)
	defer conn.Close()
	if err = conn.Send("ZADD", idxKey, h.ViewAt, h.Kid); err != nil {
		log.Error("conn.Send(ZADD %s,%d) error(%v)", key, h.Kid, err)
		return
	}
	if err = conn.Send("HSET", key, h.Kid, string(b)); err != nil {
		log.Error("conn.Send(HSET %s,%d) error(%v)", key, h.Kid, err)
		return
	}
	if err = conn.Send("EXPIRE", idxKey, d.redisExpire); err != nil {
		log.Error("conn.Send(EXPIRE) error(%v)", err)
		return
	}
	if err = conn.Send("EXPIRE", key, d.redisExpire); err != nil {
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

// AddHistoriesCache add the user to redis.
func (d *Dao) AddHistoriesCache(c context.Context, hs []*pb.AddHistoryReq) (err error) {
	if len(hs) == 0 {
		return
	}
	conn := d.redis.Get(c)
	defer conn.Close()
	var count int
	for _, h := range hs {
		var (
			b           []byte
			mid         = h.Mid
			idxKey, key = keyIndex(h.Business, mid), keyHistory(h.Business, mid)
		)
		if b, err = json.Marshal(h); err != nil {
			continue
		}
		if err = conn.Send("ZADD", idxKey, h.ViewAt, h.Kid); err != nil {
			log.Error("conn.Send(ZADD %s,%d) error(%v)", key, h.Kid, err)
			continue
		}
		count++
		if err = conn.Send("HSET", key, h.Kid, string(b)); err != nil {
			log.Error("conn.Send(HSET %s,%d) error(%v)", key, h.Kid, err)
			continue
		}
		count++
		if err = conn.Send("EXPIRE", idxKey, d.redisExpire); err != nil {
			log.Error("conn.Send(EXPIRE) error(%v)", err)
			continue
		}
		count++
		if err = conn.Send("EXPIRE", key, d.redisExpire); err != nil {
			log.Error("conn.Send(EXPIRE) error(%v)", err)
			continue
		}
		count++
	}
	if err = conn.Flush(); err != nil {
		log.Error("conn.Flush() error(%v)", err)
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

// TrimCache trim history.
func (d *Dao) TrimCache(c context.Context, business string, mid int64, limit int) (err error) {
	conn := d.redis.Get(c)
	defer conn.Close()
	aids, err := redis.Int64s(conn.Do("ZRANGE", keyIndex(business, mid), 0, -limit-1))
	if err != nil {
		log.Error("conn.Do(ZRANGE %v) error(%v)", keyIndex(business, mid), err)
		return
	}
	if len(aids) == 0 {
		return
	}
	return d.DelCache(c, business, mid, aids)
}

// DelCache delete the history redis.
func (d *Dao) DelCache(c context.Context, business string, mid int64, aids []int64) (err error) {
	var (
		key1  = keyIndex(business, mid)
		key2  = keyHistory(business, mid)
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

// pingRedis ping redis.
func (d *Dao) pingRedis(c context.Context) (err error) {
	conn := d.redis.Get(c)
	if _, err = conn.Do("SET", "PING", "PONG"); err != nil {
		log.Error("redis: conn.Do(SET,PING,PONG) error(%v)", err)
	}
	conn.Close()
	return
}
